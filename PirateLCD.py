import math
import serial
import time

# Use can use the bus pirates' scan mode to find these.
deviceAddress = [0x4e]

# Corresponds to pins
LCD_BACKLIGHT = 0x08
LCD_ENABLE    = 0b00000100
LCD_RW        = 0b00000010
LCD_REG_SEL   = 0b00000001

# Corresponds to commands/bits
LCD_FUNCTIONSET    = 0x20
LCD_2LINE          = 0x08
LCD_DISPLAYCONTROL = 0x08
LCD_DISPLAYON      = 0x04
LCD_CURSORON       = 0x02
LCD_BLINKON        = 0x01
LCD_SETDDRAMADDR   = 0x80
LCD_CLEARDISPLAY   = 0x01
LCD_RETURNHOME     = 0x02

class PirateI2C(object):
    def __init__ (self, f = "/dev/bus_pirate", vreg = True):
        self.port = serial.Serial(port = f, baudrate = 115200, timeout = 0.01)
        self.resetBP(vreg = vreg)

    def resetBP (self, vreg = False, slow = False):
        self.expected = 9999
        self.clear()
        self.port.write(bytearray([0x0F]))
        while self.port.read(5) != "BBIO1":
            self.clear(9999)
            self.port.write(bytearray([0x00]))
        self.port.write(bytearray([0x02]))
        if self.port.read(4) != "I2C1":
            raise Exception("error initializing bus pirate")

        configurePeripheralsCmd = 0b01000000
        if vreg:
	    configurePeripheralsCmd |= 1<<3 # set bit 3 to high for power
        self.port.write(bytearray([configurePeripheralsCmd]))

        self.port.write(bytearray([0b01100010])) #I2C speed = 100KHz
        self.port.write(bytearray([0x88])) # enable output pins
        self.clear(9999)

    def clear (self, more = 0):
        vals = self.port.read(self.expected + more)
        self.expected = 0
        return vals[-more:]

    def send (self, b):
        if len(b) > 4096:
          raise Exception("Can only send 0-4k bytes")
        self.port.write(bytearray([0x08, (len(b) & 0xff00) >> 8, len(b) & 0xff, 0, 0]))
        self.port.write(bytearray(b))
        ret = self.port.read(1)
        if ord(ret[0]) != 0x01:
            raise Exception("Send failed")

# These things are connected through an I2C port expander.
# We use the expander to 'bitbang' the protocol to the LCD controller.
# There are a bunch of different modes, but we are going to use the 4 bit
# mode.
class PirateLCD(object):
    def __init__(self, f = "/dev/bus_pirate", vreg = True, addr = 0x7e):
        self.p = PirateI2C(f, vreg)
        self.addr = addr
        self.backlight = True
        self.displayOn = True
        self.blink = True
        self.cursor = True

        time.sleep(0.05)
        self._writeExpander(0) # Reset all bits to 0, except backlight

        time.sleep(0.01)
        self._write4bits(3 << 4)
        time.sleep(0.01)
        self._write4bits(3 << 4)
        time.sleep(0.01)
        self._write4bits(3 << 4)
        time.sleep(0.01)
        self._write4bits(2 << 4)

        self.sendByte(LCD_FUNCTIONSET | LCD_2LINE)
        self.update()

    def clear(self):
        self.sendByte(LCD_CLEARDISPLAY)

    def home(self):
        self.sendByte(LCD_RETURNHOME)

    def setCursor(self, v):
        self.sendByte(LCD_SETDDRAMADDR | v)

    def update(self):
        self._displayCmd(0)

    def _displayCmd(self, v):
        if self.displayOn:
            v |= LCD_DISPLAYON
        if self.blink:
            v |= LCD_BLINKON
        if self.cursor:
            v |= LCD_CURSORON
        self.sendByte(LCD_DISPLAYCONTROL | v)

    def _write(self, val):
        self.p.send(self.addr + val)

    def _writeExpander(self, v):
        if self.backlight:
            self._write([v | LCD_BACKLIGHT])
        else:
            self._write([v])

    def write(self, chars):
        for b in chars:
            self.sendByte(ord(b), LCD_REG_SEL)

    def sendByte(self, v, extra=0):
	print "Sending byte: %2x" % v
        self._write4bits((v & 0xf0) | extra)
        self._write4bits(((v & 0x0f) << 4) | extra)

    def _write4bits(self, v):
        self._writeExpander(v)
        self._pulseEnable(v)

    def _pulseEnable(self, v):
        self._writeExpander(v | LCD_ENABLE)
        self._writeExpander(v)


p = PirateLCD("/dev/buspirate", True, deviceAddress)
p.home()
p.write('LOADING')
p.setCursor(64 + 0)
for x in xrange(15):
    p.write('.')
    time.sleep(1)
p.clear()
