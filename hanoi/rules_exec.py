import sched, time, os, random, threading
import conf
import cherrypy.process.plugins
import cherrypy
import sensors_store
import diff_serve
import traceback
import db
import sensors_store, rules_store
import traceback
import paging_exec

running = False
DEFAULT_INTERVAL = 160 # seconds

def latest_rule_exec(rule, group):
    exec_vars = {'rule': rule}
    for sensor_alias in rule['sensors']:
        last_sensor_state = sensors_store.getIfExists(rule['sensors'][sensor_alias], group)
        if not last_sensor_state:
            print "Aborting exec of rule '" + rule['name'] + "', missing sensor dependency", rule['sensors'][sensor_alias]
            return True, True, {}
        exec_vars[sensor_alias] = last_sensor_state

    try:
        result = eval(rule['condition'], exec_vars)
        print "Rule result: ", result
        return bool(result), False, {}
    except KeyError, e:
        if e.args[0] == 'val':
            print "Sensor not ready? aborting rule"
            return True, True, {}
        else:
            print "Exception running rule: ", e
            traceback.print_exc()
            return False, False, {'exception': str(e)}
    except Exception, e:
        print "Exception running rule: ", e
        traceback.print_exc()
        return False, False, {'exception': str(e)}


def run_rule(rule, group):
    print "Now running rule:", group.replace(' ', '_')+"."+rule['name']
    ok, noop, state = False, True, {'exception': 'Invalid type'}
    if rule['type'] == "latest":
        ok, noop, state = latest_rule_exec(rule, group)
    ruleState = rules_store.get(rule['type'], rule['name'], group)
    ruleState.ok = ok
    ruleState.noop = noop
    ruleState.state = state
    diff_serve.notify('rules')
    if paging_exec.should_run_paging_rules():
        paging_exec.exec_paging_rules()


def maxInterval(config):
    '''Returns the largest interval registered for the rule.'''
    maxi = 0
    for group in config['groups']:
        rules = config['groups'][group]['rules']
        for rule in rules:
            interval = DEFAULT_INTERVAL # seconds
            if 'interval' in rule:
                interval = rule['interval']
            maxi = max(maxi, interval)
    return maxi

def createRuleEvents(s):
    config = conf.get_config()
    biggest_interval = maxInterval(config)
    for group in config['groups']:
        rules = config['groups'][group]['rules']
        for rule in rules:
            interval = DEFAULT_INTERVAL # seconds
            if 'interval' in rule:
                interval = rule['interval']
            if interval < 1:
                raise AttributeError("Interval should be > 1")
            interval += random.randint(1, 10) #entropy
            if (biggest_interval / interval) > 1:
                # Make registrations for rules with small intervals which fit within the biggest interval.
                for x in range(1, biggest_interval / interval):
                    s.enter((x*interval)+interval, 2, run_rule, (rule, group))
            s.enter(interval, 2, run_rule, (rule, group))


def startRuleScheduler(s):
    print "startRuleScheduler()"
    global running
    running = True

    while running:
        createRuleEvents(s)
        s.run()

def start():
    global running
    wd = None
    s = sched.scheduler(time.time, time.sleep)
    def stop():
        map(s.cancel, s.queue)
        wd.cancel()
        map(s.cancel, s.queue)
        raise Exception, "Stopping exception"
    wd = cherrypy.process.plugins.BackgroundTask(1, startRuleScheduler, [s])
    cherrypy.engine.subscribe('stop', stop)
    wd.start()
