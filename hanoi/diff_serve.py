"""Responsible for keeping track of whats changed, so polling requests get a diff which they can use to update their browser state."""

import cherrypy
import conf
import json


latest_update_key = 0
latest_update_keys = {}

# On update, latest_update_key is incremented and latest_update_keys has a new value written - the component which should be updated.

def notify(component):
    global latest_update_key, latest_update_keys
    latest_update_keys[latest_update_key] = component
    latest_update_key += 1

    if latest_update_key > 100:
        del latest_update_keys[latest_update_key-98]


class DiffServe(object):

    def diff(self, update_key):
        global latest_update_key, latest_update_keys

        if update_key >= latest_update_key:
            return {'up_to_date': True, 'key': latest_update_key}

        if (latest_update_key - update_key) > 50:
            return {'up_to_date': False, 'updates': ['all'], 'key': latest_update_key}

        to_update = set()
        for x in range(update_key, latest_update_key):
            if x in latest_update_keys:
                to_update.add(latest_update_keys[x])
        return {'up_to_date': False, 'updates': list(to_update), 'key': latest_update_key}

    def notify_test(self, component):
        global latest_update_key, latest_update_keys
        notify(component)
        return latest_update_key
