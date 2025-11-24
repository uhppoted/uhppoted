import datetime
import json
import tomllib
import os
import re
import subprocess
import xml.etree.ElementTree as ET

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

                if project.packaging:
                    if project.packaging == 'go':
                        if not go_package_version(p, project, v[1:], exit):
                            ok = False
                    elif project.packaging == 'javascript':
                        if not javascript_package_version(p, project, v[1:], exit):
                            ok = False
                    elif project.packaging == 'python':
                        if not python_package_version(p, project, v[1:], exit):
                            ok = False
                    elif project.packaging == 'dotnet':
                        if not dotnet_package_version(p, project, v[1:], exit):
                            ok = False

                if exit.is_set():
                    return False
        if ok:
            break

    return True


# NTS: not tested and debugged
def go_package_version(project, info, version, exit):
    print(f'     ... {project}')

    if project != 'uhppoted-app-wild-apricot':
        return True

    path = f"{info.folder}/cmd/uhppoted-app-wild-apricot/main.go"

    if version != 'development' and os.path.isfile(f'{path}'):
        with open(path, 'r', encoding="utf-8") as f:
            for line in f:
                if match := re.search('^const VERSION string = "(v[0-9]+.[0-9]+.[0-9]+)"', line):
                    print(f'>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DEBUG {line}')
                    print(f'>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DEBUG {match}')
                    break

    #     if package['version'] != version:
    #         modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
    #         command = f'{sublime2} {path}'
    #         subprocess.run([command], shell=True)
    #
    #         print(f"     ... package version:{package['version']} - expected {version}")
    #
    #         while not exit.is_set():
    #             exit.wait(1)
    #             t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
    #             if t != modified:
    #                 break
    #
    #         return False
    #
    # return True
    return False


def javascript_package_version(project, info, version, exit):
    print(f'     ... {project}')

    path = f"{info.folder}/package.json"

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


def python_package_version(project, info, version, exit):
    print(f'     ... {project}')

    path = f"{info.folder}/pyproject.toml"

    if version != 'development' and os.path.isfile(f'{path}'):
        with open(path, 'rb') as f:
            package = tomllib.load(f)

        if package['project']['version'] != version:
            modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            command = f'{sublime2} {path}'
            subprocess.run([command], shell=True)

            print(f"     ... package version:{package['project']['version']} - expected {version}")

            while not exit.is_set():
                exit.wait(1)
                t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
                if t != modified:
                    break

            return False

    return True


def dotnet_package_version(project, info, version, exit):
    print(f'     ... {project}')

    path = f"{info.folder}/uhppoted/uhppoted/uhppoted.fsproj"

    if version != 'development' and os.path.isfile(f'{path}'):
        namespace = {'msbuild': 'http://schemas.microsoft.com/developer/msbuild/2003'}
        tree = ET.parse(path)
        root = tree.getroot()

        PackageVersion = root.find('PropertyGroup/PackageVersion', namespace)
        Version = root.find('PropertyGroup/Version', namespace)

        if PackageVersion is None or PackageVersion.text != version or Version is None or Version.text != version:
            modified = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
            command = f'{sublime2} {path}'
            subprocess.run([command], shell=True)

            if PackageVersion is None:
                print(f"     ... PackageVersion missing")
            elif PackageVersion.text != version:
                print(f"     ... PackageVersion:{PackageVersion.text} (expected {version})")

            if Version is None:
                print(f"     ... Version missing")
            elif Version.text != version:
                print(f"     ... Version:{Version.text} (expected {version})")

            while not exit.is_set():
                exit.wait(1)
                t = datetime.datetime.fromtimestamp(os.path.getmtime(path)).strftime('%Y-%m-%d %H:%M:%S')
                if t != modified:
                    break

            return False

    return True
