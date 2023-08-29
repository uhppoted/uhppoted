import datetime
import json
import os
import subprocess

sublime2 = '"/Applications/Sublime Text 2.app/Contents/SharedSupport/bin/subl"'
ignore = []


def uncommitted(projects, version, exit):
    print(f'>>>> checking for uncommitted changes ({version})')

    while True:
        ok = True
        for p in projects:
            if not p in ignore:
                project = projects[p]
                v = version.version(p)

                if not _uncommitted(p, project, v[1:], exit):
                    ok = False

                if exit.is_set():
                    return False
        if ok:
            break

    return True


def _uncommitted(project, info, version, exit):
    print(f'     ... {project}')
    try:
        command = f"cd {info['folder']} && git remote update"
        subprocess.run(command, shell=True, check=True)

        command = f"cd {info['folder']} && git status -uno"
        result = subprocess.check_output(command, shell=True)

        if (not project in ignore) and 'Changes not staged for commit' in str(result):
            raise Exception(f"{project} has uncommitted changes")

    except subprocess.CalledProcessError:
        raise Exception(f"{project}: command 'git status' failed")

    return True
