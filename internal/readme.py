import datetime
import os
import re
import subprocess
import time

sublime2 = '"/Applications/Sublime Text 2.app/Contents/SharedSupport/bin/subl"'


def READMEs(projects, versions, exit):
    print(f'>>>> checking READMEs')

    while True:
        ok = True
        for project in projects:
            versions = versions.version(project)

            if not readme(p, project, version, exit):
                ok = False

            if exit.is_set():
                return False
        if ok:
            break

    return True


def readme(project, info, version, exit):
    print(f'     ... {project}')

    path = f"{info.folder}/README.md"
    README = ''

    with open(path, 'r', encoding="utf-8") as f:
        README = f.read()

    if re.compile(f'{version}').search(README) == None:
        modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
        command = f'{sublime2} {path}'
        subprocess.run([command], shell=True)

        print(f'     ... {project} README has not been updated for release')

        while not exit.is_set():
            exit.wait(1)
            t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            if t != modified:
                break

        return False

    return True
