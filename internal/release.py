#!python3

import argparse
import subprocess
import sys
import os
import shutil
import json
import re
import hashlib
import signal
import time
import datetime
import itertools
import traceback
import tempfile
import build

from threading import Event

from projects import projects
from version import Version
import github
import npm
from changelog import CHANGELOGs
from readme import READMEs
from packaging import package_versions
from git import uncommitted
from misc import say

sublime2 = '"/Applications/Sublime Text 2.app/Contents/SharedSupport/bin/subl"'
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
    parser.add_argument('--version', type=str, default='development', help='release version e.g. v0.8.10')
    parser.add_argument('--node-red', type=str, default='development', help='NodeRED release version e.g. v1.1.9')
    parser.add_argument('--release', action='store_true', help="builds the release versions")
    parser.add_argument('--bump', action='store_true', help="bumps version and cleans up after a release")

    args = parser.parse_args()
    versions = Version(args.version, args.node_red)

    print(f'VERSION: {versions}')
    print()

    try:
        # ... get release state
        state = get_release_info(versions)
        plist = projects()

        # ... get unreleased project list
        if not 'unreleased' in state:
            it = itertools.filterfalse(lambda p: github.already_released(p, plist[p], versions.version(p)), plist)
            state['unreleased'] = [p for p in it]
            save_release_info(versions, state)
            print()

        unreleased = {p: plist[p] for p in state['unreleased']}

        # ... CHANGELOG.md
        if not 'changelogs' in state or state['changelogs'] != 'ok':
            CHANGELOGs(unreleased.values(), versions, exit)
            print()
            if exit.is_set():
                return -1
            else:
                state['changelogs'] = 'ok'
                save_release_info(versions, state)

        # ... README.md
        if not 'readmes' in state or state['readmes'] != 'ok':
            READMEs(unreleased.values(), versions, exit)
            print()
            if exit.is_set():
                return -1
            else:
                state['readmes'] = 'ok'
                save_release_info(versions, state)

        # ... package versions
        if not 'package-versions' in state or state['package-versions'] != 'ok':
            package_versions(unreleased, versions, exit)
            print()
            if exit.is_set():
                return -1
            else:
                state['package-versions'] = 'ok'
                save_release_info(versions, state)

        # ... uncommitted changes
        if not 'uncommitted-changes' in state or state['uncommitted-changes'] != 'ok':
            ok = uncommitted(unreleased, versions, exit)
            print()
            if exit.is_set():
                return -1
            elif ok:
                state['uncommitted-changes'] = 'ok'
                save_release_info(versions, state)

        # ... 'prepare' builds
        if not 'prepared' in state or state['prepared'] != 'ok':
            build.prepare(unreleased, versions, exit)
            state['uncommitted-changes'] = '??'
            save_release_info(versions, state)
            ok = uncommitted(unreleased, versions, exit)
            print()
            if exit.is_set():
                return -1
            elif ok:
                state['prepared'] = 'ok'
                state['uncommitted-changes'] = 'ok'
                save_release_info(versions, state)

        # ... build release versions and publish
        if args.release:
            if not 'release-notes' in state or state['release-notes'] != 'ok':
                ok = github.release_notes(unreleased, versions, exit)
                print()
                if exit.is_set():
                    return -1
                elif ok:
                    state['release-notes'] = 'ok'
                    save_release_info(versions, state)

            if 'uhppote-core' in unreleased:
                ok = build.release('uhppote-core', plist['uhppote-core'], versions, exit)
                if exit.is_set():
                    return -1
                elif ok:
                    ok = github.publish('uhppote-core', plist['uhppote-core'], versions, exit)
                    if exit.is_set():
                        return -1
                    elif ok:
                        state['unreleased'] = [v for v in state['unreleased'] if v != 'uhppote-core']
                        del unreleased['uhppote-core']
                        save_release_info(versions, state)

            if 'uhppoted-lib' in unreleased:
                ok = build.release('uhppoted-lib', plist['uhppoted-lib'], versions, exit)
                if exit.is_set():
                    return -1
                elif ok:
                    ok = github.publish('uhppoted-lib', plist['uhppoted-lib'], versions, exit)
                    if exit.is_set():
                        return -1
                    elif ok:
                        state['unreleased'] = [v for v in state['unreleased'] if v != 'uhppoted-lib']
                        del unreleased['uhppoted-lib']
                        save_release_info(versions, state)

            if 'uhppoted-lib-python' in unreleased:
                ok = build.release('uhppoted-lib-python', plist['uhppoted-lib-python'], versions, exit)
                if exit.is_set():
                    return -1
                elif ok:
                    ok = github.publish('uhppoted-lib-python', plist['uhppoted-lib-python'], versions, exit)
                    if exit.is_set():
                        return -1
                    elif ok:
                        # FIXME publish to testpy
                        # FIXME publish to pypi
                        state['unreleased'] = [v for v in state['unreleased'] if v != 'uhppoted-lib-python']
                        del unreleased['uhppoted-lib-python']
                        save_release_info(versions, state)

            if 'uhppoted-lib-nodejs' in unreleased:
                ok = build.release('uhppoted-lib-nodejs', plist['uhppoted-lib-nodejs'], versions, exit)
                if exit.is_set():
                    return -1
                elif ok:
                    ok = github.publish('uhppoted-lib-nodejs', plist['uhppoted-lib-nodejs'], versions, exit)
                    if exit.is_set():
                        return -1
                    elif ok:
                        ok = npm.publish('uhppoted-lib-nodejs', plist['uhppoted-lib-nodejs'], versions, exit)
                        if exit.is_set():
                            return -1
                        elif ok:
                            state['unreleased'] = [v for v in state['unreleased'] if v != 'uhppoted-lib-nodejs']
                            del unreleased['uhppoted-lib-nodejs']
                            save_release_info(versions, state)

            if 'node-red-contrib-uhppoted' in unreleased:
                ok = build.release('node-red-contrib-uhppoted', plist['node-red-contrib-uhppoted'], versions, exit)
                if exit.is_set():
                    return -1
                elif ok:
                    ok = github.publish('node-red-contrib-uhppoted', plist['node-red-contrib-uhppoted'], versions, exit)
                    if exit.is_set():
                        return -1
                    elif ok:
                        ok = npm.publish('node-red-contrib-uhppoted', plist['node-red-contrib-uhppoted'], versions,
                                         exit)
                        if exit.is_set():
                            return -1
                        elif ok:
                            state['unreleased'] = [v for v in state['unreleased'] if v != 'node-red-contrib-uhppoted']
                            del unreleased['node-red-contrib-uhppoted']
                            save_release_info(versions, state)

            ignore = ['uhppoted', 'uhppoted-lib-nodejs', 'node-red-contrib-uhppoted', 'uhppoted-lib-python']

            it = itertools.filterfalse(lambda p: p in ignore, unreleased)
            rl = {p: plist[p] for p in it}
            for p in rl:
                ok = build.release(p, plist[p], versions, exit)
                if exit.is_set():
                    return -1
                elif ok:
                    ok = uncommitted({p: plist[p]}, versions, exit)
                    if exit.is_set():
                        return -1
                    elif ok:
                        ok = github.publish(p, plist[p], versions, exit)
                        if exit.is_set():
                            return -1
                        elif ok:
                            state['unreleased'] = [v for v in state['unreleased'] if v != p]
                            del unreleased[p]
                            save_release_info(versions, state)

            if 'uhppoted' in unreleased:
                ok = build.release('uhppoted', plist['uhppoted'], versions, exit)
                if exit.is_set():
                    return -1
                elif ok:
                    # ... confirm uhppoted and submodule binary checksums match
                    # print(f'     >>> verifying checksums')
                    # ignore = ['uhppoted', 'uhppoted-lib-nodejs', 'node-red-contrib-uhppoted', 'uhppoted-lib-python']
                    # it = itertools.filterfalse(lambda p: p in ignore, plist)
                    # for p in it:
                    #     if not build.checksum(p, plist[p], versions.version(p)):
                    #         raise Exception(f"{p} 'dist' checksums differ")

                    ok = uncommitted({'uhppoted': plist['uhppoted']}, versions, exit)
                    if exit.is_set():
                        return -1
                    elif ok:
                        ok = github.publish('uhppoted', plist['uhppoted'], versions, exit)
                        if exit.is_set():
                            return -1
                        elif ok:
                            state['unreleased'] = [v for v in state['unreleased'] if v != 'uhppoted']
                            del unreleased['uhppoted']
                            save_release_info(versions, state)

        # ... post-release cleanup
        if args.bump:
            if len(unreleased) != 0:
                raise Exception(f'Projects {unreleased} have not been released')

            for p in plist:
                project = plist[p]
                clean_release_notes(p, project)
                bump_changelog(p, project)
                clean_dist(p, project)
                bump_version(p, project, versions, exit)
                print(f'>>>> bumped {p}')

            print()

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
        json.dump(info, f, skipkeys=True, ensure_ascii=False, indent=4)


