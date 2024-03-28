import hashlib
import itertools
import os
import re
import subprocess

from misc import say

sublime2 = '"/Applications/Sublime Text 2.app/Contents/SharedSupport/bin/subl"'
ignore = []


def prepare(projects, version, exit):
    print(f'>>>> rebuilding all projects ({version})')

    it = itertools.filterfalse(lambda p: p in ignore, projects)
    plist = {p: projects[p] for p in it}

    while True:
        ok = True

        for p in plist:
            print(f'     ... {p}')
            if not checkout(p, plist[p]):
                ok = False
            elif not update(p, plist[p]):
                ok = False
            elif not build(p, plist[p]):
                ok = False

            if exit.is_set():
                return False

        if ok:
            return True


# def prerelease(projects, version, exit):
#     print('>>>> prerelease builds (v{version})')
#
#     it = itertools.filterfalse(lambda p: p in ignore, projects)
#     plist = {p: projects[p] for p in it}
#
#     while True:
#         ok = True
#
#         # # ... update and build
#         # print(f'     ... rebuilding all projects')
#         # for p in plist:
#         #     print(f'     ... {p}')
#         #     if not update(p, plist[p]):
#         #         ok = False
#
#         #     if exit.is_set():
#         #         return False
#
#         # for p in plist:
#         #     print(f'     ... {p}')
#         #     if not build(p, plist[p]):
#         #         ok = False
#
#         #     if exit.is_set():
#         #         return False
#
#         # # ... uncommitted changes?
#         # print('     ... checking for uncommitted changes')
#         # for p in plist:
#         #     print(f'     ... {p}')
#         #     uncommitted(p, plist[p])
#
#         #     if exit.is_set():
#         #         return False
#
#         # # ... checkout and rebuild all
#         # print('     ... checking out latest github version')
#         # for p in plist:
#         #     print(f'     ... {p}')
#         #     checkout(p, plist[p])
#
#         #     if exit.is_set():
#         #         return False
#
#         # for p in plist:
#         #     print(f'     ... {p}')
#         #     if not build_all(p, plist[p]):
#         #         ok = False
#
#         #     if exit.is_set():
#         #         return False
#
#         # ... confirm uhppoted and submodule binary checksums match
#         for p in plist:
#             print(f'     ... {p}')
#             if not checksum(p, plist[p], 'development'):
#                 ok = False
#
#             if exit.is_set():
#                 return False
#
#     return True


def release(project, p, version, exit):
    print(f'>>>> build release for {project} ({version.version(project)})')

    # ... update for release and build
    if not update_release(project, p):
        return False

    # ... confirm go.mod has release versions of uhppote-core and uhppoted-lib
    if not updated_for_release(project, p, version):
        return False

    if not _release(project, p, version.version(project)):
        return False

    return True


#        #     if not uncommitted(p, plist[p]):
#         #         print(f'     ... {p} has uncommitted changes')
#         #         say(f'{p} has uncommitted changes')
#         #         exit.wait(10)
#         #         while not uncommitted(p, plist[p]) and not exit.is_set():
#         #             exit.wait(10)
#         #
#         #     if exit.is_set():
#         #         return False
#
#         # # ... confirm uhppoted and submodule binary checksums match
#         # print(f'     >>> verifying checksums')
#         # for p in plist:
#         #     print(f'     ... {p}')
#         #     if not checksum(p, plist[p], version.version(p)):
#         #         ok = False
#         #         print(f'     ... {p} binary checksums do not match')
#         #         say(f'{p} binary checksums do not match')
#
#         if exit.is_set():
#             return False
#
#         if ok:
#             return True
#
#     return True

