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

exit = Event()


def quit(signo, _frame):
    print("Interrupted by %d, shutting down" % signo)
    exit.set()


def main():
    print()
    print("*** uhppoted release-all")
    print()

    if len(sys.argv) < 2:
        usage()
        return -1

    parser = argparse.ArgumentParser(description='release --version=<version> --no-edit')

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
    no_edit = args.no_edit
    interim = args.interim
    version = args.version
    nodered = args.node_red

    if version != 'development' and not args.version.startswith('v'):
        version = f'v{args.version}'

    print(f'VERSION: {version}')

    plist = projects()
    it = itertools.filterfalse(
        lambda p: already_released(p, plist[p], version if p != 'node-red-contrib-uhppoted' else f'v{nodered}'), plist)
    unreleased = {p: plist[p] for p in it}

    while not exit.is_set():
        project = ''

        try:
            if args.bump:
                if len(unreleased) != 0:
                    raise Exception(f'Projects {unreleased} have not been released')

                for p in plist:
                    print(f'>>>> bumping {p}')
                    project = plist[p]
                    clean_release_notes(p, project)
                    bump_changelog(p, project)
                break

            print(f'>>>> initialise: checking CHANGELOGs, READMEs and uncommitted changes ({version})')
            for p in unreleased:
                print(f'>>>> checking {p}')
                project = unreleased[p]

                if p == 'node-red-contrib-uhppoted':
                    package_version(p, project, nodered, no_edit)
                else:
                    package_version(p, project, version[1:], no_edit)

                changelog(p, project, version[1:], no_edit)
                readme(p, project, version, no_edit)
                uncommitted(p, project, interim)

            if args.prepare or args.prerelease or args.release:
                print(f'>>>> prepare: checking builds ({version})')
                for p in unreleased:
                    print(f'... releasing {p}')
                    update(p, unreleased[p])
                    checkout(p, unreleased[p])
                    build(p, unreleased[p])

            if args.prerelease or args.release:
                print(f'>>>> prerelease: final check for consistent library and binary versions ({version})')
                for p in unreleased:
                    checksum(p, unreleased[p], 'development')
                    git(p, unreleased[p], interim)

            if args.release:
                print(f'>>>> release: building release versions ({version})')
                for p in unreleased:
                    release_notes(p, unreleased[p], version)
                    release(p, unreleased[p], version)
                    git(p, unreleased[p], interim)

            if args.release:
                print(f'>>>> release: verify checksums of release versions ({version})')
                for p in unreleased:
                    checksum(p, unreleased[p], version)

            # TODO publish

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


def projects():
    return {
        'uhppote-core': {
            'folder': './uhppote-core',
            'branch': 'master'
        },
        'uhppoted-lib': {
            'folder': './uhppoted-lib',
            'branch': 'master',
        },
        'uhppote-simulator': {
            'folder': './uhppote-simulator',
            'branch': 'main',
            'binary': 'uhppote-simulator'
        },
        'uhppote-cli': {
            'folder': './uhppote-cli',
            'branch': 'master',
            'binary': 'uhppote-cli'
        },
        'uhppoted-rest': {
            'folder': './uhppoted-rest',
            'branch': 'main',
            'binary': 'uhppoted-rest'
        },
        'uhppoted-mqtt': {
            'folder': './uhppoted-mqtt',
            'branch': 'master',
            'binary': 'uhppoted-mqtt'
        },
        'uhppoted-httpd': {
            'folder': './uhppoted-httpd',
            'branch': 'master',
            'binary': 'uhppoted-httpd'
        },
        'uhppoted-tunnel': {
            'folder': './uhppoted-tunnel',
            'branch': 'main',
            'binary': 'uhppoted-tunnel'
        },
        'uhppoted-dll': {
            'folder': './uhppoted-dll',
            'branch': 'main',
        },
        'uhppoted-codegen': {
            'folder': './uhppoted-codegen',
            'branch': 'main',
            'binary': 'uhppoted-codegen'
        },
        'uhppoted-app-s3': {
            'folder': './uhppoted-app-s3',
            'branch': 'main',
            'binary': 'uhppoted-app-s3'
        },
        'uhppoted-app-sheets': {
            'folder': './uhppoted-app-sheets',
            'branch': 'main',
            'binary': 'uhppoted-app-sheets'
        },
        'uhppoted-app-wild-apricot': {
            'folder': './uhppoted-app-wild-apricot',
            'branch': 'main',
            'binary': 'uhppoted-app-wild-apricot'
        },
        'uhppoted-app-db': {
            'folder': './uhppoted-app-db',
            'branch': 'main',
            'binary': 'uhppoted-app-db'
        },
        'uhppoted-nodejs': {
            'folder': './uhppoted-nodejs',
            'branch': 'main',
        },
        'node-red-contrib-uhppoted': {
            'folder': './node-red-contrib-uhppoted',
            'branch': 'main',
        },
        'uhppoted-wiegand': {
            'folder': './uhppoted-wiegand',
            'branch': 'main',
        },
        'uhppoted': {
            'folder': '.',
            'branch': 'master'
        }
    }


