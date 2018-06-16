import cherrypy
import conf
import json
import sensors_store, rules_store, paging_store

class ConfigServe(object):
    def config(self):
        return conf.get_config()

    def sensors(self):
        groups = conf.get_config()['groups'].keys()
        out = {}
        for group in groups:
            out[group] = conf.get_config()['groups'][group]['sensors']
            for x in xrange(len(out[group])):
                sensor_state = sensors_store.get(out[group][x]['type'], out[group][x]['name'], group)
                out[group][x]['ok'] = sensor_state.ok
                out[group][x]['last_result'] = sensor_state.last_result
                out[group][x]['err_msg'] = sensor_state.err_msg
        return out

    def rules(self):
        groups = conf.get_config()['groups'].keys()
        out = {}
        for group in groups:
            out[group] = conf.get_config()['groups'][group]['rules']
            for x in xrange(len(out[group])):
                sensor_state = rules_store.get(out[group][x]['type'], out[group][x]['name'], group)
                out[group][x]['ok'] = sensor_state.ok
                out[group][x]['noop'] = sensor_state.noop
                out[group][x]['state'] = sensor_state.state
        return out

    def paging_rules(self):
        groups = conf.get_config()['groups'].keys()
        out = {}
        for group in groups:
            out[group] = conf.get_config()['groups'][group]['paging_rules']
            for x in xrange(len(out[group])):
                sensor_state = paging_store.get(out[group][x]['name'], group)
                out[group][x]['ok'] = not sensor_state.failing
                out[group][x]['ruleState'] = sensor_state.rules_by_state
                out[group][x]['lastPage'] = sensor_state.last_page_epoch_seconds
        return out
