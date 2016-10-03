//this class keeps track of all interceptable objects, and does all collision testing, intercepts, ray casting etc. Basically responsible for all physics.
class CollisionCore {
  ArrayList<SolidLine> lines;
  int RANGE = 8000;
  CollisionCore()
  {
    lines = new ArrayList<SolidLine>();
  }
  
  public void RegisterLine(SolidLine line)
  {
    println("Line " + lines.size() + " registered.");
    lines.add(line);
  }
  
  public void draw()
  {
    for(SolidLine s: lines)
      s.draw();
  }
  
  //given an angle and a starting point, ray cast from that point finding all collisions/intercepts with objects.
  public RangerData ranger(float ang, float x, float y)
  {
    Vec2 lineEnd = guage(RANGE, ang);//cast the ray which we will test the intercept along
    lineEnd.x += x;//relative co-ords, make them absolute
    lineEnd.y += y;//relative co-ords, make them absolute
    
    SolidLine ray = new SolidLine(int(x),int(y),int(lineEnd.x),int(lineEnd.y));//create the line to test the intercept along 
    float shortest = 99999999;//keeps track of the shortest distance seen for an intercept - the shortest is returned as the range.
    Vec2 inter = null;//keeps track of the intercept-coords which the ray intercepted.
    ArrayList<Vec2> allIntercepts = new ArrayList<Vec2>();
    
    for(SolidLine l: lines)//iterate all lines it could collide with and test their distance. The closest one is the one it collided with.
    {
      Vec2 intercept = ray.intercept(l);
      if (intercept.x != -1)
      {
        float dist = sqrt(sq(intercept.x-x)+sq(intercept.y-y));//find the distance between the intercept and the robot.
        if(dist < shortest)//if the distance is the shortest so far, store it.
        {
          shortest = dist;
          inter = intercept;
        }
        allIntercepts.add(intercept);
      }
    }
    return new RangerData(x,y,inter,allIntercepts,shortest);
  }
}
