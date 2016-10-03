
//responsible for:
//1. Printing search path and start/end positions on the screen
//2. Actually calculating the search paths
class Searcher {
  int startX, startY, goalX, goalY;
  int pathAnimCounter, pathAnimCounterDelay;
  boolean startSet, goalSet, showPath;
  ArrayList<Node> path = new ArrayList<Node>();
  
  void draw(int screenx, int screeny, int cols, int rows)//draws start/end, as well as animates the search path
  {
    float quantax = screenx / cols;
    float quantay = screeny / rows;
    if (startSet)
    {
      fill(200, 200, 50);
      ellipse(quantax * startX + (quantax/2), quantay * startY + (quantay/2), quantax/4, quantay/4);
    }
    if (goalSet)
    {
      fill(50, 200, 50);
      ellipse(quantax * goalX + (quantax/2), quantay * goalY + (quantay/2), quantax/4, quantay/4);
    }
    
    if(showPath)
    {
      pathAnimCounterDelay++;
      pathAnimCounterDelay = pathAnimCounterDelay % 5;
      if (pathAnimCounterDelay == 0)
        pathAnimCounter++;
      pathAnimCounter = pathAnimCounter % path.size();
      Node n = path.get(pathAnimCounter);
        fill(200, 200, 200, 150);
        ellipse(quantax * n.x + (quantax/2), quantay * n.y + (quantay/2), quantax/4, quantay/4);
    }
  }
  
  void setStart(int x, int y)
  {
    if ((x==startX)&&(y==startY)&&(startSet))
      this.startSet = false;
    else
      this.startSet = true;
    this.startX = x;
    this.startY = y;
    showPath = false;
  }
  
  void setGoal(int x, int y)
  {
    if ((x==goalX)&&(y==goalY)&&(goalSet))
      this.goalSet = false;
    else
      this.goalSet = true;
    this.goalX = x;
    this.goalY = y;
    showPath = false;
  }
  
  void search()
  {
    println("Commencing search ... Please wait.");
    nodeGraph[startX][startY].resetSearches();
    
    ResultsSet results = recursiveSearch(nodeGraph[startX][startY], new ResultsSet(false, new ArrayList<Node>()));
    
    if (!results.successful)
      println("Search failed.");
    else{
      println("Search success!");
      this.path = results.path;
      showPath = true;
      for(Node point: results.path)
        println(point);
    }
  }
  
  ResultsSet recursiveSearch(Node current, ResultsSet ret)
  {
    if(current == null)return new ResultsSet(false, null);
    if (current.hasBeenSearched)return new ResultsSet(false, null);
    current.hasBeenSearched = true;

    if(current == nodeGraph[goalX][goalY])return ret.add(current);
    
    ResultsSet temp;
    
    if(!current.wallNorth)
      if ((temp = recursiveSearch(current.getNorth(), ret)).successful)
        return temp.add(current);
    if(!current.wallSouth)
      if ((temp = recursiveSearch(current.getSouth(), ret)).successful)
        return temp.add(current);
    if(!current.wallWest)
      if ((temp = recursiveSearch(current.getWest(), ret)).successful)
        return temp.add(current);
    if(!current.wallEast)
      if ((temp = recursiveSearch(current.getEast(), ret)).successful)
        return temp.add(current);

    return new ResultsSet(false, null);
  }
}


class ResultsSet{//used to propergate the results back up the recursive function.
  boolean successful;
  ArrayList<Node> path;
  ResultsSet(boolean successful, ArrayList<Node> path)
  {
    this.successful = successful;
    this.path = path;
  }
  ResultsSet add(Node node)
  {
    this.successful = true;
    this.path.add(node);
    return this;
  }
}
