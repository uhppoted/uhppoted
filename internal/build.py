import itertools
import subprocess

ignore = []


def prerelease(projects, version, exit):
    print('>>>> prerelease builds (v{version})')

    it = itertools.filterfalse(lambda p: p in ignore, projects)
    plist = {p: projects[p] for p in it}

    while True:
        ok = True
        # ... update and build
        print(f'     ... rebuilding all projects')
        for p in plist:
            print(f'     ... {p}')
            if not update(p, plist[p]):
                ok = False

            if exit.is_set():
                return False

        for p in plist:
            print(f'     ... {p}')
            if not build(p, plist[p]):
                ok = False

            if exit.is_set():
                return False

        # ... uncommitted changes?
        print('     ... checking for uncommitted changes')
        for p in plist:
            print(f'     ... {p}')
            uncommitted(p, plist[p])

            if exit.is_set():
                return False

        # ... checkout and rebuild
        print('     ... checking out latest github version')
        for p in plist:
            print(f'     ... {p}')
            checkout(p, plist[p])

            if exit.is_set():
                return False

        for p in plist:
            print(f'     ... {p}')
            if not build(p, plist[p]):
                ok = False

            if exit.is_set():
                return False

        if ok:
            break

    return True


def update(project, info):
    try:
        folder = info['folder']
        command = f'cd {folder} && make update'
        subprocess.run(command, shell=True, check=True)
        return True
    except subprocess.CalledProcessError:
        raise Exception(f"command 'update {project}' failed")


def uncommitted(project, info):
    try:
        command = f"cd {info['folder']} && git remote update"
        subprocess.run(command, shell=True, check=True)

        command = f"cd {info['folder']} && git status -uno"
        result = subprocess.check_output(command, shell=True)

        if (not project in ignore) and 'Changes not staged for commit' in str(result):
            raise Exception(f"{project} has uncommitted changes")

    except subprocess.CalledProcessError:
        raise Exception(f"{project}: command 'git status' failed")


def checkout(project, info):
    try:
        command = f"cd {info['folder']} && git checkout {info['branch']}"
        subprocess.run(command, shell=True, check=True)
    except subprocess.CalledProcessError:
        raise Exception(f"command 'checkout {project}' failed")

def build(project, info):
    try:
        folder = info['folder']
        command = f'cd {folder} && make build'
        subprocess.run(command, shell=True, check=True)
        return True
    except subprocess.CalledProcessError:
        raise Exception(f"command 'update {project}' failed")


