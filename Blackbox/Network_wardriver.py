import wifiscan
import subprocess
import RPi.GPIO as GPIO
import time

class Network_wardriver(object):
	def __init__(self, interface, main):
		self.interface = interface
		self.main = main
		self.selection = 0
		self.wifimenu = self.main.addMenu(" WIFI WARDRIVER ","Press to start.", self.startMenuSelected)

	def startMenuSelected(self):
		self.start()


	def setup(self):
		time.sleep(0.5)
		subprocess.check_output(["/usr/local/sbin/airmon-ng", "check", "kill", self.interface])
		subprocess.check_output(["/usr/bin/killall", "gpsd"])
		subprocess.check_output(["/usr/sbin/gpsd", "/dev/ttyUSB0", "-F", "/var/run/gpsd.sock"])
		time.sleep(0.5)
		subprocess.check_output(["ifconfig", self.interface, "down"])
		time.sleep(1)
		subprocess.check_output(["/usr/local/sbin/airmon-ng", "check", "kill", self.interface])
		time.sleep(0.5)
		subprocess.check_output(["/usr/local/sbin/airmon-ng", "start", self.interface])
		time.sleep(2)
		subprocess.check_output(["/usr/local/sbin/airmon-ng", "check", "kill", self.interface])
		time.sleep(0.5)

	
	def start(self):
		self.main.displayMsg(" WIFI WARDRIVER ", "SETUP: "+self.interface)
		self.setup()
		self.captureProcess = subprocess.Popen(["/usr/local/sbin/airodump-ng", "--gpsd", "--output-format", "csv", "mon0"])
		self.main.updateMenu(self.wifimenu, " WIFI WARDRIVER ", "Capture started.", None)
