"""Stores intantenous state for each of the rules."""

rules_by_ident = {}

class RuleState(object):
    def __init__(self, typ, name, group):
        self.type = typ
        self.name = name
        self.group = group
        self.ok = True
        self.noop = True
        self.state = {}

    def ident(self):
        return self.group + '.' + self.name


def get(typ, name, group):
    '''fetches a RuleState object, registering one if it doesnt exist.'''
    global rules_by_ident
    ident = group + '.' + name
    if ident in rules_by_ident:
        if rules_by_ident[ident].type != typ:
            raise Exception('Rule already registered with different type')
        return rules_by_ident[ident]
    else:
        rules_by_ident[ident] = RuleState(typ, name, group)
        return rules_by_ident[ident]

def getIfExists(name, group):
    '''fetches a RuleState object, returning None if it doesnt exist.'''
    global rules_by_ident
    ident = group + '.' + name
    if ident in rules_by_ident:
        return rules_by_ident[ident]
    return None
