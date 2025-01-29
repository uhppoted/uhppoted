import itertools
import subprocess
import re

from misc import say
from git import _uncommitted

import github

ignore = []


def publish(project, p, version, exit):
    print(f'>>>> publishing {project} ({version})')

    # ... confirm no uncommitted changes
    print(f'     >>> checking for uncommitted changes')
    if not _uncommitted(project, p, version, exit):
        raise Exception(f"{project} has uncommitted changes")

    # ... confirm changes have been pushed to github repo
    if not github.pushed(project, p):
        print(f'     ... {project} has unpushed changes')
        say(f'{project} has unpushed changes')
        exit.wait(10)
        while not pushed(project, p) and not exit.is_set():
            exit.wait(10)

        if exit.is_set():
            return False

    # ... confirm github release has been published
    exit.wait(10)
    while not github.published(project, p, version.version(p)) and not exit.is_set():
        exit.wait(10)

    if exit.is_set():
        return False

    # ... publish release to npm
    _publish(project, p, version.version(p))
    print(f'     ... {project} is waiting for release on npm')
    say(f'{project} is waiting for release on npm')

    exit.wait(10)
    while not published(project, p, version.version(p)) and not exit.is_set():
        exit.wait(10)

    if exit.is_set():
        return False

    say(f'{project} has been published to npm')
    return True


def _publish(project, info, version):
    try:
        folder = info.folder
        command = f"cd {folder} && make publish-npm"
        subprocess.run(command, shell=True, check=True)
    except subprocess.CalledProcessError:
        raise Exception(f"command 'publish-npm {project}' failed")

    return True


# FIXME
def published(project, info, version):
    # command = f"cd {info.folder} && git fetch --tags"
    # result = subprocess.call(command, shell=True)
    # if result != 0:
    #     raise Exception(f"command 'git fetch --tags' failed")
    # else:
    #     command = f"cd {info.folder} && git tag {version} --list"
    #     result = subprocess.check_output(command, shell=True)
    #
    #     if f'{version}' in str(result):
    #         print(f'     +++ {project} has been released')
    #         return True
    #     else:
    #         print(f'     ... {project} has not been released')
    #         return False
    return False