def package_version(project, info, version, no_edit):
    if version != 'development' and os.path.isfile(f"{info['folder']}/package.json"):
        with open(f"{info['folder']}/package.json", 'r', encoding="utf-8") as f:
            package = json.load(f)
            if package['version'] != version:
                print(f">> package version:{package['version']} - expected {version}")
                if not no_edit:
                    command = f"sublime2 {info['folder']}/package.json"
                    subprocess.run(['/bin/zsh', '-i', '-c', command])

                raise Exception(f'{project} package.json has not been updated for release')


def changelog(project, info, version, no_edit):
    with open(f"{info['folder']}/CHANGELOG.md", 'r', encoding="utf-8") as f:
        CHANGELOG = f.read()
        if 'Unreleased' in CHANGELOG:
            rest = CHANGELOG
            for i in range(3):
                line, _, rest = rest.partition('\n')
                print(f'>> {line}')

            if not no_edit:
                command = f"sublime2 {info['folder']}/CHANGELOG.md"
                subprocess.run(['/bin/zsh', '-i', '-c', command])

            raise Exception(f'{project} CHANGELOG has not been updated for release')

    if project == 'node-red-contrib-uhppoted':
        return

    with open(f"{info['folder']}/CHANGELOG.md", 'r', encoding="utf-8") as f:
        CHANGELOG = f.read()
        if not CHANGELOG.startswith(f'# CHANGELOG\n\n## [{version}]'):
            rest = CHANGELOG
            for i in range(3):
                line, _, rest = rest.partition('\n')
                print(f'>> {line}')

            if not no_edit:
                command = f"sublime2 {info['folder']}/CHANGELOG.md"
                subprocess.run(['/bin/zsh', '-i', '-c', command])

            raise Exception(f'{project} CHANGELOG has not been updated for release')


def readme(project, info, version, no_edit):
    ignore = ['uhppoted-nodejs', 'node-red-contrib-uhppoted']
    if project in ignore:
        return

    with open(f"{info['folder']}/README.md", 'r', encoding="utf-8") as f:
        README = f.read()
        if re.compile(f'\|\s*{version}\s*\|').search(README) == None:
            if not no_edit:
                command = f"sublime2 {info['folder']}/README.md"
                subprocess.run(['/bin/zsh', '-i', '-c', command])

            raise Exception(f'{project} README has not been updated for release')


def uncommitted(project, info, interim):
    ignore = []
    if interim:
        ignore = ['uhppoted']

    try:
        command = f"cd {info['folder']} && git remote update"
        subprocess.run(command, shell=True, check=True)

        command = f"cd {info['folder']} && git status -uno"
        result = subprocess.check_output(command, shell=True)

        if (not project in ignore) and 'Changes not staged for commit' in str(result):
            raise Exception(f"{project} has uncommitted changes")

    except subprocess.CalledProcessError:
        raise Exception(f"{project}: command 'git status' failed")


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


def build(project, info):
    command = f"cd {info['folder']} && make update-release && make build-all"
    result = subprocess.call(command, shell=True)
    if result != 0:
        raise Exception(f"command 'build {project}' failed")


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


def already_released(project, info, version):
    command = f"cd {info['folder']} && git fetch --tags"
    result = subprocess.call(command, shell=True)
    if result != 0:
        raise Exception(f"command 'git fetch --tags' failed")
    else:
        command = f"cd {info['folder']} && git tag {version} --list"
        result = subprocess.check_output(command, shell=True)

        if f'{version}' in str(result):
            print(f'     ... {project} has been released')
            return True
        else:
            print(f'     ... {project} has not been released')
            return False


def release_notes(project, info, version):
    regex = r'##\s+\[(.*?)\](?:.*?)\n(.*?)##\s+\[(.*?)\]'
    file = f"{info['folder']}/release-notes.md"

    try:
        with open(file, 'xt', encoding="utf-8") as f:
            with open(f"{info['folder']}/CHANGELOG.md", 'r', encoding="utf-8") as changelog:
                CHANGELOG = changelog.read()
                match = re.search(regex, CHANGELOG, re.MULTILINE | re.DOTALL)

                current = match.group(1)
                notes = match.group(2).strip()
                previous = match.group(3)

                if notes == '':
                    notes = 'Maintenance release for version compatibility.'

                # print(f'Current  release {current}')
                # print(f'Previous release {previous}')
                # print(f'Release notes\n----\n{notes}\n----\n')

                f.write('### Release Notes\n')
                f.write('\n')
                f.write(notes)
                f.write('\n')
    except FileExistsError:
        f"... keeping existing {info['folder']}/release-notes.md"


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


def say(msg):
    transliterated = msg.replace('nodejs','node js') \
                        .replace('codegen', 'code gen') \
                        .replace('Errno','error number') \
                        .replace('exe','e x e')
    subprocess.call(f'say {transliterated}', shell=True)


def usage():
    print()
    print('  Usage: python release.py <options>')
    print()
    print('  Options:')
    print('    --version <version>  Release version e.g. v0.8.1')
    print()


if __name__ == '__main__':
    for sig in ('TERM', 'HUP', 'INT'):
        signal.signal(getattr(signal, 'SIG' + sig), quit)

    main()
