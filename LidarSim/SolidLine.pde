class SolidLine {
  int x1, x2, y1, y2;
  int lowX, highX, lowY, highY;
  float m,b;
  SolidLine(int x1, int y1, int x2, int y2)
  {
    if (x1 < x2)//store the lowest and highest x value for bounding tests.
    {
      lowX = x1;
      highX = x2;
    }else{
      lowX = x2;
      highX = x1;
    }
    
    if (y1 < y2)//store the lowest and highest y value for bounding tests.
    {
      lowY = y1;
      highY = y2;
    }else{
      lowY = y2;
      highY = y1;
    }
    this.x1 = x1;
    this.x2 = x2;
    this.y1 = y1;
    this.y2 = y2;
    calcPointinterceptValues();
  }
  
  void calcPointinterceptValues()
  {
    this.m = float(this.y2-this.y1)/float(this.x2-this.x1); //gradient
    this.b = (float)this.y1-(this.m*(float)this.x1); //b term in y = mx + b
  }
  
  void draw()
  {
    stroke(255,255,255);
    strokeWeight(2);
    line(x1,y1,x2,y2);
  }
  
  void print()
  {
    println("Line details:");
    println("\tX1:"+x1);
    println("\tX2:"+x2);
    println("\tY1:"+y1);
    println("\tY2:"+y2);
    println("\t B:"+b);
    println("\t M:"+m);
  }
  
  float subX(float x)
  {
    return x*m + b;
  }


  Vec2 intercept(SolidLine l)
  {
    Vec2 inte = _intercept(l);
    
    if (!this.pointOnLine(inte.x,inte.y))return new Vec2(-1,-1);  //check if the point is on the line.
    if (!l.pointOnLine(inte.x,inte.y))return new Vec2(-1,-1);    //check if the point is on the line.
    return inte;
  }

  
  Vec2 _intercept(SolidLine l)
  {
    float xIntercept = (b-l.b) / (l.m - m);
    float yIntercept = xIntercept * l.m + l.b;
    
    if(x1 == x2)//special case - vertical line hence m = infinity
      return new Vec2(x1, l.subX(x1));
    else if(l.x1 == l.x2)//special case - vertical line hence m = infinity
      return new Vec2(l.x1, subX(l.x1));
    else
      return new Vec2(xIntercept, yIntercept);
  }
  
  boolean pointOnLine(float x, float y)//assumes gradient check has been done, we are just checking bounds.
  {
    if(x < lowX)return false;
    if(x > highX)return false;
    if(y < lowY)return false;
    if(y > highY)return false;
    return true;
  }
}

    
    
    
