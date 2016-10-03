import wifiscan
import interfaces
import imp
import lcd
import time
import Network_scanforknown, Network_hostap, Network_jammer
import RPi.GPIO as GPIO
import subprocess

class BlackBox(object):
	VERSION="v0.2 alpha"
	SHORTVERSION="v0.2 A"

	def __init__(self):
		print "[S] Blackbox version "+self.VERSION+" now initialising."
		self.controlinterf = None
                self.display = lcd.lcd(0x27, 1, 1, 8)
		self.msg1 = ""
		self.msg2 = ""
		self.selection = 0
		self.menus = [None]
		self.selectionPresses = 0
		self.inMenu = False
                self.displayMsg("Blackbox "+self.SHORTVERSION, "Initializing...")
                GPIO.setmode(GPIO.BCM)
                GPIO.setup(22, GPIO.IN)
		GPIO.setup(11, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
                time.sleep(2)
                if not GPIO.input(22):#if the button is held - connect to a known network
			self.displayMsg(None, " WIFI CONNECT ")
                	self.connectKnown()
			time.sleep(4)
		self.netmenu = self.addMenu("  WIFI CONFIG  ", "  Change Mode?  ", self.interfaceConfigMenu)
              	self.nominalPrint()
              	
              	#Load the Jammer module
              	self.jammer = Network_jammer.Network_jammer("wlan1", self)
              	
              	#Load the Wardriver module
              	self.wardriver = Network_wardriver.Network_wardriver("wlan1", self)
              	
		GPIO.add_event_detect(11, GPIO.FALLING, callback=self.changeSelectionButtonPress, bouncetime=220)
		self.mainloop()

	def changeSelectionButtonPress(self, pin):
		print "Selectbutton falling", pin
		if not self.inMenu:
			self.selection += 1
			self.selection = self.selection % len(self.menus)
			while (self.selection > 0) and (self.menus[self.selection] == None):
				self.selection += 1
				self.selection = self.selection % len(self.menus)
			self.nominalPrint()
		self.selectionPresses += 1

	def mainloop(self):
		while True:
			GPIO.wait_for_edge(22, GPIO.FALLING)
			print "Mainbutton Falling"
			if not self.handlePwroffHold():
				print "Selection"
				if (self.selection != 0) and (self.menus[self.selection][2] != None):
					#if we are on a menu page run its callback
					self.inMenu = True
					#try:
					self.menus[self.selection][2]()
					#except:
					#	pass
					self.inMenu = False
					self.nominalPrint()


	def handlePwroffHold(self):
                self.displayMsg("   POWER OFF   ", "  Hold for 3")
                starthold = time.time()
                while time.time() < (starthold+3):
                        if GPIO.input(22):#button released
                                self.nominalPrint()
                                return False
                        self.displayMsg(None, "  Hold for "+str(round(2.99-time.time()+starthold,1)))
                        time.sleep(0.1)
                self.displayMsg(" PWR DOWN EVENT", " ")
		subprocess.check_output(["poweroff"])
		return True
		
	def nominalPrint(self):
		if self.selection >= len(self.menus):
			self.selection = 0
		if self.selection == 0:
			self.displayMsg("Blackbox "+self.SHORTVERSION, "WAITING FOR COMM")
		else:
			self.displayMsg(self.menus[self.selection][0], self.menus[self.selection][1])

        def connectKnown(self):
        	self.controlinterf = Network_scanforknown.Network_scanforknown("wlan0", self)

	def addMenu(self, topline,bottomline,cb=None):
		self.menus.append([topline,bottomline,cb])
		return len(self.menus)-1
		
	def updateMenu(self, index,topline,bottomline,cb):
		self.menus[index] = [topline,bottomline,cb]
		
	def deleteMenu(self, index):
		self.menus[index] = None
		while (self.selection > 0) and (self.menus[self.selection] == None):
			self.selection += 1
			self.selection = self.selection % len(self.menus)

	def displayMsg(self, msg1, msg2):
		time.sleep(0.07)
		self.display.lcd_clear()
		if msg1 == None:
			msg1 = self.msg1
		self.display.lcd_puts(msg1, 1)
		if msg2 == None:
			msg2 = self.msg2
		self.display.lcd_puts(msg2, 2)
		self.msg1 = msg1
		self.msg2 = msg2
		self.display.lcd_puts("",3)

	def interfaceConfigMenu(self):
		self.displayMsg("  WHICH MODE?  ", "    AP/SCAN    ")
		currentPresses = self.selectionPresses
		while (self.selectionPresses == currentPresses) and (GPIO.input(22) == 1):
			time.sleep(0.1)
		if self.controlinterf != None:
			self.controlinterf.shutdown()
		if self.selectionPresses == currentPresses:
			self.connectKnown()
		else:
			self.controlinterf = Network_hostap.Network_hostap("wlan0", "blackbox", self)

BlackBox()
