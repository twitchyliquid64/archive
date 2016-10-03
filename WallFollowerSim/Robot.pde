
class Robot {
  float x,y;
  float ang;
  int wallXPos = 55;
  Robot(int x, int y, float ang, int wallXPos)
  {
    this.x = x;
    this.y = y;
    this.ang = ang;
    this.wallXPos = wallXPos;
  }
  
  void draw()
  {
    //draw robot
    strokeWeight(1);
    stroke(200,200,200);
    fill(200, 50, 50);
    ellipse(x,y,12,12);
    drawFront();
    moveForward(0.1);
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
  
  float calcwallDistance()
  {
    float rad = radians(270+ang);
    float dist = (x-wallXPos) / cos(rad);
    stroke(200, 200, 0);
    if (dist>0 && dist<450)
    {
      drawArrowAtPoint(dist,int(x),int(y),ang+90);//print the line
    }

    return dist;
  }
}
