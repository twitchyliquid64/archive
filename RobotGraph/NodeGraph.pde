
int screenx = 680;
int screeny = 680;
int cols = 14;
int rows = 14;

Node[][] nodeGraph = {};
Searcher s = new Searcher();
ArdInterface ard = new ArdInterface(this, screenx, screeny);

void makeGrid(int cols, int rows)//creates the node graph
{
  nodeGraph = new Node[cols][rows];
  for(int x = 1; x <= cols; x++)
    for(int y = 1; y <= rows; y++)
      nodeGraph[x-1][y-1] = new Node(x, y, cols, rows);
}

void setup() {
  size(screenx+300, screeny);
  makeGrid(cols, rows);
}



void draw() {
  background(200);
  int gridx = mouseX/(screenx/cols) + 1;
  int gridy = mouseY/(screeny/cols) + 1;
  
  for(int x = 1; x <= cols; x++)
    for(int y = 1; y <= rows; y++)
      if ((x==gridx) && (y==gridy))
        nodeGraph[x-1][y-1].draw(true);
      else
        nodeGraph[x-1][y-1].draw(false);
  
  s.draw(screenx, screeny, cols, rows);
  
  textSize(16);
  fill(0, 102, 153);
  text("X: " + gridx + " Y: " + gridy, screenx-100, screeny-30);
  
  ard.draw();
}


void keyPressed() {
  
  if (key == CODED) {
    if (keyCode == UP)
      nodeGraph[mouseX/(screenx/cols)][mouseY/(screeny/cols)].toggleNorth();
    if (keyCode == DOWN)
      nodeGraph[mouseX/(screenx/cols)][mouseY/(screeny/cols)].toggleSouth();
    if (keyCode == LEFT)
      nodeGraph[mouseX/(screenx/cols)][mouseY/(screeny/cols)].toggleWest();
    if (keyCode == RIGHT)
      nodeGraph[mouseX/(screenx/cols)][mouseY/(screeny/cols)].toggleEast();
  }else{
    if (key == ' ')
      nodeGraph[mouseX/(screenx/cols)][mouseY/(screeny/cols)].setVehicleLocation();
    if (key == 's')
      s.setStart(mouseX/(screenx/cols), mouseY/(screeny/cols));
    if (key == 'g')
      s.setGoal(mouseX/(screenx/cols), mouseY/(screeny/cols));
    if (key == '\n')
      s.search();
  }
  
  ard.keyPressed();
}
