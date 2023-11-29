#!python3

import argparse
import subprocess
import sys
import os
import json
import re
import hashlib
import signal
import time
import itertools
import traceback
import tempfile

from threading import Event

from projects import projects
from version import Version
import github
from changelog import CHANGELOGs
from readme import READMEs
from javascript import package_versions
from git import uncommitted
from misc import say
import build

exit = Event()


def quit(signo, _frame):
    print("Interrupted by %d, shutting down" % signo)
    exit.set()


def main():
    print()
    print("*** uhppoted release script (v2) ***")
    print()

    if len(sys.argv) < 2:
        usage()
        return -1

    parser = argparse.ArgumentParser(description='release --version=<version>')
    parser.add_argument('--version', type=str, default='development', help='release version e.g. v0.8.7')
    parser.add_argument('--node-red', type=str, default='development', help='NodeRED release version e.g. v1.1.2')

    args = parser.parse_args()
    version = Version(args.version, args.node_red)

    print(f'VERSION: {version}')
    print()

    try:
        # ... get release state
        state = get_release_info(version)
        plist = projects()

        # ... get unreleased project list
        if not 'unreleased' in state:
            it = itertools.filterfalse(lambda p: github.already_released(p, plist[p], version.version(p)), plist)
            state['unreleased'] = {p: plist[p] for p in it}
            save_release_info(version, state)
            print()

        unreleased = state['unreleased']

        # ... CHANGELOG.md
        if not 'changelogs' in state or state['changelogs'] != 'ok':
            CHANGELOGs(unreleased, version, exit)
            print()
            if not exit.is_set():
                state['changelogs'] = 'ok'
                save_release_info(version, state)

        # ... REAMDME.md
        if not 'readmes' in state or state['readmes'] != 'ok':
            READMEs(unreleased, version, exit)
            print()
            if not exit.is_set():
                state['readmes'] = 'ok'
                save_release_info(version, state)

        # ... package versions
        if not 'package-versions' in state or state['package-versions'] != 'ok':
            package_versions(unreleased, version, exit)
            print()
            if not exit.is_set():
                state['package-versions'] = 'ok'
                save_release_info(version, state)

        # ... uncommitted changes
        if not 'uncommitted-changes' in state or state['uncommitted-changes'] != 'ok':
            uncommitted(unreleased, version, exit)
            print()
            if not exit.is_set():
                state['uncommitted-changes'] = 'ok'
                save_release_info(version, state)

        # ... 'prepare' build
        if not 'prepare' in state or state['prepare'] != 'ok':
            build.prepare(unreleased, version, exit)
            print()
            if not exit.is_set():
                state['prepare'] = 'ok'
                save_release_info(version, state)

        #     if 'prerelease' in ops:
        #         build.prerelease(unreleased, version, exit)
        #         print()

        #     if 'release' in ops:
        #         github.release_notes(unreleased, version, exit)
        #         build.release(unreleased, version, exit)
        #         github.publish(unreleased, version, exit)
        #         print()

        #     if 'bump' in ops:
        #         if len(unreleased) != 0:
        #             raise Exception(f'Projects {unreleased} have not been released')

        #         for p in plist:
        #             print(f'>>>> bumping {p}')
        #             project = plist[p]
        #             clean_release_notes(p, project)
        #             bump_changelog(p, project)
        #         print()

        #     print()
        #     print(f'*** OK!')
        #     print()
        say('OK')

    except BaseException as x:
        print(traceback.format_exc())

        msg = f'{x}'
        msg = msg.replace('uhppoted-','')                         \
                 .replace('uhppote-','')                          \
                 .replace('uhppoted','umbrella project')          \
                 .replace('cli','[[char LTRL]]cli[[char NORM]]')  \
                 .replace('git','[[inpt PHON]]git[[input TEXT]]') \
                 .replace('codegen','code gen')

        print()
        print(f'*** ERROR  {x}')
        print()

        say('ERROR')
        say(msg)


def get_release_info(version):
    file = f'.release-{version}'

    try:
        with open(file) as bytes:
            return json.load(bytes)
    except FileNotFoundError:
        say(f'no release file found')
        say(f'starting new release {version}')

    return {}


def save_release_info(version, info):
    file = f'.release-{version}'

    with open(file, 'w', encoding='utf-8') as f:
        json.dump(info, f, ensure_ascii=False, indent=4)


# def has_uncommitted_changes(folder):
#     command = f"git -C {folder} status  -uno"
#     result = subprocess.check_output(command, shell=True)
#
#     if 'Changes not staged for commit' in str(result):
#         return True
#
#     return False


