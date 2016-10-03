#include <Adafruit_NeoPixel.h>
#ifdef __AVR__
  #include <avr/power.h>
#endif

#define PIN 6
Adafruit_NeoPixel strip = Adafruit_NeoPixel(150, PIN, NEO_GRB + NEO_KHZ800);



#include <SPI.h>         // needed for Arduino versions later than 0018
#include <Ethernet.h>
#include <EthernetUdp.h>         // UDP library from: bjoern@cs.stanford.edu 12/30/2008

byte mac[] = { 0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED };
IPAddress ip(192, 168, 0, 177);
unsigned int localPort = 8888;
#define MAX_PKT_SIZE (3*140)+1
unsigned char packetBuffer[MAX_PKT_SIZE]; //incoming packet,
EthernetUDP Udp;




void setup() {
  // start the Ethernet and UDP:
  Ethernet.begin(mac, ip);
  Udp.begin(localPort);

  Serial.begin(115200);
  strip.begin();
  strip.clear();
  strip.show(); // Initialize all pixels to 'off'
}

void loop() {
  // if there's data available, read a packet
  int packetSize = Udp.parsePacket();
  if (packetSize)
  {
    Udp.read(packetBuffer, MAX_PKT_SIZE);

    if((packetBuffer[0] == 7) && (((packetSize-1) % 3) == 0))//RAW COMMAND CODE
    {
      for(int i = 1; i < packetSize; i+=3){
        //printPix(((i-1)/3), packetBuffer[i], packetBuffer[i+1], packetBuffer[i+2]);
        strip.setPixelColor(((i-1)/3), strip.Color(packetBuffer[i], packetBuffer[i+1], packetBuffer[i+2]));
      }
      strip.show();
    }
  }
}



void printPix(unsigned char i, unsigned char r, unsigned char g, unsigned char b){
  Serial.print("Got pixel: [");
  Serial.print(i);
  Serial.print("] (");
  Serial.print(r);
  Serial.print(",");
  Serial.print(g);
  Serial.print(",");
  Serial.print(b);
  Serial.println(")");
}

