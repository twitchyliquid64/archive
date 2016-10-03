//Responsible for doing a single LiDAR sweep and plotting it.
class ScanPlot {
  int x,y;
  int sizeX,sizeY;
  int plotCentreX, plotCentreY;
  PGraphics plotBuffer;
  float distance[];
  Robot robot;
  
  ScanPlot(int x, int y, Robot robot)
  {
    this.x = x;
    this.y = y;
    this.sizeX = 250;
    this.sizeY = 250;
    this.robot = robot;
    
    this.plotCentreX = x+(sizeX/2);
    this.plotCentreY = y+(sizeY/2);
  }
  
  void draw()
  {
    //nasty hack as createGraphics doesnt work in a constructor
    if(plotBuffer == null)plotBuffer = createGraphics(sizeX,sizeY);
    
    //plot scan box
    fill(0,0,0);
    stroke(255,255,255);
    rect(x,y,sizeX,sizeY);
    
    //write "LIDAR"
    fill(255,255,255);
    text("LiDAR", plotCentreX - textWidth("LiDAR")/2, y - 20);
    
    //update
    sweep();
    
    //draw our plot to the window
    image(plotBuffer,x+2,y+2, sizeX-4, sizeY-4);
    
    //draw direction arrow
    stroke(0,0,255);
    drawArrowAtPoint(25,plotCentreX,plotCentreY,90);
  }
  
  void sweep()//updates our plot
  {
    //nasty hack as createGraphics doesnt work in a constructor
    if(plotBuffer == null)plotBuffer = createGraphics(sizeX,sizeY);
    
    //setup our back buffer
    plotBuffer.beginDraw();
    plotBuffer.background(0);
    plotBuffer.stroke(255,255,255);
    
    //get a ranger value for each degree, and plot it to our back buffer.
    distance = new float[360];
    for(int i=0;i < 360;i++)
    {
      RangerData ranger = robot.rangerDistance(engine, i);
      distance[i] = ranger.distance;
      Vec2 lineEnd = guage(distance[i] / 3, i+90);
      lineEnd.x += sizeX / 2;//relative -> absolute
      lineEnd.y += sizeY / 2;
      plotBuffer.line(sizeX / 2,sizeY / 2,lineEnd.x,lineEnd.y);
    }
    plotBuffer.endDraw();
  }
}
