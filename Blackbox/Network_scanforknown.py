import csv
import wifiscan
import subprocess
import os
import time
import RPi.GPIO as GPIO

class Network_scanforknown(object):
	VERSION = "0.01 alpha"
	def __init__(self, interface, main):
		print "      \tNetwork-scanforknown"
		print "      \tVersion: "+self.VERSION
		print "      \tNow scanning for networks"
		self.interface = interface
		self.main = main
		self.net = ""
		
		self.netConnected = False
		self.DHCPleased = False
		self.wpaProcess = None
		
		self.scan()

	def scan(self):
		self.loadKnown()
		self.results = wifiscan.scan(self.interface)
		self.search()

	def search(self):
		for network in self.results:
			for known in self.known:
				if network["Name"] == known[0]:
					print "     \tFound wifi candidate: "+str(network)
					self.main.displayMsg(None, "C: "+str(network["Name"]))
					self.connect(known)

	def loadKnown(self):
		self.known = []
		with open('/knownNetworks.csv', 'rb') as inputfile:
			reader = csv.reader(inputfile, delimiter=',',quotechar='"')
			for network in reader:
				self.known.append(network)

	def connect(self,networkdetails):
		SSID = networkdetails[0]
		self.net = SSID
		passkey = networkdetails[2]
		outconf = subprocess.check_output(["wpa_passphrase", SSID, passkey])
		fo = open("/wpa.conf", "w+")
		fo.write( "ctrl_interface=DIR=/var/run/wpa_supplicant\n" )
		fo.write( outconf )
		fo.close()
		fpath = "/wpa.conf"
		self.wpaProcess = subprocess.Popen(["wpa_supplicant", "-c"+fpath, "-i"+self.interface])
		timeout = time.time() + 40
		while timeout > time.time():
			time.sleep(1)
			chk = subprocess.check_output(['wpa_cli', 'status'])
			if "COMPLETED" in chk:
				self.netConnected = True
				self.DHCP()
				return
		self.main.displayMsg(None, "C: ASSOC ERR")
		time.sleep(2)
		self.shutdown()

	def DHCP(self):
		self.main.displayMsg(None, "C: Getting IP")
		
		subprocess.Popen(["ip", "route", "del", "default", "dev", "eth0"]).communicate()
		
		proc = subprocess.Popen(["dhclient", self.interface])
		proc.communicate()
		if proc.returncode != 0:
			self.main.displayMsg(None, "C: DHCP ERR")
			time.sleep(2)
			self.shutdown()
			return
		self.DHCPleased = True
		
		ipout = subprocess.check_output(['ip','addr','show',self.interface]).split(" ")
		for x in xrange(0,len(ipout)):
			if ipout[x] == "inet":
				ipspl = ipout[x+1].split("/")
				self.main.displayMsg("C: "+self.net, " "+ipspl[0])
				self.wifimenu = self.main.addMenu("  WIFI CONNECT  ","C: "+self.net, self.wifiMenuSelected)
				self.dhcpmenu = self.main.addMenu("  WIFI CONNECT  ", " "+ipspl[0], self.dhcpMenuSelected)
				
	def dhcpMenuSelected(self):
		self.main.displayMsg("  REVOKE DHCP?  ", "     NO/YES     ")
		currentPresses = self.main.selectionPresses
		while (self.main.selectionPresses == currentPresses) and (GPIO.input(22) == 1):
			time.sleep(0.1)
		if self.main.selectionPresses == currentPresses:
			self.shutdownDHCP()

	def wifiMenuSelected(self):
		self.main.displayMsg(" POWEROFF WIFI? ", "     NO/YES     ")
		currentPresses = self.main.selectionPresses
		while (self.main.selectionPresses == currentPresses) and (GPIO.input(22) == 1):
			time.sleep(0.1)
		if self.main.selectionPresses == currentPresses:
			self.shutdown()
			
	def shutdownDHCP(self):
		if self.DHCPleased:
			self.main.displayMsg(" Transmitting: ", "   rDHCP   ")
			subprocess.Popen(["dhclient", "-r", self.interface]).communicate()
			self.DHCPleased = False
			self.main.deleteMenu(self.dhcpmenu)
	
	def shutdown(self):
		self.shutdownDHCP()
		self.main.displayMsg("Disassociating: ", "   "+self.net)
		try:
			self.wpaProcess.terminate()
		except:
			pass
		if self.netConnected:
			self.main.deleteMenu(self.wifimenu)
		self.netConnected = False
