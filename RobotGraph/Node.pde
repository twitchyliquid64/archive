//represents a single grid square, which may have 4 corners and possibly contain a vehicle.
class Node {
  int x,y,cols,rows;
  boolean containsVehicle, wallNorth, wallSouth, wallEast, wallWest;
  boolean hasBeenSearched;//used by searcher to determine if a node has been searched or not.
  Node previous;//used by searcher to retrace the path once search complete.
  
  Node(int x, int y, int cols, int rows)
  {
    this.x = x-1;
    this.y = y-1;
    this.cols = cols;//to be able to know our bounds on the screen
    this.rows = rows;//to be able to know our bounds on the screen
  }
  
  public void draw(boolean isActive)
  {
    if (isActive)
      strokeWeight(3);
    else
      strokeWeight(1);
      
    float quantax = screenx / cols;
    float quantay = screeny / rows;
    stroke(0,0,0);
    line(quantax * x, quantay * y,quantax * x + quantax, quantay * y);
    line(quantax * x, quantay * y,quantax * x, quantay * y + quantay);
    if (isActive)
    {
      line(quantax * x, quantay * y + quantay,quantax * x + quantax, quantay * y + quantay);
      line(quantax * x + quantax, quantay * y,quantax * x + quantax, quantay * y + quantay);
    }
    if(containsVehicle)
    {
      fill(200, 50, 50);
      ellipse(quantax * x + (quantax/2), quantay * y + (quantay/2), quantax/2, quantay/2);
    }
    strokeWeight(2);
    stroke(100,0,0);
    if(wallNorth)
      line(quantax * x, quantay * y+3,quantax * x + quantax, quantay * y+3);
    if(wallSouth)
      line(quantax * x, quantay * y + quantay-3,quantax * x + quantax, quantay * y + quantay-3);
    if(wallEast)
      line(quantax * x + quantax-3, quantay * y,quantax * x + quantax-3, quantay * y + quantay);
    if(wallWest)
      line(quantax * x +3, quantay * y,quantax * x +3, quantay * y + quantay);
  }
  
  public void setVehicleLocation()
  {
    for(int x = 1; x <= cols; x++)
      for(int y = 1; y <= rows; y++)
        nodeGraph[x-1][y-1].containsVehicle = false;
    this.containsVehicle = true;
  }
  
  public void resetSearches()
  {
    for(int x = 1; x <= cols; x++)
      for(int y = 1; y <= rows; y++)
      {
        nodeGraph[x-1][y-1].hasBeenSearched = false;
        nodeGraph[x-1][y-1].previous = null;
      }
  }
  
  public void toggleNorth()
  {
    this.wallNorth = !this.wallNorth;
    if(this.withinBounds(this.x, this.y-1))
      nodeGraph[this.x][this.y-1].wallSouth = this.wallNorth;
  }
  
  public void toggleSouth()
  {
    this.wallSouth = !this.wallSouth;
    if(this.withinBounds(this.x, this.y+1))
      nodeGraph[this.x][this.y+1].wallNorth = this.wallSouth;
  }
  
  public void toggleEast()
  {
    this.wallEast = !this.wallEast;
    if(this.withinBounds(this.x+1, this.y))
      nodeGraph[this.x+1][this.y].wallWest = this.wallEast;
  }
  
  public void toggleWest()
  {
    this.wallWest = !this.wallWest;
    if(this.withinBounds(this.x-1, this.y))
      nodeGraph[this.x-1][this.y].wallEast = this.wallWest;
  }
  
  private boolean withinBounds(int x, int y)//checks if a given relative co-ordinate is within bounds of our graph
  {
    //println("Wanted - X: "+x+" Y: "+y);
    if(x < 0)return false;
    if(y < 0)return false;
    if(y >= nodeGraph[this.x].length)return false;
    if(x >= nodeGraph.length)return false;
    return true;
  }
  
  
  public Node getNorth()
  {
    if(this.withinBounds(this.x, this.y-1))
      return nodeGraph[this.x][this.y-1];
    return null;
  }
  
  public Node getSouth()
  {
    if(this.withinBounds(this.x, this.y+1))
      return nodeGraph[this.x][this.y+1];
    return null;
  }
  
  public Node getEast()
  {
    if(this.withinBounds(this.x+1, this.y))
      return nodeGraph[this.x+1][this.y];
    return null;
  }
  
  public Node getWest()
  {
    if(this.withinBounds(this.x-1, this.y))
      return nodeGraph[this.x-1][this.y];
    return null;
  }
  
  public String toString()
  {
    return "("+x+","+y+") - "+wallNorth+" "+wallEast+" "+wallSouth+" "+wallWest;
  }
}

