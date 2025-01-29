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
    Project(name='uhppoted-lib-python',         folder='./uhppoted-lib-python',       packaging='python'),
    Project(name='uhppoted-lib-dotnet',         folder='./uhppoted-lib-dotnet',       packaging='dotnet'),
    Project(name='uhppoted-nodejs',             folder='./uhppoted-nodejs',           packaging='nodejs'),
    Project(name='node-red-contrib-uhppoted',   folder='./node-red-contrib-uhppoted', packaging='nodejs'),
    Project(name='uhppoted-wiegand',            folder='./uhppoted-wiegand'),
    Project(name='uhppoted',                    folder='.'),
]
# yapf: enable


def projects():
    return {p.name: p for p in PROJECTS}
