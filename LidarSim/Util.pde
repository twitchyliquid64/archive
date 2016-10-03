class Vec2 {
  float x, y;
  boolean hasErrored;
  Vec2(float x, float y) {
    this.x = x;
    this.y = y;
  }
}

//holds the result of a ranging operation.
class RangerData {
  float startX, startY;
  Vec2 firstIntercept;
  ArrayList<Vec2> allIntercepts;
  float distance;
  RangerData(float startX, float startY, Vec2 firstIntercept, ArrayList<Vec2> allIntercepts, float distance)
  {
    this.startX = startX;
    this.startY = startY;
    this.firstIntercept = firstIntercept;
    this.allIntercepts = allIntercepts;
    this.distance = distance;
  }
}




//draws an arrow from a point, to an angle and size.
void drawArrowAtPoint(float size, int x, int y, float angle)
{
  Vec2 offset = guage(size, angle);
  arrow(x,y,int(offset.x+x),int(y+offset.y));
}

//given an angle and a distance, returns co-ordinates that project that angle. Starts at (0,0), projecting relative in the direction of the angle.
Vec2 guage(float r, float angle) {
  // Conver from Degree -> Rad
  angle = -angle*(PI/180) ;
    // Convert Polar -> Cartesian
  float x = r * cos(angle);
  float y = r * sin(angle);
 return new Vec2(x,y);
}
 
//draws an arrow from a point to a point.
void arrow(int x1, int y1, int x2, int y2) {
  line(x1, y1, x2, y2);
  pushMatrix();
  translate(x2, y2);
  float a = atan2(x1-x2, y2-y1);
  rotate(a);
  line(0, 0, -10, -10);
  line(0, 0, 10, -10);
  popMatrix();
} 



