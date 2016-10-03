#!/usr/bin/python
import socket
import time

UDP_IP = "192.168.0.177"
UDP_PORT = 8888

def genPixPkt(lst):
	out = ""
	for pix in lst:
		out += chr(pix[0])
		out += chr(pix[1])
		out += chr(pix[2])
	return out

sock = socket.socket(socket.AF_INET, # Internet
                     socket.SOCK_DGRAM) # UDP

def send(pix):
	sock.sendto("\x07" + genPixPkt(pix), (UDP_IP, UDP_PORT))


if __name__ == "__main__":
	for i in xrange(100):
		sock.sendto("\x07" + genPixPkt([(0,0,int(100*(i/100.0))) for x in range(135)]), (UDP_IP, UDP_PORT))
		time.sleep(0.01)