# def release(projects, version, exit):
#     print('>>>> release builds ({version})')
#
#     plist = {p: projects[p] for p in itertools.filterfalse(lambda p: p in ignore, projects)}
#
#     while True:
#         ok = True
#
#         # print(f'     >>> building releases')
#         # for p in plist:
#         #     # ... update for release and build
#         #     print(f'     ... {p}')
#         #     if not update_release(p, plist[p]):
#         #         ok = False
#         #
#         #     if not _release(p, plist[p], version.version(p)):
#         #         ok = False
#         #
#         #     if not uncommitted(p, plist[p]):
#         #         print(f'     ... {p} has uncommitted changes')
#         #         say(f'{p} has uncommitted changes')
#         #         exit.wait(10)
#         #         while not uncommitted(p, plist[p]) and not exit.is_set():
#         #             exit.wait(10)
#         #
#         #     if exit.is_set():
#         #         return False
#
#         # # ... confirm uhppoted and submodule binary checksums match
#         # print(f'     >>> verifying checksums')
#         # for p in plist:
#         #     print(f'     ... {p}')
#         #     if not checksum(p, plist[p], version.version(p)):
#         #         ok = False
#         #         print(f'     ... {p} binary checksums do not match')
#         #         say(f'{p} binary checksums do not match')
#
#         if exit.is_set():
#             return False
#
#         if ok:
#             return True
#
#     return True


def checkout(project, info):
    try:
        command = f"cd {info['folder']} && git checkout {info['branch']}"
        subprocess.run(command, shell=True, check=True)
    except subprocess.CalledProcessError:
        raise Exception(f"command 'checkout {project}' failed")

    return True


def update(project, info):
    try:
        folder = info['folder']
        command = f'cd {folder} && make update'
        subprocess.run(command, shell=True, check=True)
    except subprocess.CalledProcessError:
        raise Exception(f"command 'update {project}' failed")

    return True


def update_release(project, info):
    try:
        folder = info['folder']
        command = f'cd {folder} && make update-release'
        subprocess.run(command, shell=True, check=True)
    except subprocess.CalledProcessError:
        raise Exception(f"command 'update {project}' failed")

    return True


def updated_for_release(project, info, version):
    try:
        folder = info['folder']
        path = os.path.join(folder, 'go.mod')

        if os.path.isfile(path):
            core = None
            lib = None
            r = re.compile('(?:require\s+)?(\S+)\s+(\S+)')

            with open(path, 'rt') as f:
                while line := f.readline():
                    match = r.match(line.strip())
                    if match:
                        if match.group(1) == 'github.com/uhppoted/uhppote-core':
                            core = match.group(2)
                        if match.group(1) == 'github.com/uhppoted/uhppoted-lib':
                            lib = match.group(2)

            if core and f'{core}' != f'{version}':
                raise Exception(f"'{project}' has not been updated to the release version of uhppote-core")

            if lib and f'{lib}' != f'{version}':
                raise Exception(f"{project}' has not been updated to the release version of uhppoted-lib")

    except subprocess.CalledProcessError:
        raise Exception(f"command 'update {project}' failed")

    return True


def uncommitted(project, info):
    try:
        command = f"cd {info['folder']} && git remote update"
        subprocess.run(command, shell=True, check=True)

        command = f"cd {info['folder']} && git status -uno"
        result = subprocess.check_output(command, shell=True)

        if not 'Changes not staged for commit' in str(result):
            return True

    except subprocess.CalledProcessError:
        raise Exception(f"{project}: command 'git status' failed")

    return False


def build(project, info):
    try:
        folder = info['folder']
        command = f'cd {folder} && make build'
        subprocess.run(command, shell=True, check=True)
    except subprocess.CalledProcessError:
        raise Exception(f"command 'update {project}' failed")

    return True


def build_all(project, info):
    try:
        folder = info['folder']
        command = f'cd {folder} && make build-all'
        subprocess.run(command, shell=True, check=True)
        return True
    except subprocess.CalledProcessError:
        raise Exception(f"command 'update {project}' failed")


def _release(project, info, version):
    try:
        folder = info['folder']
        command = f'cd {folder} && make release DIST={project}_{version}'
        subprocess.run(command, shell=True, check=True)
    except subprocess.CalledProcessError:
        raise Exception(f"command 'release {project}' failed")

    return True


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

    return True


def hash(file):
    hash = hashlib.sha256()

    with open(file, "rb") as f:
        bytes = f.read(65536)
        hash.update(bytes)

    return hash.hexdigest()
