int screenX = 800;
int screenY = 600;
int wx = screenX/2;//for the robot sim box
int wy = screenY/2;

  
Robot robot = new Robot(50,wy-50, 0);
CollisionCore engine = new CollisionCore();
ScanPlot plot = new ScanPlot(wx+30, 50, robot);

void setup()
{
  size(800,600);
  newLine(1, wy, wx, wy);
  newLine(1, 1, wx, 1);
  newLine(wx, 1, wx, wy);
  newLine(wx-20, 10, wx-20, wy-10);
  newLine(1, 1, 1, wy);
  newLine(30, 50, 400,400);
  for (int j = 0; j<15; j++){ //make sum random boxes
      int s = (int)random(10, 50);
     int x = (int)random(0, wx-s), y = (int)random(0, wy-s);
     newLine(x, y, x+s, y);
     newLine(x+s, y, x+s, y+s);
     newLine(x+s, y+s, x, y+s);
     newLine(x, y+s, x, y);
  }
}


void draw()
{
  background(0);
  
  engine.draw();
  robot.draw();
  RangerData ranger = robot.rangerDistance(engine, 0);
  line(ranger.startX, ranger.startY, ranger.firstIntercept.x, ranger.firstIntercept.y);
  for(Vec2 pos: ranger.allIntercepts)
    ellipse(pos.x,pos.y,15,15);
  //text("Ranger Distance: " + ranger.distance, wx+30,30);
  plot.draw();
}


void keyPressed()
{
  if (key == CODED)
  {
    if(keyCode == LEFT)
      robot.ang += 2;
    if(keyCode == RIGHT)
      robot.ang -= 2;
    if(keyCode == UP)
      robot.moveForward(2);
    if(keyCode == DOWN)
      robot.moveForward(-2);
  }
}


void newLine(int x1, int y1, int x2, int y2)
{
  engine.RegisterLine(new SolidLine(x1, y1, x2, y2));
}
