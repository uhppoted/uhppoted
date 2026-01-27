import yaml

import projects


class Version:

    @classmethod
    def read(cls, file):
        path = file if not None else '.versions'

        with open(path, 'r', encoding='utf-8') as f:
            yml = yaml.safe_load(f)

            version = yml['versions']['default']
            wild_apricot = yml['versions']['uhppoted-app-wild-apricot']
            node_red = yml['versions']['node-red-contrib-uhppoted']
            home_assistant = yml['versions']['uhppoted-app-home-assistant']

            return Version(version, wild_apricot, node_red, home_assistant)

    def __init__(self, version, wild_apricot, node_red, home_assistant):
        if version != 'development' and not version.startswith('v'):
            self._version = f'v{version}'
        else:
            self._version = f'{version}'

        if wild_apricot != 'development' and not wild_apricot.startswith('v'):
            self._uhppoted_app_wild_apricot = f'v{wild_apricot}'
        else:
            self._uhppoted_app_wild_apricot = f'{wild_apricot}'

        if node_red != 'development' and not node_red.startswith('v'):
            self._node_red = f'v{node_red}'
        else:
            self._node_red = f'{node_red}'

        if home_assistant != 'development' and not home_assistant.startswith('v'):
            self._uhppoted_app_home_assistant = f'v{home_assistant}'
        else:
            self._uhppoted_app_home_assistant = f'{home_assistant}'

    def __str__(self):
        return self._version

    def version(self, project):
        name = project.name if isinstance(project, projects.Project) else f'{project}'

        if name == 'node-red-contrib-uhppoted':
            return self._node_red
        elif name == 'uhppoted-app-wild-apricot':
            return self._uhppoted_app_wild_apricot
        elif name == 'uhppoted-app-home-assistant':
            return self._uhppoted_app_home_assistant
        else:
            return self._version
