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
    print("*** uhppoted release script ***")
    print()

    if len(sys.argv) < 2:
        usage()
        return -1

    parser = argparse.ArgumentParser(description='release --version=<version> --no-edit <command>')

    parser.add_argument('command', type=str, help='command')
    parser.add_argument('--version', type=str, default='development', help='release version e.g. v0.8.4')
    parser.add_argument('--node-red', type=str, default='development', help='NodeRED release version e.g. v1.1.2')
    parser.add_argument('--prepare', action='store_true', help="executes only the 'prepare release' operation")
    parser.add_argument('--prerelease',
                        action='store_true',
                        help="executes the 'prepare and prerelease builds' operation")
    parser.add_argument('--release',
                        action='store_true',
                        help="executes only the 'prepare, prerelease and build releases' operation")
    parser.add_argument('--bump', action='store_true', help="executes only the 'post-release' operation")
    parser.add_argument('--no-edit',
                        action='store_true',
                        help="doesn't automatically invoke the editor for e.g. CHANGELOG.md'")
    parser.add_argument('--interim', action='store_false', help="doesn't insist on changes being pushed to github")

    args = parser.parse_args()

    ops = args.command.split(',')
    no_edit = args.no_edit
    interim = args.interim
    version = Version(args.version, args.node_red)

    print(f"VERSION: {version}")
    print()

    plist = projects()
    it = itertools.filterfalse(lambda p: github.already_released(p, plist[p], version.version(p)), plist)
    unreleased = {p: plist[p] for p in it}
    print()

    if 'prepare' in ops:
        CHANGELOGs(unreleased, version, exit)
        print()
        READMEs(unreleased, version, exit)
        print()
        package_versions(unreleased, version, exit)
        print()
        uncommitted(unreleased, version, exit)
        print()
        build.prepare(unreleased, version, exit)
        print()

    if 'prerelease' in ops:
        build.prerelease(unreleased, version, exit)
        print()

    if 'release' in ops:
        github.release_notes(unreleased, version, exit)
        build.release(unreleased, version, exit)
        github.publish(unreleased, version, exit)
        print()

    while not exit.is_set():
        project = ''

        try:
            if 'bump' in ops:
                if len(unreleased) != 0:
                    raise Exception(f'Projects {unreleased} have not been released')

                for p in plist:
                    print(f'>>>> bumping {p}')
                    project = plist[p]
                    clean_release_notes(p, project)
                    bump_changelog(p, project)
                break

            print()
            print(f'*** OK!')
            print()
            say('OK')
            break

        except BaseException as x:
            print(traceback.format_exc())

            msg = f'{x}'
            msg = msg.replace('uhppoted-','')                        \
                     .replace('uhppote-','')                         \
                     .replace('uhppoted','umbrella project')         \
                     .replace('cli','[[char LTRL]]cli[[char NORM]]') \
                     .replace('git','[[inpt PHON]]git[[input TEXT]]') \
                     .replace('codegen','code gen')

            print()
            print(f'*** ERROR  {x}')
            print()

            say('ERROR')
            say(msg)

            if args.prepare and not exit.is_set():
                sys.stdout.write('  waiting ')
                sys.stdout.flush()
                timestamps = {
                    'changelog': os.stat(f"{project['folder']}/CHANGELOG.md").st_mtime,
                    'readme': os.stat(f"{project['folder']}/README.md").st_mtime,
                    'git': has_uncommitted_changes(project['folder']),
                }

                for i in range(60):
                    if exit.is_set():
                        break
                    elif os.stat(f"{project['folder']}/CHANGELOG.md").st_mtime != timestamps['changelog']:
                        break
                    elif os.stat(f"{project['folder']}/README.md").st_mtime != timestamps['readme']:
                        break
                    elif has_uncommitted_changes(project['folder']) != timestamps['git']:
                        break
                    else:
                        sys.stdout.write('.')
                        sys.stdout.flush()
                        time.sleep(1)

                if os.stat(f"{project['folder']}/CHANGELOG.md").st_mtime != timestamps['changelog']:
                    print('CHANGELOG updated')
                    say('CHANGELOG updated')
                    continue
                elif os.stat(f"{project['folder']}/README.md").st_mtime != timestamps['readme']:
                    print('README updated')
                    say('README updated')
                    continue
                elif has_uncommitted_changes(project['folder']) != timestamps['git']:
                    print('git updated')
                    say('git updated')
                    continue
                else:
                    print()
                    break

        sys.exit(1)


def has_uncommitted_changes(folder):
    command = f"git -C {folder} status  -uno"
    result = subprocess.check_output(command, shell=True)

    if 'Changes not staged for commit' in str(result):
        return True

    return False


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
