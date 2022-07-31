#!python3

import subprocess
import sys

def main():
    print()
    print("*** uhppoted release-all")
    print()

    try:
        list = projects()
        for p in list:
            update(p,list[p])

    except BaseException as x:
        print()
        print(f'*** ERROR  {x}')
        print()

        sys.exit(1)

def projects():
    return {
        'uhppote-core': {
            'folder': './uhppote-core'
        },
        'uhppoted-lib': {
            'folder': './uhppoted-lib'
        }
    }

def update(project,info):
    command = f"cd {info['folder']} && make update"
    result = subprocess.call(command, shell = True)
    if result != 0:
       raise Exception(f"command 'make update {project}' failed")


if __name__ == '__main__':
    main()

