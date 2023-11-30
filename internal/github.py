import itertools
import subprocess
import re

from build import checksum
from misc import say
from git import _uncommitted

ignore = []


def already_released(project, info, version):
    command = f"cd {info['folder']} && git fetch --tags"
    result = subprocess.call(command, shell=True)
    if result != 0:
        raise Exception(f"command 'git fetch --tags' failed")
    else:
        command = f"cd {info['folder']} && git tag {version} --list"
        result = subprocess.check_output(command, shell=True)

        if f'{version}' in str(result):
            print(f'     +++ {project} has been released')
            return True
        else:
            print(f'     ... {project} has not been released')
            return False


def release_notes(projects, version, exit):
    print(f'>>>> generating release notes ({version})')

    while True:
        ok = True
        for p in projects:
            if not p in ignore:
                project = projects[p]
                v = version.version(p)

                if not _release_notes(p, project, v[1:], exit):
                    ok = False

                if exit.is_set():
                    return False
        if ok:
            break

    return True


def _release_notes(project, info, version, exit):
    print(f'     ... {project}')

    regex = r'##\s+\[(.*?)\](?:.*?)\n(.*?)##\s+\[(.*?)\]'
    path = f"{info['folder']}/release-notes.md"
    changelog = f"{info['folder']}/CHANGELOG.md"

    try:
        with open(changelog, 'r', encoding="utf-8") as f:
            CHANGELOG = f.read()

        with open(path, 'xt', encoding="utf-8") as f:
            match = re.search(regex, CHANGELOG, re.MULTILINE | re.DOTALL)

            current = match.group(1)
            notes = match.group(2).strip()
            previous = match.group(3)

            if notes == '':
                notes = 'Maintenance release for version compatibility.'

            f.write('### Release Notes\n')
            f.write('\n')
            f.write(notes)
            f.write('\n')

    except FileExistsError:
        f"... keeping existing {info['folder']}/release-notes.md"

    return True


def publish(project, p, version, exit):
    print(f'>>>> publishing {project} ({version})')

    # ... confirm uhppoted and submodule binary checksums match
    print(f'     >>> verifying checksums')
    if not checksum(project, p, version.version(p)):
        raise Exception(f"{project} 'dist' checksums differ")

    # ... confirm no uncommitted changes
    print(f'     >>> checking for uncommitted changes')
    if not _uncommitted(project, p, version, exit):
        raise Exception(f"{project} has uncommitted changes")

    # ... confirm changes have been pushed to github repo
    if not pushed(project, p):
        print(f'     ... {project} has unpushed changes')
        say(f'{project} has unpushed changes')
        exit.wait(10)
        while not pushed(project, p) and not exit.is_set():
            exit.wait(10)

        if exit.is_set():
            return False

    # ... publish release
    _publish(project, p, version.version(p))
    print(f'     ... {project} is waiting for release on github')
    say(f'{project} is waiting for release on github')

    exit.wait(10)
    while not published(project, p, version.version(p)) and not exit.is_set():
        exit.wait(10)

    if exit.is_set():
        return False

    return True


# def publish(projects, version, exit):
#     print(f'>>>> publishing releases ({version})')
#
#     plist = {p: projects[p] for p in itertools.filterfalse(lambda p: p in ignore, projects)}
#
#     while True:
#         ok = True
#
#         # ... confirm uhppoted and submodule binary checksums match
#         print(f'     >>> verifying checksums')
#         for p in plist:
#             print(f'     ... {p}')
#             if not checksum(p, plist[p], version.version(p)):
#                 raise Exception(f"{project} 'dist' checksums differ")
#
#             if exit.is_set():
#                 return False
#
#         for p in plist:
#             # ... remote updated?
#             print(f'     ... {p}')
#
#             if not pushed(p, plist[p]):
#                 print(f'     ... {p} has unpushed changes')
#                 say(f'{p} has unpushed changes')
#                 ok = False
#                 exit.wait(10)
#                 while not pushed(p, plist[p]) and not exit.is_set():
#                     exit.wait(10)
#
#             if exit.is_set():
#                 return False
#
#             _publish(p, plist[p], version.version(p))
#             print(f'     ... {p} is waiting for release on github')
#             say(f'{p} is waiting for release on github')
#             exit.wait(10)
#             while not published(p, plist[p], version.version(p)) and not exit.is_set():
#                 exit.wait(10)
#
#             if exit.is_set():
#                 return False
#
#         if ok:
#             return True
#
#     return True


def pushed(project, info):
    try:
        command = f"cd {info['folder']} && git remote update"
        subprocess.run(command, shell=True, check=True)

        command = f"cd {info['folder']} && git status -uno"
        result = subprocess.check_output(command, shell=True)

        if "Your branch is up to date with 'origin/" in str(result):
            return True

    except subprocess.CalledProcessError:
        raise Exception(f"{project}: command 'git status' failed")

    return False


def _publish(project, info, version):
    try:
        folder = info['folder']
        command = f"cd {folder} && make publish DIST={project}_{version} VERSION={version}"
        subprocess.run(command, shell=True, check=True)
    except subprocess.CalledProcessError:
        raise Exception(f"command 'update {project}' failed")

    return True


def published(project, info, version):
    command = f"cd {info['folder']} && git fetch --tags"
    result = subprocess.call(command, shell=True)
    if result != 0:
        raise Exception(f"command 'git fetch --tags' failed")
    else:
        command = f"cd {info['folder']} && git tag {version} --list"
        result = subprocess.check_output(command, shell=True)

        if f'{version}' in str(result):
            print(f'     +++ {project} has been released')
            return True
        else:
            print(f'     ... {project} has not been released')
            return False
