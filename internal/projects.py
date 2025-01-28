from dataclasses import dataclass
# from submodules.uhppote_core import uhppote_core_changelog


@dataclass
class Project:
    name: str
    folder: str
    branch: str = 'main'
    binary: str = None
    packaging: str = None


# yapf: disable
PROJECTS = [
    Project(name='uhppote-core',                folder='./uhppote-core'),
    Project(name='uhppoted-lib',                folder='./uhppoted-lib'),
    Project(name='uhppote-simulator',           folder='./uhppote-simulator',           binary='uhppote-simulator'),
    Project(name='uhppote-cli',                 folder='./uhppote-cli',                 binary='uhppote-cli'),
    Project(name='uhppoted-rest',               folder='./uhppoted-rest',               binary='uhppoted-rest'),
    Project(name='uhppoted-mqtt',               folder='./uhppoted-mqtt',               binary='uhppoted-mqtt'),
    Project(name='uhppoted-httpd',              folder='./uhppoted-httpd',              binary='uhppoted-httpd'),
    Project(name='uhppoted-tunnel',             folder='./uhppoted-tunnel',             binary='uhppoted-tunnel'),
    Project(name='uhppoted-dll',                folder='./uhppoted-dll'),
    Project(name='uhppoted-codegen',            folder='./uhppoted-codegen',            binary='uhppoted-codegen'),
    Project(name='uhppoted-app-s3',             folder='./uhppoted-app-s3',             binary='uhppoted-app-s3'),
    Project(name='uhppoted-app-sheets',         folder='./uhppoted-app-sheets',         binary='uhppoted-app-sheets'),
    Project(name='uhppoted-app-wild-apricot',   folder='./uhppoted-app-wild-apricot',   binary='uhppoted-app-wild-apricot'),
    Project(name='uhppoted-app-db',             folder='./uhppoted-app-db',             binary='uhppoted-app-db'),
    Project(name='uhppoted-app-home-assistant', folder='./uhppoted-app-home-assistant', binary='uhppoted-app-home-assistant'),
    Project(name='uhppoted-python',             folder='./uhppoted-python',           packaging='python'),
    Project(name='uhppoted-nodejs',             folder='./uhppoted-nodejs',           packaging='nodejs'),
    Project(name='node-red-contrib-uhppoted',   folder='./node-red-contrib-uhppoted', packaging='nodejs'),
    Project(name='uhppoted-lib-dotnet',         folder='./uhppoted-lib-dotnet',       packaging='dotnet'),
    Project(name='uhppoted-wiegand',            folder='./uhppoted-wiegand'),
    Project(name='uhppoted',                    folder='.'),
]
# yapf: enable


def projects():
    return {
        p.name: {
            'folder': p.folder,
            'branch': p.branch,
            'binary': p.binary,
            'packaging': p.packaging,
        }
        for p in PROJECTS
    }

    # return {
    #     'uhppote-core': {
    #         'folder': './uhppote-core',
    #         'branch': 'main',
    #         # 'changelog': uhppote_core_changelog,
    #     },
    #     'uhppoted-lib': {
    #         'folder': './uhppoted-lib',
    #         'branch': 'main',
    #     },
    #     'uhppote-simulator': {
    #         'folder': './uhppote-simulator',
    #         'branch': 'main',
    #         'binary': 'uhppote-simulator'
    #     },
    #     'uhppote-cli': {
    #         'folder': './uhppote-cli',
    #         'branch': 'main',
    #         'binary': 'uhppote-cli'
    #     },
    #     'uhppoted-rest': {
    #         'folder': './uhppoted-rest',
    #         'branch': 'main',
    #         'binary': 'uhppoted-rest'
    #     },
    #     'uhppoted-mqtt': {
    #         'folder': './uhppoted-mqtt',
    #         'branch': 'main',
    #         'binary': 'uhppoted-mqtt'
    #     },
    #     'uhppoted-httpd': {
    #         'folder': './uhppoted-httpd',
    #         'branch': 'main',
    #         'binary': 'uhppoted-httpd'
    #     },
    #     'uhppoted-tunnel': {
    #         'folder': './uhppoted-tunnel',
    #         'branch': 'main',
    #         'binary': 'uhppoted-tunnel'
    #     },
    #     'uhppoted-dll': {
    #         'folder': './uhppoted-dll',
    #         'branch': 'main',
    #     },
    #     'uhppoted-codegen': {
    #         'folder': './uhppoted-codegen',
    #         'branch': 'main',
    #         'binary': 'uhppoted-codegen'
    #     },
    #     'uhppoted-app-s3': {
    #         'folder': './uhppoted-app-s3',
    #         'branch': 'main',
    #         'binary': 'uhppoted-app-s3'
    #     },
    #     'uhppoted-app-sheets': {
    #         'folder': './uhppoted-app-sheets',
    #         'branch': 'main',
    #         'binary': 'uhppoted-app-sheets'
    #     },
    #     'uhppoted-app-wild-apricot': {
    #         'folder': './uhppoted-app-wild-apricot',
    #         'branch': 'main',
    #         'binary': 'uhppoted-app-wild-apricot'
    #     },
    #     'uhppoted-app-db': {
    #         'folder': './uhppoted-app-db',
    #         'branch': 'main',
    #         'binary': 'uhppoted-app-db'
    #     },
    #     'uhppoted-app-home-assistant': {
    #         'folder': './uhppoted-app-home-assistant',
    #         'branch': 'main',
    #         'binary': 'uhppoted-app-home-assistant'
    #     },
    #     'uhppoted-nodejs': {
    #         'folder': './uhppoted-nodejs',
    #         'branch': 'main',
    #     },
    #     'uhppoted-python': {
    #         'folder': './uhppoted-python',
    #         'branch': 'main',
    #     },
    #     'node-red-contrib-uhppoted': {
    #         'folder': './node-red-contrib-uhppoted',
    #         'branch': 'main',
    #     },
    #     'uhppoted-lib-dotnet': {
    #         'folder': './uhppoted-lib-dotnet',
    #         'branch': 'main',
    #     },
    #     'uhppoted-wiegand': {
    #         'folder': './uhppoted-wiegand',
    #         'branch': 'main',
    #     },
    #     'uhppoted': {
    #         'folder': '.',
    #         'branch': 'main'
    #     }
    # }