def clean_release_notes(project, info):
    file = f"{info.folder}/release-notes.md"
    if os.path.isfile(file):
        os.remove(file)
        print(f'     ... {project} removed release-notes.md')


def bump_changelog(project, info):
    with open(f"{info.folder}/CHANGELOG.md", 'r', encoding="utf-8") as f:
        CHANGELOG = f.read()
        if 'Unreleased' in CHANGELOG:
            return

    tmpfile = tempfile.NamedTemporaryFile(mode="w+t", delete=False)

    try:
        with open(f"{info.folder}/CHANGELOG.md", 'r', encoding="utf-8") as f:
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

            os.rename(f"{info.folder}/CHANGELOG.md", f"{info.folder}/CHANGELOG.bak")
            os.rename(tmpfile.name, f"{info.folder}/CHANGELOG.md")
            if os.path.isfile(f"{info.folder}/CHANGELOG.bak"):
                os.remove(f"{info.folder}/CHANGELOG.bak")

            print(f'     ... {project} updated CHANGELOG for next dev cycle')
    finally:
        tmpfile.close()
        if os.path.isfile(tmpfile.name):
            os.remove(tmpfile.name)


def clean_dist(project, info):
    folder = f'{info.folder}/dist'

    if os.path.exists(folder) and os.listdir(folder):
        for filename in os.listdir(folder):
            file_path = os.path.join(folder, filename)

            if os.path.isfile(file_path) or os.path.islink(file_path):
                os.unlink(file_path)
            elif os.path.isdir(file_path):
                shutil.rmtree(file_path)

        print(f"     ... {project} cleaned 'dist' folder")


