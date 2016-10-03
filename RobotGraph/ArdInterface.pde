import processing.serial.*;

int MODE_SETUP = 0;
int MODE_RUN = 1;

class ArdInterface {
  
  int x, y;
  Serial serial;
  int mode = MODE_SETUP;
  NodeGraph n;
  
  ArdInterface(NodeGraph n, int x, int y)
  {
    this.x = x;
    this.y = y;
    this.n = n;
  }
  
  void draw()
  {
    textSize(16);
    fill(0, 0, 0);

    if(mode == MODE_SETUP)
      text("Arduino Status: SETUP", x+80, y/20);
    else if(mode == MODE_RUN)
      text("Arduino Status: ONLINE", x+80, y/20);
  }
  
  
  void keyPressed()
  {
    if (keyCode == SHIFT)
      if (mode == MODE_SETUP)
      {
         serial = new Serial(this.n, Serial.list()[0], 9600);
         mode = MODE_RUN;
      }
  }
}
