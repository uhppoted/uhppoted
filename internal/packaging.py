import datetime
import json
import os
import subprocess

sublime2 = '"/Applications/Sublime Text 2.app/Contents/SharedSupport/bin/subl"'
ignore = []


def package_versions(projects, version, exit):
    print(f'>>>> checking package versions ({version})')

    while True:
        ok = True
        for p in projects:
            if not p in ignore:
                project = projects[p]
                v = version.version(p)

                if 'packaging' in project:
                    if project['packaging'] == 'javascript':
                        if not javascript_package_version(p, project, v[1:], exit):
                            ok = False
                    elif project['packaging'] == 'python':
                        return False
                    elif project['packaging'] == 'dotnet':
                        return False

                if exit.is_set():
                    return False
        if ok:
            break

    return True


def javascript_package_version(project, info, version, exit):
    print(f'     ... {project}')

    path = f"{info['folder']}/package.json"

    if version != 'development' and os.path.isfile(f'{path}'):
        with open(path, 'r', encoding="utf-8") as f:
            package = json.load(f)

        if package['version'] != version:
            modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            command = f'{sublime2} {path}'
            subprocess.run([command], shell=True)

            print(f"     ... package version:{package['version']} - expected {version}")

            while not exit.is_set():
                exit.wait(1)
                t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
                if t != modified:
                    break

            return False

    return True
