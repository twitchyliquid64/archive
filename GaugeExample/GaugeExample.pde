import processing.serial.*;
/**a
a simple gauge simulator by <a href="http://www.local-guru.net/blog">guru</a>
*/

Serial serial;

void setup() {
  size(400,300);
  serial = new Serial(this, "/dev/ttyACM0", 9600);
  //smooth();
}

int[] val = new int[255];
boolean mode = false;
void serialEvent(Serial port) {
  // Data from the Serial port is read in serialEvent() using the readStringUntil() function with * as the end character.
  String input = port.readStringUntil('\n');
  if (input != null)
  {
    println(input);
    String[] vals = (splitTokens(trim(input), ",*"));
    println(vals);
    int[] val2 = new int[vals.length+1];
    for(int i=0; i<val2.length;i++)
      val2 = int(vals);
    val[val2[0]] = val2[1];
  }
}

void keyPressed()
{
  if (!mode)
  {
    mode = true;
    serial.write("SSS");
  }
}

void draw() {
  background(255);
  pushMatrix();
  translate( width/2,height/2-20);
  drawGauge( map( val[1], 0,175,0,1 ));  
  popMatrix();

  pushMatrix();
  translate( width/2,height-20);
  drawGauge( map( val[0], 0,175,0,1 ));  
  popMatrix();

  text("Online: "+mode, 10, 10);
  for(int i=0;i<6;i++)
    text("Channel "+i+": "+val[i], 10, 22 + (12*i));
}

void drawGauge( float val ) {
  stroke(0);

  for( int i=0; i<31; i++) {
    float a = radians(210 + i * 120 / 30);
    float r1 = 100;
    float r2 = 90;
    r2 = i % 5 == 0 ? 85 : r2;
    r2 = i % 10 == 0 ? 80 : r2;

    line( r1*cos(a), r1*sin(a), r2*cos(a), r2*sin(a));
  }
  stroke( 255,0,0 );
  float b = radians( 210 + val * 120 );
  fill(255,0,0);
  ellipse(0,0,10,10);
  line( -10*cos(b),-10*sin(b),100 * cos(b), 100 * sin(b));
}