def bump_version(project, info, version, exit):
    if project == 'uhppote-core':
        path = f'{info.folder}/uhppote/uhppote.go'
        modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
        command = f'{sublime2} {path}'
        subprocess.run([command], shell=True)

        print(f'     ... {project} bump VERSION')
        while not exit.is_set():
            exit.wait(1)
            t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            if t != modified:
                break

    if project == 'uhppoted-httpd':
        path = f'{info.folder}/package.json'
        modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
        command = f'{sublime2} {path}'
        subprocess.run([command], shell=True)

        print(f'     ... {project} bump package json')
        while not exit.is_set():
            exit.wait(1)
            t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            if t != modified:
                break

    if project == 'uhppoted-lib-nodejs':
        path = f'{info.folder}/package.json'
        modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
        command = f'{sublime2} {path}'
        subprocess.run([command], shell=True)

        print(f'     ... {project} bump package json')
        while not exit.is_set():
            exit.wait(1)
            t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            if t != modified:
                break

    if project == 'node-red-contrib-uhppoted':
        path = f'{info.folder}/package.json'
        modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
        command = f'{sublime2} {path}'
        subprocess.run([command], shell=True)

        print(f'     ... {project} bump package json')
        while not exit.is_set():
            exit.wait(1)
            t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            if t != modified:
                break

    if project == 'uhppoted-lib-python':
        path = f'{info.folder}/pyproject.toml'
        modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
        command = f'{sublime2} {path}'
        subprocess.run([command], shell=True)

        print(f'     ... {project} bump package json')
        while not exit.is_set():
            exit.wait(1)
            t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            if t != modified:
                break

    if project == 'uhppoted':
        path = f'{info.folder}/Makefile'
        modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
        command = f'{sublime2} {path}'
        subprocess.run([command], shell=True)

        print(f'     ... {project} bump Makefile versions')
        while not exit.is_set():
            exit.wait(1)
            t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            if t != modified:
                break


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
