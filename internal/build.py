import subprocess

ignore = []


def prerelease(projects, version, exit):
    print(f'>>>> prerelease builds (v{version})')

    while True:
        ok = True
        # ... build
        print(f'     ... building all projects')
        for p in projects:
            if not p in ignore:
                project = projects[p]

                print(f'     ... {p}')
                if not update(p, project):
                    ok = False

                if exit.is_set():
                    return False

        # ... uncommitted changes?
        print(f'     ... checking for uncommitted changesbuilding all projects')
        for p in projects:
            if not p in ignore:
                project = projects[p]

                print(f'     ... {p}')
                uncommitted(p, project)

                if exit.is_set():
                    return False

        # ... checkout and rebuild
        print(f'     ... checking out latest github version')
        for p in projects:
            if not p in ignore:
                project = projects[p]

                print(f'     ... {p}')
                checkout(p, project)

                if exit.is_set():
                    return False

        if ok:
            break

    return True


def update(project, info):
    try:
        folder = info['folder']
        command = f'cd {folder} && make update && make build'
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
