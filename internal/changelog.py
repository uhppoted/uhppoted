import datetime
import os
import signal
import subprocess
import time

sublime2 = '"/Applications/Sublime Text 2.app/Contents/SharedSupport/bin/subl"'


def changelogs(projects, version, nodered, exit):
    print(f'>>>> checking CHANGELOGs ({version})')

    ok = True

    for p in projects:
        project = projects[p]
        v = nodered if p == 'node-red-contrib-uhppoted' else version

        if not changelog(p, project, v, exit):
            print(f'{p} CHANGELOG has not been updated for release')
            ok = False

        if exit.is_set():
            return False

    return ok


def changelog(project, info, version, exit):
    print(f'>>>> checking {project}')

    path = f"{info['folder']}/CHANGELOG.md"
    CHANGELOG = ''

    with open(path, 'r', encoding="utf-8") as f:
        CHANGELOG = f.read()

    if 'Unreleased' in CHANGELOG:
        rest = CHANGELOG
        for i in range(3):
            line, _, rest = rest.partition('\n')
            print(f'>> {line}')

        modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
        command = f"{sublime2} {info['folder']}/CHANGELOG.md"
        subprocess.run([command], shell=True)

        print(f'   ... please update {project} CHANGELOG for release')

        while not exit.is_set():
            exit.wait(1)
            t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            if t != modified:
                break

        with open(path, 'r', encoding="utf-8") as f:
            CHANGELOG = f.read()

        if 'Unreleased' in CHANGELOG:
            return False

    with open(path, 'r', encoding="utf-8") as f:
        CHANGELOG = f.read()

    if not CHANGELOG.startswith(f'# CHANGELOG\n\n## [{version}]'):
        rest = CHANGELOG
        for i in range(3):
            line, _, rest = rest.partition('\n')
            print(f'>> {line}')

        command = f"{sublime2} {info['folder']}/CHANGELOG.md"
        subprocess.run([command], shell=True)

        print(f'   ... please fix {project} CHANGELOG version for release {version}')

        while not exit.is_set():
            exit.wait(1)
            t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            if t != modified:
                break

        with open(path, 'r', encoding="utf-8") as f:
            CHANGELOG = f.read()

        if not CHANGELOG.startswith(f'# CHANGELOG\n\n## [{version}]'):
            raise Exception(f'{project} CHANGELOG has not been updated for release')

    return True
