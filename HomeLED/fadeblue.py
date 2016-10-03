#!/usr/bin/python
from sendUDP import send
import time

if __name__ == "__main__":
        for i in xrange(100):
                send([(0,0,int(255*(i/100.0))) for x in range(135)])
                time.sleep(0.01)
