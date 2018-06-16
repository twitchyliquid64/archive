"""Stores intantenous state for each of the paging rules."""

paging_rules_by_name = {}

class PagingRuleState(object):
    def __init__(self, name, group):
        self.name = name
        self.group = group
        self.rules_by_state = {}
        self.failing = False
        self.last_page_epoch_seconds = 0

    def ident(self):
        return self.group + '.' + self.name


def get(name, group):
    '''fetches a PagingRuleState object, registering one if it doesnt exist.'''
    global paging_rules_by_name
    ident = group + '.' + name
    if ident in paging_rules_by_name:
        return paging_rules_by_name[ident]
    else:
        paging_rules_by_name[ident] = PagingRuleState(name, group)
        return paging_rules_by_name[ident]
