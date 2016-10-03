import uuid, subprocess, time

class Network_hostap(object):
	def __init__(self, interface, name, main):
		self.interface = interface
		self.name = name
		self.hostapProc = None
		self.main = main
		self.menu = self.main.addMenu("AP: "+name, "PW: "+self.startRandPass())

	def startRandPass(self):
		self.password = str(uuid.uuid4().get_hex().upper()[0:9])
		config = self.createConfig(self.interface, self.name, self.password)
		fo = open("/hostap.conf", "w+")
                fo.write( config )
                fo.close()
		self.hostapProc = subprocess.Popen(["hostapd", "-B", "/hostap.conf"])
		time.sleep(3)
		subprocess.check_output(["ifconfig", self.interface, "192.168.2.198", "netmask", "255.255.255.0", "up"])
		time.sleep(3)
		subprocess.check_output(["udhcpd"])
		return self.password

	def shutdown(self):
		subprocess.check_output(["killall", "udhcpd"])
		time.sleep(1)
		self.hostapProc.terminate()
		self.main.deleteMenu(self.menu)

	def createConfig(self, interface, name, password):
		return """ctrl_interface=/var/run/hostapd
interface=%s
driver=rtl871xdrv
country_code=AU
ctrl_interface_group=0
ssid=%s
hw_mode=g
channel=1
wpa=3
wpa_passphrase=%s
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP
beacon_int=100
auth_algs=3
macaddr_acl=0
wmm_enabled=1
eap_reauth_period=360000000""" % (interface, name, password)
