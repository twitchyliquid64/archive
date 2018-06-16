// Door is a xilo object which represents a lockable door.
// Takes a name and an object with options.
// options:
//  - position {x,y}
//  - size - in both width and height
//  - locked - boolean symbolising the initial state of the door.
//  - upNode - points to the node which is connected upstream
//  - animationSpeedSeconds - seconds the door takes to open - default 0.3.
function Door(name, opts) {
  this.name = name || "door";

  this.pos = opts.pos || {x: 0, y: 0};
  this.upNode = opts.up;

  this.size = opts.size || 100;

  // state
  this.locked = opts.locked || false;
  this.isAnimating = false;
  this.animationStartTime = null;
  this.animationSpeedSeconds = opts.animationSpeedSeconds || 0.3;

  this.paintRoot = null;
}

// state-aware API methods
Door.prototype.stateDescription = function(){
  return [
    {
      name: "locked",
      type: "NUMBER",
      get: this.getState.bind(this),
      set: this.setState.bind(this),
    },
  ];
}

// connection-aware API methods
Door.prototype.setUpstreamNode = function(up, upLink){
  this.upNode = up;
  this.upLink = upLink;
}

Door.prototype.bounds = function(){
  return this.paintRoot.bounds;
}

// Xilo integration methods
Door.prototype.addToScene = function(surface){
  if (!this.paintRoot) {
    var half = this.size/2;
    this.outline = new surface.Path.Rectangle({point: [0, 0], size: [this.size/2, this.size*4/5], radius: this.size/30});
    this.knob = new surface.Path.Circle({center: [half*2/8, half*4/5], radius: this.size/25, fillColor: xiloFillColor('Door', this.name)})
    this.lockOutlineTop = new surface.Path.Rectangle({point: [-this.size/30, half/10], size: [half/13, this.size/5], radius: this.size/50});
    this.lockOutlineBottom = new surface.Path.Rectangle({point: [-this.size/30, half+(half/10)], size: [half/13, this.size/5], radius: this.size/50});
    this.frame = new surface.Path.Rectangle({point: [-this.size/12, 0], size: [5, this.size*4/5], radius: this.size/50});

    this.lock = new surface.Group({
      children: [
        this.lockOutlineTop,
        this.lockOutlineBottom,
      ]
    });

    this.paintRoot = new surface.Group({
      children: [
        this.outline,
        this.knob,
        this.lock,
        this.frame,
      ],
      strokeColor: xiloStrokeColor('Door', this.name),
      fillColor: xiloFillColor('Door', this.name),
    });
    this.lock.fillColor = lerpColor('#33BB33', '#BB3333', this.locked);
    this.lock.strokeColor = lerpColor('#33BB33', '#BB3333', this.locked);
    this.paintRoot.onFrame = this.onFrame.bind(this);
  }
  this.paintRoot.position = this.pos;
  this.paintRoot.name = this.name;
  surface.project.activeLayer.addChild(this.paintRoot);
}


Door.prototype.setState = function(locked){
  if (locked != this.locked)
    this.startAnimating = true;
  this.locked = !!locked;
}
Door.prototype.getState = function(){
  return this.locked;
}

// implements animation of state change
Door.prototype.onFrame = function(event){
  if (this.startAnimating) {
    this.startAnimating = false;
    this.animationStartTime = event.time;
    this.isAnimating = true;
    return;
  }

  if (this.isAnimating) {
    if ((this.animationStartTime+this.animationSpeedSeconds) < event.time) {
      this.isAnimating = false;
      this.lock.fillColor = this.lock.strokeColor = lerpColor('#33BB33', '#BB3333', this.locked);
    } else {
      this.lock.fillColor = this.lock.strokeColor = lerpColor('#33BB33', '#BB3333', this.locked ? ((event.time-this.animationStartTime)/this.animationSpeedSeconds) : (1-((event.time-this.animationStartTime)/this.animationSpeedSeconds)));
    }
  }
}
