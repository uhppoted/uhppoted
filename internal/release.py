#!python3

import argparse
import subprocess
import sys


def main():
    print()
    print("*** uhppoted release-all")
    print()

    if len(sys.argv) < 2:
        usage()
        return -1

    parser = argparse.ArgumentParser(description='release --version=<version>')

    parser.add_argument('--version',
                        type=str,
                        default='development',
                        help='release version e.g. v0.8.1')

    args = parser.parse_args()
    version = args.version

    try:
        print(f'VERSION: {version}')

        list = projects()
        for p in list:
            print(f'... releasing {p}')
            print()

            checkout(p, list[p])
            update(p, list[p])
            build(p, list[p])
            release(p, list[p], version)

        print()
        print(f'*** OK!')
        print()
        say('OK')

    except BaseException as x:
        print()
        print(f'*** ERROR  {x}')
        print()
        say('ERROR')

        sys.exit(1)


def projects():
    return {
        'uhppote-core': {
            'folder': './uhppote-core',
            'branch': 'master'
        },
        'uhppote-simulator': {
            'folder': './uhppote-simulator',
            'branch': 'master'
        },
        'uhppoted-dll': {
            'folder': './uhppoted-dll',
            'branch': 'master'
        },
        'uhppoted-lib': {
            'folder': './uhppoted-lib',
            'branch': 'master'
        },
        'uhppote-cli': {
            'folder': './uhppote-cli',
            'branch': 'master'
        },
        'uhppoted-rest': {
            'folder': './uhppoted-rest',
            'branch': 'master'
        },
        'uhppoted-mqtt': {
            'branch': 'master',
            'folder': './uhppoted-mqtt',
        },
        'uhppoted-httpd': {
            'folder': './uhppoted-httpd',
            'branch': 'master'
        },
        'uhppoted-tunnel': {
            'folder': './uhppoted-tunnel',
            'branch': 'master'
        },
        'uhppoted-app-s3': {
            'folder': './uhppoted-app-s3',
            'branch': 'master'
        },
        'uhppoted-app-sheets': {
            'folder': './uhppoted-app-sheets',
            'branch': 'master'
        },
        'uhppoted-app-wild-apricot': {
            'folder': './uhppoted-app-wild-apricot',
            'branch': 'master'
        },
        'uhppoted-nodejs': {
            'folder': './uhppoted-nodejs',
            'branch': 'master'
        },
        'node-red-contrib-uhppoted': {
            'folder': './node-red-contrib-uhppoted',
            'branch': 'master'
        },
        'uhppoted': {
            'folder': '.',
            'branch': 'master'
        }
    }


def checkout(project, info):
    command = f"cd {info['folder']} && git checkout {info['branch']}"
    result = subprocess.call(command, shell=True)
    if result != 0:
        raise Exception(f"command 'checkout {project}' failed")


def update(project, info):
    command = f"cd {info['folder']} && make update && make build"
    result = subprocess.call(command, shell=True)
    if result != 0:
        raise Exception(f"command 'update {project}' failed")


def build(project, info):
    command = f"cd {info['folder']} && make update-release && make build-all"
    result = subprocess.call(command, shell=True)
    if result != 0:
        raise Exception(f"command 'build {project}' failed")


def release(project, info, version):
    command = f"cd {info['folder']} && make release DIST={project}_{version}"
    result = subprocess.call(command, shell=True)
    if result != 0:
        raise Exception(f"command 'build {project}' failed")


def say(msg):
    subprocess.call(f'say {msg}', shell=True)

def usage():
    print()
    print('  Usage: python release.py <options>')
    print()
    print('  Options:')
    print('    --version <version>  Release version e.g. v0.8.1')
    print()


if __name__ == '__main__':
    main()
