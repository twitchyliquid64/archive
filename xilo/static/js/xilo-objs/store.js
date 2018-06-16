// Store is a xilo object which represents a physical tank/store
// Takes a name and an object with options.
// options:
//  - position {x,y}
//  - size - in both width and height
//  - capacity
//  - currentLevel
//  - upNode - points to the node which is connected upstream
//  - downNode - points to the node which is connected downstream
function Store(name, opts) {
  this.name = name || "store";

  this.pos = opts.pos || {x: 0, y: 0};
  this.upNode = opts.up;
  this.downNodes = opts.down || [];

  this.size = opts.size || 60;

  // state
  this.capacity = opts.capacity || 100;
  this.currentLevel = opts.currentLevel || 0;

  this.paintRoot = null;
}

// state-aware API methods
Store.prototype.stateDescription = function(){
  return [
    {
      name: "level",
      type: "NUMBER",
      get: this.getState.bind(this),
      set: this.setState.bind(this),
    },
    {
      name: "capacity",
      type: "NUMBER",
      get: this.getCapacity.bind(this),
    },
  ];
}

// connection-aware API methods
Store.prototype.setUpstreamNode = function(up){
  this.upNode = up;
}
Store.prototype.setDownstreamNode = function(down){
  this.downNodes[this.downNodes.length] = down;
}
Store.prototype.connectionDescription = function(){
  return [
    {
      name: "upstream",
      pos: {x: this.pos.x, y: this.pos.y-this.size},
      wireNormal: 0,
    },
    {
      name: "downstream",
      pos: {x: this.pos.x, y: this.pos.y+(this.size/2)},
      wireNormal: 180,
    },
  ];
}

Store.prototype.bounds = function(){
  return this.paintRoot.bounds;
}

// Xilo integration methods
Store.prototype.addToScene = function(surface){
  if (!this.paintRoot) {
    var half = this.size/2;
    this.top = new surface.Path.Ellipse({
      point: [-half, -half],
      size: [this.size, half],
      fillColor: '#FFFFFF',
    });
    this.middle = new surface.Path.Rectangle({point: [-half, -half+(half/2)], size: [this.size, this.size], fillColor: '#FFFFFF'});
    // this.middle.segments[this.middle.segments.length-2].remove();
    this.bottom = new surface.Path.Ellipse({
      point: [-half, half],
      size: [this.size, half],
      fillColor: '#FFFFFF',
    });

    this.paintRoot = new surface.Group({
      children: [
        this.middle,
        this.bottom,
        this.top,
      ],
      strokeColor: xiloStrokeColor('Store', this.name),
      fillColor: xiloFillColor('Store', this.name),
    });
    this.paintRoot.position = this.pos;
  }
  this.paintRoot.name = this.name;
  surface.project.activeLayer.addChild(this.paintRoot);
}



Store.prototype.setState = function(currentLevel){
  this.currentLevel = currentLevel;
}
Store.prototype.getState = function(){
  return this.currentLevel;
}
Store.prototype.getCapacity = function(){
  return this.capacity;
}
