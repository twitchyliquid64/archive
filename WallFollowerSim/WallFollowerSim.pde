
Robot robot;
Controller controller;
boolean pause = true;

void setup()
{
  size(600,600);
  robot = new Robot(130, 450, 130, 55);
  controller = new Controller(robot, 40);
}


void draw()
{
  if(pause)
  {
    background(0);
    drawWall();
    robot.draw();
    controller.process();
  }

}

void keyPressed()
{
  pause = !pause;
}

void drawWall()
{
  strokeWeight(5);
  stroke(255,255,255);
  line(55,25,50,550);
}
