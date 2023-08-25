import subprocess


def already_released(project, info, version):
    command = f"cd {info['folder']} && git fetch --tags"
    result = subprocess.call(command, shell=True)
    if result != 0:
        raise Exception(f"command 'git fetch --tags' failed")
    else:
        command = f"cd {info['folder']} && git tag {version} --list"
        result = subprocess.check_output(command, shell=True)

        if f'{version}' in str(result):
            print(f'     ... {project} has been released')
            return True
        else:
            print(f'     ... {project} has not been released')
            return False
