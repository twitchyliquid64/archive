//logically encapsulates a robot.
class Robot {
  float x,y;
  float ang;
  Robot(int x, int y, float ang)
  {
    this.x = x;
    this.y = y;
    this.ang = ang;
  }
  
  void draw()
  {
    //draw robot
    strokeWeight(1);
    stroke(200,200,200);
    fill(200, 50, 50);
    ellipse(x,y,12,12);
    drawFront();
  }
  
  void drawFront()
  {
    stroke(0,0,200);
    drawArrowAtPoint(25,int(x),int(y),ang);
  }
  
  void moveForward(float distance)
  {
    float rad = radians(180+ang);
    float deltaY = sin(rad) * distance;
    float deltaX = -cos(rad) * distance;
    x += deltaX;
    y += deltaY;
  }
  
  RangerData rangerDistance(CollisionCore engine, float ang)
  {
    return engine.ranger(ang+this.ang, x, y);
  }
}
