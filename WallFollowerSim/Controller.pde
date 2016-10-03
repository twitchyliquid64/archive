class Controller {
  int desiredDistance = 50;
  float lastDistance;
  
  float sizeCoeff = 0.1;
  
  RollingAverage RdeltaDistance = new RollingAverage();
  
  Robot robot;
  Controller(Robot r, int desiredDistance)
  {
    this.robot = r;
    this.desiredDistance = desiredDistance;
  }
  
  void process()
  {
    float wallDistance = robot.calcwallDistance();
    float deltaDistance = RdeltaDistance.ad(lastDistance-wallDistance);
    float error = 0;
    lastDistance = wallDistance;
    
    float multiplier = 0;
    if ((wallDistance > desiredDistance) && (deltaDistance<0))//then deltadistance should be positive
      multiplier = -1;
    else
      multiplier = 1;
    
    robot.ang += multiplier;

    if (wallDistance>0 && wallDistance<450)
    {
      fill(255,255,255);
      text("Wall Distance: "+wallDistance, 400, 50);
      //text("Error: "+error, 400, 70);
      text("Delta Distance: "+deltaDistance, 400, 90);
      //text("Wanted Delta Distance: "+wantedDeltaDistance, 400, 110);
      //text("P term: "+pTerm, 400, 130);
      text("Angle: "+robot.ang, 400, 150);
    }
  }
}
