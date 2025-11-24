from dataclasses import dataclass
from dataclasses import field
from typing import List

# from submodules.uhppote_core import uhppote_core_changelog


@dataclass
class State:
    released: bool = False
    changelog: bool = False
    readme: bool = False
    version: bool = False
    committed: bool = False
    pushed: bool = False
    prepared: bool = False
    release_notes: bool = False
    published: bool = False


@dataclass
class Project:
    name: str
    folder: str
    branch: str = 'main'
    binary: str = None
    packaging: str = None
    package: str = None
    dependencies: List[str] = field(default_factory=list)
    state: State = field(default_factory=State)


# yapf: disable
PROJECTS = [
    Project(name='uhppote-core',
            folder='./uhppote-core'),

    Project(name='uhppoted-lib',
            folder='./uhppoted-lib',
            dependencies=['uhppote-core']),

    Project(name='uhppoted-lib-python',
            folder='./uhppoted-lib-python',
            packaging='python'),

    Project(name='uhppoted-lib-dotnet',
            folder='./uhppoted-lib-dotnet',
            packaging='dotnet'),

    Project(name='uhppoted-lib-nodejs',
            folder='./uhppoted-lib-nodejs',
            packaging='nodejs',
            package='uhppoted'),

    Project(name='node-red-contrib-uhppoted',
            folder='./node-red-contrib-uhppoted',
            packaging='nodejs',
            package='node-red-contrib-uhppoted'),

    Project(name='uhppote-simulator',
            folder='./uhppote-simulator',
            binary='uhppote-simulator',
            dependencies=['uhppote-core']),

    Project(name='uhppote-cli',
            folder='./uhppote-cli',
            binary='uhppote-cli',
            dependencies=['uhppote-core', 'uhppoted-lib']),

    Project(name='uhppoted-rest',
            folder='./uhppoted-rest',
            binary='uhppoted-rest',
            dependencies=['uhppote-core', 'uhppoted-lib']),

    Project(name='uhppoted-mqtt',
            folder='./uhppoted-mqtt',
            binary='uhppoted-mqtt',
            dependencies=['uhppote-core', 'uhppoted-lib']),

    Project(name='uhppoted-httpd',
            folder='./uhppoted-httpd',
            binary='uhppoted-httpd',
            dependencies=['uhppote-core', 'uhppoted-lib']),

    Project(name='uhppoted-tunnel',
            folder='./uhppoted-tunnel',
            binary='uhppoted-tunnel',
            dependencies=['uhppote-core', 'uhppoted-lib']),

    Project(name='uhppoted-dll',
            folder='./uhppoted-dll',
            dependencies=['uhppote-core', 'uhppoted-lib']),

    Project(name='uhppoted-codegen',
            folder='./uhppoted-codegen',
            binary='uhppoted-codegen',
            dependencies=['uhppote-core', 'uhppoted-lib']),

    Project(name='uhppoted-app-s3',
            folder='./uhppoted-app-s3',
            binary='uhppoted-app-s3',
            dependencies=['uhppote-core', 'uhppoted-lib']),

    Project(name='uhppoted-app-sheets',
            folder='./uhppoted-app-sheets',
            binary='uhppoted-app-sheets',
            dependencies=['uhppote-core', 'uhppoted-lib']),

    Project(name='uhppoted-app-wild-apricot',
            folder='./uhppoted-app-wild-apricot',
            binary='uhppoted-app-wild-apricot',
            dependencies=['uhppote-core', 'uhppoted-lib'],
            packaging='go'),

    Project(name='uhppoted-app-db',
            folder='./uhppoted-app-db',
            binary='uhppoted-app-db',
            dependencies=['uhppote-core', 'uhppoted-lib']),

    Project(name='uhppoted-app-home-assistant',
            folder='./uhppoted-app-home-assistant',
            binary='uhppoted-app-home-assistant',
            dependencies=['uhppoted-lib-python']),

    Project(name='uhppoted-wiegand',
            folder='./uhppoted-wiegand'),

    Project(name='uhppoted',
            folder='.',
            dependencies=['*']),
]
# yapf: enable


def projects():
    return {p.name: p for p in PROJECTS}
