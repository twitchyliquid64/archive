"""Stores intantenous state for each of the sensors."""

sensors_by_ident = {}

class SensorState(object):
    def __init__(self, typ, name, group):
        self.type = typ
        self.name = name
        self.group = group
        self.ok = True
        self.err_msg = ''
        self.last_result = {}

    def ident(self):
        return self.group + '.' + self.name


def get(typ, name, group):
    '''fetches a SensorState object, registering one if it doesnt exist.'''
    global sensors_by_ident
    ident = group + '.' + name
    if ident in sensors_by_ident:
        if sensors_by_ident[ident].type != typ:
            raise Exception('Sensor already registered with different type')
        return sensors_by_ident[ident]
    else:
        sensors_by_ident[ident] = SensorState(typ, name, group)
        return sensors_by_ident[ident]

def getIfExists(name, group):
    '''fetches a SensorState object, returning None if it doesnt exist.'''
    global sensors_by_ident
    ident = group + '.' + name
    if ident in sensors_by_ident:
        return sensors_by_ident[ident]
    return None
