import datetime
import os
import re
import subprocess
import time

editor = '"/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl"'


def READMEs(projects, versions, exit):
    print(f'>>>> checking READMEs')

    while True:
        ok = True
        for project in projects:
            version = versions.version(project.name)

            if not readme(project, version, exit):
                ok = False

            if exit.is_set():
                return False
        if ok:
            break

    return True


def readme(project, version, exit):
    print(f'     ... {project.name}')

    path = f"{project.folder}/README.md"
    README = ''

    with open(path, 'r', encoding="utf-8") as f:
        README = f.read()

    if re.compile(f'{version}').search(README) == None:
        copy_release_notes(project, version)

        modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
        command = f'{editor} {path}'
        subprocess.run([command], shell=True)

        print(f'     ... {project.name} README has not been updated for release')

        while not exit.is_set():
            exit.wait(1)
            t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            if t != modified:
                break

        return False

    return True


def copy_release_notes(project, version):
    path = f'{project.folder}/CHANGELOG.md'

    with open(path, 'r', encoding="utf-8") as f:
        changelog = f.read()
        matches = list(re.finditer(r'(?m)^## \[[^\]]+\].*?(?=^## \[|\Z)', changelog, flags=re.DOTALL | re.MULTILINE))

        if matches:
            today = datetime.date.today().isoformat()
            tomorrow = (datetime.date.today() + datetime.timedelta(days=1)).isoformat()
            heading = f'**[{version}](https://github.com/uhppoted/{project.name}/releases/tag/{version}) - {today}**'
            release_notes = matches[0].group().strip() + '\n'
            clip = tag + '\n\n' + release_notes + '\n'
            subprocess.run("pbcopy", text=True, input=clip)
        else:
            subprocess.run("pbcopy", text=True, input='')
