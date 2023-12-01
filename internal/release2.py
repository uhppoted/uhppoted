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
    parser.add_argument('--release', action='store_true', help="builds the release versions")

    args = parser.parse_args()
    version = Version(args.version, args.node_red)
    release = args.release

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
            ok = uncommitted(unreleased, version, exit)
            print()
            if ok and not exit.is_set():
                state['uncommitted-changes'] = 'ok'
                save_release_info(version, state)

        # ... 'prepare' builds
        if not 'prepared' in state or state['prepared'] != 'ok':
            build.prepare(unreleased, version, exit)
            state['uncommitted-changes'] = '??'
            save_release_info(version, state)
            ok = uncommitted(unreleased, version, exit)
            print()
            if ok and not exit.is_set():
                state['prepared'] = 'ok'
                state['uncommitted-changes'] = 'ok'
                save_release_info(version, state)

        # ... 'release' builds
        if release:
            if not 'release-notes' in state or state['release-notes'] != 'ok':
                ok = github.release_notes(unreleased, version, exit)
                print()
                if ok and not exit.is_set():
                    state['release-notes'] = 'ok'
                    save_release_info(version, state)

            if 'uhppote-core' in unreleased:
                ok = build.release('uhppote-core', plist['uhppote-core'], version, exit)
                if ok and not exit.is_set():
                    ok = github.publish('uhppote-core', plist['uhppote-core'], version, exit)
                    if ok and not exit.is_set():
                        del unreleased['uhppote-core']
                        # del state['unreleased']['uhppote-core']
                        save_release_info(version, state)

            if 'uhppoted-lib' in unreleased:
                ok = build.release('uhppoted-lib', plist['uhppoted-lib'], version, exit)
                if ok and not exit.is_set():
                    ok = github.publish('uhppoted-lib', plist['uhppoted-lib'], version, exit)
                    if ok and not exit.is_set():
                        del unreleased['uhppoted-lib']
                        save_release_info(version, state)

            ignore = ['uhppoted', 'uhppoted-nodejs', 'node-red-contrib-uhppoted', 'uhppoted-python']
            it = itertools.filterfalse(lambda p: p in ignore, unreleased)
            rl = {p: plist[p] for p in it}
            for p in rl:
                ok = build.release(p, plist[p], version, exit)
                if ok and not exit.is_set():
                    ok = uncommitted({p: plist[p]}, version, exit)
                    if ok and not exit.is_set():
                        ok = github.publish(p, plist[p], version, exit)
                        if ok and not exit.is_set():
                            del unreleased[p]
                            save_release_info(version, state)

            if 'uhppoted' in unreleased:
                ok = build.release('uhppoted', plist['uhppoted'], version, exit)
                if ok and not exit.is_set():
                    # # ... confirm uhppoted and submodule binary checksums match
                    # print(f'     >>> verifying checksums')
                    # ignore = [ 'uhppoted', 'uhppoted-nodejs', 'node-red-contrib-uhppoted', 'uhppoted-python' ]
                    # it = itertools.filterfalse(lambda p: p in ignore, plist)
                    # for p in it:
                    #     if not build.checksum(p, plist[p], version.version(p)):
                    #         raise Exception(f"{p} 'dist' checksums differ")

                    ok = uncommitted({'uhppoted': plist['uhppoted']}, version, exit)
                    if ok and not exit.is_set():
                        ok = github.publish('uhppoted', plist['uhppoted'], version, exit)
                        if ok and not exit.is_set():
                            del unreleased['uhppoted']
                            save_release_info(version, state)

            ### TODO release nodejs, node-red and python
            ### TODO remove dist files
            ### TODO bump version
            ### TODO remove release notes

        #     if 'bump' in ops:
        #         if len(unreleased) != 0:
        #             raise Exception(f'Projects {unreleased} have not been released')

        #         for p in plist:
        #             print(f'>>>> bumping {p}')
        #             project = plist[p]
        #             clean_release_notes(p, project)
        #             bump_changelog(p, project)
        #         print()

        print()
        print(f'*** OK!')
        print()
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