def checkout(project, info):
    try:
        command = f"cd {info['folder']} && git checkout {info['branch']}"
        subprocess.run(command, shell=True, check=True)
    except subprocess.CalledProcessError:
        raise Exception(f"command 'checkout {project}' failed")


def update(project, info):
    try:
        command = f"cd {info['folder']} && make update && make build"
        subprocess.run(command, shell=True, check=True)
    except subprocess.CalledProcessError:
        raise Exception(f"command 'update {project}' failed")


def git(project, info, interim):
    try:
        command = f"cd {info['folder']} && git remote update"
        subprocess.run(command, shell=True, check=True)

        command = f"cd {info['folder']} && git status -uno"
        result = subprocess.check_output(command, shell=True)

        if 'Changes not staged for commit' in str(result):
            raise Exception(f"{project} has uncommitted changes")
        elif (not interim) and 'Your branch is ahead' in str(result):
            raise Exception(f"{project} has commits that have not been pushed")

    except subprocess.CalledProcessError:
        raise Exception(f"{project}: command 'git status' failed")


def release(project, info, version):
    command = f"cd {info['folder']} && make release DIST={project}_{version}"
    result = subprocess.call(command, shell=True)
    if result != 0:
        raise Exception(f"command 'build {project}' failed")


def checksum(project, info, version):
    if 'binary' in info:
        binary = info['binary']
        root = f"{info['folder']}"
        platforms = ['linux', 'darwin', 'windows', 'arm', 'arm7']

        for platform in platforms:
            filename = binary
            if platform == 'windows':
                filename = f'{binary}.exe'

            exe = os.path.join(root, 'dist', f"{project}_{version}", platform, filename)
            combined = os.path.join('dist', platform, f"uhppoted_{version}", filename)

            if version == 'development':
                exe = os.path.join(root, 'dist', f"{version}", platform, filename)
                combined = os.path.join('dist', platform, f"{version}", filename)

            if hash(combined) != hash(exe):
                print(f'{project:<25}  {exe:<82}  {hash(exe)}')
                print(f'{"":<25}  {combined:<82}  {hash(combined)}')
                raise Exception(f"{project} 'dist' checksums differ")


def clean_release_notes(project, info):
    file = f"{info['folder']}/release-notes.md"
    if os.path.isfile(file):
        os.remove(file)
        print(f'     ... {project} removed release-notes.md')


def bump_changelog(project, info):
    with open(f"{info['folder']}/CHANGELOG.md", 'r', encoding="utf-8") as f:
        CHANGELOG = f.read()
        if 'Unreleased' in CHANGELOG:
            return

    tmpfile = tempfile.NamedTemporaryFile(mode="w+t", delete=False)

    try:
        with open(f"{info['folder']}/CHANGELOG.md", 'r', encoding="utf-8") as f:
            CHANGELOG = f.read()
            rest = CHANGELOG
            heading, _, rest = rest.partition('\n')
            spacer, _, rest = rest.partition('\n')

            tmpfile.write(f'{heading}\n')
            tmpfile.write(f'\n')
            tmpfile.write(f'## Unreleased\n')
            tmpfile.write(f'\n')
            tmpfile.write(f'\n')
            tmpfile.write(rest)
            tmpfile.close()

            os.rename(f"{info['folder']}/CHANGELOG.md", f"{info['folder']}/CHANGELOG.bak")
            os.rename(tmpfile.name, f"{info['folder']}/CHANGELOG.md")
            if os.path.isfile(f"{info['folder']}/CHANGELOG.bak"):
                os.remove(f"{info['folder']}/CHANGELOG.bak")
    finally:
        tmpfile.close()
        if os.path.isfile(tmpfile.name):
            os.remove(tmpfile.name)


def hash(file):
    hash = hashlib.sha256()

    with open(file, "rb") as f:
        bytes = f.read(65536)
        hash.update(bytes)

    return hash.hexdigest()


# def say(msg):
#     transliterated = msg.replace('nodejs','node js') \
#                         .replace('codegen', 'code gen') \
#                         .replace('Errno','error number') \
#                         .replace('exe','e x e')
#     subprocess.call(f'say {transliterated}', shell=True)


def usage():
    print()
    print('  Usage: python release.py <options> <command>')
    print()
    print('  Options:')
    print('    --version <version>  Release version e.g. v0.8.1')
    print()


if __name__ == '__main__':
    for sig in ['SIGTERM', 'SIGHUP', 'SIGINT']:
        signal.signal(getattr(signal, sig), quit)

    main()
