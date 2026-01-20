import datetime
import os
import signal
import subprocess
import time

editor = '"/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl"'


def CHANGELOGs(projects, versions, exit):
    print(f'>>>> checking CHANGELOGs')

    while True:
        ok = True
        for project in projects:
            version = versions.version(project)

            if not changelog(project.name, project, version[1:], exit):
                ok = False

            if exit.is_set():
                return False
        if ok:
            break

    return True


def changelog(project, info, version, exit):
    print(f'     ... {project}')

    path = f"{info.folder}/CHANGELOG.md"
    CHANGELOG = ''

    with open(path, 'r', encoding="utf-8") as f:
        CHANGELOG = f.read()

    if 'Unreleased' in CHANGELOG:
        rest = CHANGELOG
        for i in range(3):
            line, _, rest = rest.partition('\n')
            print(f'>> {line}')

        modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
        command = f'{editor} {path}'
        subprocess.run([command], shell=True)

        print(f'     ... {project} CHANGELOG has not been updated for release')

        while not exit.is_set():
            exit.wait(1)
            t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            if t != modified:
                break

        return False

    with open(path, 'r', encoding="utf-8") as f:
        CHANGELOG = f.read()

    if not CHANGELOG.startswith(f'# CHANGELOG\n\n## [{version}]'):
        rest = CHANGELOG
        for i in range(3):
            line, _, rest = rest.partition('\n')
            print(f'>> {line}')

        modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
        command = f'{editor} {path}'
        subprocess.run([command], shell=True)

        print(f'     ... {project} CHANGELOG has not been updated for release')

        while not exit.is_set():
            exit.wait(1)
            t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            if t != modified:
                break

        return False

    return True
