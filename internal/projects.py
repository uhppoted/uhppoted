from submodules.uhppote_core import uhppote_core_changelog


def projects():
    return {
        'uhppote-core': {
            'folder': './uhppote-core',
            'branch': 'main',
            # 'changelog': uhppote_core_changelog,
        },
        'uhppoted-lib': {
            'folder': './uhppoted-lib',
            'branch': 'main',
        },
        'uhppote-simulator': {
            'folder': './uhppote-simulator',
            'branch': 'main',
            'binary': 'uhppote-simulator'
        },
        'uhppote-cli': {
            'folder': './uhppote-cli',
            'branch': 'main',
            'binary': 'uhppote-cli'
        },
        'uhppoted-rest': {
            'folder': './uhppoted-rest',
            'branch': 'main',
            'binary': 'uhppoted-rest'
        },
        'uhppoted-mqtt': {
            'folder': './uhppoted-mqtt',
            'branch': 'main',
            'binary': 'uhppoted-mqtt'
        },
        'uhppoted-httpd': {
            'folder': './uhppoted-httpd',
            'branch': 'main',
            'binary': 'uhppoted-httpd'
        },
        'uhppoted-tunnel': {
            'folder': './uhppoted-tunnel',
            'branch': 'main',
            'binary': 'uhppoted-tunnel'
        },
        'uhppoted-dll': {
            'folder': './uhppoted-dll',
            'branch': 'main',
        },
        'uhppoted-codegen': {
            'folder': './uhppoted-codegen',
            'branch': 'main',
            'binary': 'uhppoted-codegen'
        },
        'uhppoted-app-s3': {
            'folder': './uhppoted-app-s3',
            'branch': 'main',
            'binary': 'uhppoted-app-s3'
        },
        'uhppoted-app-sheets': {
            'folder': './uhppoted-app-sheets',
            'branch': 'main',
            'binary': 'uhppoted-app-sheets'
        },
        'uhppoted-app-wild-apricot': {
            'folder': './uhppoted-app-wild-apricot',
            'branch': 'main',
            'binary': 'uhppoted-app-wild-apricot'
        },
        'uhppoted-app-db': {
            'folder': './uhppoted-app-db',
            'branch': 'main',
            'binary': 'uhppoted-app-db'
        },
        'uhppoted-app-home-assistant': {
            'folder': './uhppoted-app-home-assistant',
            'branch': 'main',
            'binary': 'uhppoted-app-home-assistant'
        },
        'uhppoted-nodejs': {
            'folder': './uhppoted-nodejs',
            'branch': 'main',
        },
        'uhppoted-python': {
            'folder': './uhppoted-python',
            'branch': 'main',
        },
        'node-red-contrib-uhppoted': {
            'folder': './node-red-contrib-uhppoted',
            'branch': 'main',
        },
        'uhppoted-lib-dotnet': {
            'folder': './uhppoted-lib-dotnet',
            'branch': 'main',
        },
        'uhppoted-wiegand': {
            'folder': './uhppoted-wiegand',
            'branch': 'main',
        },
        'uhppoted': {
            'folder': '.',
            'branch': 'main'
        }
    }
