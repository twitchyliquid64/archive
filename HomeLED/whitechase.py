#!/usr/bin/python
from sendUDP import send
import time
NUM_PIX = 135

if __name__ == "__main__":
        for i in xrange(NUM_PIX):
		out = [(0,0,0) for z in range(NUM_PIX)]
		out[i] = (255,255,255)
                send(out)
                time.sleep(1.00/24.0/3)
