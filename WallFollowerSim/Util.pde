class Vec2 {
  float x, y;
  Vec2(float x, float y) {
    this.x = x;
    this.y = y;
  }
}


void drawArrowAtPoint(float size, int x, int y, float angle)
{
  Vec2 offset = guage(size, angle);
  arrow(x,y,int(offset.x+x),int(y+offset.y));
}

Vec2 guage(float r, float angle) {
  // Conver from Degree -> Rad
  angle = -angle*(PI/180) ;
    // Convert Polar -> Cartesian
  float x = r * cos(angle);
  float y = r * sin(angle);
 return new Vec2(x,y);
}
 
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
