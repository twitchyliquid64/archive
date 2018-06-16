import os
import json

config_path = 'configuration.json'

Web_conf = {
    '/': {
        'tools.sessions.on': True,
        'tools.staticdir.root': os.path.abspath(os.getcwd())
    },
    '/css': {
        'tools.staticdir.on': True,
        'tools.staticdir.dir': './css'
    },
    '/fonts': {
        'tools.staticdir.on': True,
        'tools.staticdir.dir': './fonts'
    },
    '/js': {
        'tools.staticdir.on': True,
        'tools.staticdir.dir': './js'
    }
}

config = None

def _readConfig():
    global config, config_path
    with open(config_path, mode='r') as conf_file:
        config = json.load(conf_file)
        _validateConfig(config)


class InvalidConfiguration(Exception):
    pass

def _validateConfig(conf_obj):
    if 'groups' not in conf_obj:
        raise InvalidConfiguration('Missing top level keys')
    _validateGroups(conf_obj['groups'])
    _validateSensors(conf_obj['groups'])
    _validateRules(conf_obj['groups'])
    _validatePagingRules(conf_obj['groups'])

def _validateGroups(groups_obj):
    for group_name in groups_obj:
        if not all (k in groups_obj[group_name] for k in ('sensors', 'rules', 'paging_rules')):
            raise InvalidConfiguration('Missing keys in group item')

def _validateSensors(groups_obj):
    for group_name in groups_obj:
        for sensor in groups_obj[group_name]['sensors']:
            if 'name' not in sensor:
                raise InvalidConfiguration('Missing key \'name\' in sensor item: ' + str(sensor))

def _validateRules(groups_obj):
    for group_name in groups_obj:
        for rule in groups_obj[group_name]['rules']:
            if 'name' not in rule:
                raise InvalidConfiguration('Missing key \'name\' in rule item: ' + str(rule))
            if 'type' not in rule:
                raise InvalidConfiguration('Missing key \'type\' in rule item: ' + str(rule))


def _validatePagingRules(groups_obj):
    for group_name in groups_obj:
        for rule in groups_obj[group_name]['paging_rules']:
            if 'name' not in rule:
                raise InvalidConfiguration('Missing key \'name\' in paging_rules item: ' + str(rule))
            if 'rules' not in rule:
                raise InvalidConfiguration('Missing key \'rules\' in paging_rules item: ' + str(rule))
            if 'page' not in rule:
                raise InvalidConfiguration('Missing key \'page\' in paging_rules item: ' + str(rule))
            if rule['page']['type'] == "GMAIL":
                if 'address' not in rule['page']:
                    raise InvalidConfiguration('Missing key \'page.address\' in paging_rules item: ' + str(rule))


def get_config():
    global config
    if not config:
        _readConfig()
    return config
