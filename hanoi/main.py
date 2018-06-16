import cherrypy
import conf
from config_serve import ConfigServe
from diff_serve import DiffServe
import sensors_exec
import rules_exec
import sys
import json
from db import db_start, query

conf.config_path = sys.argv[2]

class Root(object):
    def __init__(self):
        self.conf_serve = ConfigServe()
        self.diff_serve = DiffServe()

    @cherrypy.expose
    def index(self):
        return open('index.html')

    @cherrypy.expose
    def chart(self, name):
        return open('chart.html')

    @cherrypy.expose
    def sensor_history(self, name):
        return json.dumps(query(name))

    @cherrypy.expose
    @cherrypy.tools.json_out()
    def config(self):
        return self.conf_serve.config()

    @cherrypy.expose
    @cherrypy.tools.json_out()
    def diff(self, update_key=0):
        return self.diff_serve.diff(int(update_key))

    @cherrypy.expose
    @cherrypy.tools.json_out()
    def notify_test(self, component='test'):
        return self.diff_serve.notify_test(component)

    @cherrypy.expose
    @cherrypy.tools.json_out()
    def sensors(self):
        return self.conf_serve.sensors()

    @cherrypy.expose
    @cherrypy.tools.json_out()
    def rules(self):
        return self.conf_serve.rules()

    @cherrypy.expose
    @cherrypy.tools.json_out()
    def paging_rules(self):
        return self.conf_serve.paging_rules()

cherrypy.config.update({
    'server.socket_host': '127.0.0.1',
    'server.socket_port': int(sys.argv[1]),
})

sensors_exec.start()
rules_exec.start()
db_start()
cherrypy.quickstart(Root(), '/', conf.Web_conf)
