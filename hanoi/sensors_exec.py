import sched, time, os, random, threading
import conf
import cherrypy.process.plugins
import cherrypy
import sensors_store
import diff_serve
import traceback
import db

running = False
DEFAULT_INTERVAL = 120 # seconds

def run_sensor(sensor, group):
    print "Now running sensor:", group.replace(' ', '_')+"."+sensor['name']
    ss = sensors_store.get(sensor['type'], sensor['name'], group)
    filename = os.path.join(os.getcwd(), 'sensors/' + sensor['type'] + '.py')
    try:
        with open(filename) as fp:
            code = compile(fp.read(), filename, "exec")
            def write():
                diff_serve.notify('sensors')
                db.write_sensor_result(group+"."+sensor['name'], True, ss.last_result)
            exec code in {'group': group, 'sensor': sensor, 'state': ss, 'write': write}
    except Exception, e:
        print "Exception running sensor: ", e
        traceback.print_exc()
        ss.ok = False
        ss.exception = e
        ss.err_msg = str(e)
        diff_serve.notify('sensors')
        db.write_sensor_result(group+"."+sensor['name'], False, ss.err_msg)


def maxInterval(config):
    '''Returns the largest interval registered for the sensors.'''
    maxi = 0
    for group in config['groups']:
        sensors = config['groups'][group]['sensors']
        for sensor in sensors:
            interval = DEFAULT_INTERVAL # seconds
            if 'interval' in sensor:
                interval = sensor['interval']
            maxi = max(maxi, interval)
    return maxi

def createSensorEvents(s):
    config = conf.get_config()
    biggest_interval = maxInterval(config)
    for group in config['groups']:
        sensors = config['groups'][group]['sensors']
        for sensor in sensors:
            interval = DEFAULT_INTERVAL # seconds
            if 'interval' in sensor:
                interval = sensor['interval']
            if interval < 1:
                raise AttributeError("Interval should be > 1")
            interval += random.randint(1, 10) #entropy
            if (biggest_interval / interval) > 1:
                # Make registrations for sensors with small intervals which fit within the biggest interval.
                for x in range(1, biggest_interval / interval):
                    s.enter((x*interval)+interval, 2, run_sensor, (sensor, group))
            s.enter(interval, 2, run_sensor, (sensor, group))


def startSensorScheduler(s):
    print "startSensorScheduler()"
    global running
    running = True

    while running:
        createSensorEvents(s)
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
    wd = cherrypy.process.plugins.BackgroundTask(1, startSensorScheduler, [s])
    cherrypy.engine.subscribe('stop', stop)
    wd.start()
