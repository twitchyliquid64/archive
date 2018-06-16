// Valve is a xilo object which represents a physical valve/relay/connector.
// Takes a name and an object with options.
// options:
//  - position {x,y}
//  - isVertical (bool)
//  - size - radius of the valve symbol
//  - open - if the valve is open or not. Default closed.
//  - upNode - points to the node which is connected upstream
//  - downNode - points to the node which is connected downstream
//  - animationSpeedSeconds - seconds the valve takes to open - default 0.5.
function Valve(name, opts) {
  this.name = name || "valve";
  this.isVertical = !!opts.isVertical;

  this.pos = opts.pos || {x: 0, y: 0};
  this.upNode = opts.upNode;
  this.downNode = opts.downNode;

  this.size = opts.size || 12; // radius

  // state + animation
  this.closed = !opts.open;
  this.isAnimating = false;
  this.animationStartTime = null;
  this.animationSpeedSeconds = opts.animationSpeedSeconds || 0.5;

  this.paintRoot = null;
  this.rotation = 0;
}

// state-aware API methods
Valve.prototype.stateDescription = function(){
  return [
    {
      name: "open",
      type: "BOOL",
      get: this.getState.bind(this),
      set: this.setState.bind(this),
    },
  ];
}

// connection-aware API methods
Valve.prototype.setUpstreamNode = function(up){
  this.upNode = up;
}
Valve.prototype.setDownstreamNode = function(down){
  this.downNode = down;
}
Valve.prototype.connectionDescription = function(){
  if (!isVertical) { //horizontal
    return [
      {
        name: "upstream",
        pos: {x: this.pos.x-this.size, y: this.pos.y},
        wireNormal: 90,
      },
      {
        name: "downstream",
        pos: {x: this.pos.x+this.size, y: this.pos.y},
        wireNormal: -90,
      },
    ];
  } else {
    return [
      {
        name: "upstream",
        pos: {x: this.pos.x, y: this.pos.y-this.size},
        wireNormal: 90,
      },
      {
        name: "downstream",
        pos: {x: this.pos.x, y: this.pos.y+this.size},
        wireNormal: -90,
      },
    ];
  }
}

Valve.prototype.bounds = function(){
  return this.paintRoot.bounds;
}

// Xilo integration methods
Valve.prototype.addToScene = function(surface){
  if (!this.paintRoot) {
    this.circ = new surface.Path.Circle({radius: this.size, strokeWidth: 1});
    this.line = new surface.Path({segments: [{x: -this.size, y: 0}, {x: this.size, y: 0}], strokeWidth: 1});
    this.paintRoot = new surface.Group({
  		children: [
  			this.circ,
        this.line,
  		],
  		strokeColor: xiloStrokeColor('Valve', this.name),
      fillColor: xiloFillColor('Valve', this.name),
  	});
    this.paintRoot.position = this.pos;
    this.paintRoot.onFrame = this.onFrame.bind(this);
  }
  this._rotate(this.closed ? 90 : 0);
  this.paintRoot.name = this.name;
  surface.project.activeLayer.addChild(this.paintRoot);
}


Valve.prototype.toggleState = function(){
  this.closed = !this.closed;
  this.startAnimating = true;
}

Valve.prototype.setState = function(state){
  if (!state == this.closed){
    return;
  }
  this.toggleState();
}
Valve.prototype.getState = function(){
  return !this.closed;
}

// internal method to rotate to a absolute value - we store the state because the paper.js API only supports absolute rotation
Valve.prototype._rotate = function(amt){
  amt += this.isVertical ? 90 : 0;
  this.line.rotate(amt - this.rotation);
  this.rotation = amt;
}

// implements animation of state change
Valve.prototype.onFrame = function(event){
  if (this.startAnimating) {
    this.startAnimating = false;
    this.animationStartTime = event.time;
    this.isAnimating = true;
    return;
  }

  if (this.isAnimating) {
    if ((this.animationStartTime+this.animationSpeedSeconds) < event.time) {
      this.isAnimating = false;
      this._rotate(this.closed ? 90 : 0);
    } else {
      if (this.closed) {
        this._rotate( 90 * ((event.time-this.animationStartTime)/this.animationSpeedSeconds));
      } else {
        this._rotate(90 - (90 * ((event.time-this.animationStartTime)/this.animationSpeedSeconds)));
      }
    }
  }
}
