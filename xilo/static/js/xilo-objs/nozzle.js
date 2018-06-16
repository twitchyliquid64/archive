// Nozzle is a xilo object which represents a physical nozzle.
// Takes a name and an object with options.
// options:
//  - position {x,y}
//  - isVertical (bool)
//  - size - radius of the valve symbol
//  - isMetered - enable flow animations
//  - upNode - points to the node which is connected upstream
function Nozzle(name, opts) {
  this.name = name || "nozzle";
  this.isVertical = !!opts.isVertical;

  this.pos = opts.pos || {x: 0, y: 0};
  this.upNode = opts.up;

  this.size = opts.size || 20;
  this.isMetered = opts.isMetered || false;
  this.flowRate = opts.flowRate || 0;

  this.paintRoot = null;
}

// state-aware API methods
Nozzle.prototype.stateDescription = function(){
  return [];
}

// connection-aware API methods
Nozzle.prototype.setUpstreamNode = function(up){
  this.upNode = up;
}
Nozzle.prototype.connectionDescription = function(){
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

Nozzle.prototype.bounds = function(){
  return this.paintRoot.bounds;
}

// Xilo integration methods
Nozzle.prototype.addToScene = function(surface){
  if (!this.paintRoot) {
    var l = this.size/1.3;
    var q = this.size/2.5;
    var t = this.size/5;

    this.outline = new surface.Path({segments: [
      {x: -q, y: 0},
      {x: q, y: 0},
      {x: q, y: l},
      {x: q-t, y: l+t},
      {x: q-t, y: l+t*2},
      {x: t-q, y: l+t*2},
      {x: t-q, y: l+t},
      {x: -q, y: l},
      {x: -q, y: 0},
    ], strokeWidth: 1});
    this.paintRoot = new surface.Group({
  		children: [
        this.outline,
  		],
  		strokeColor: xiloStrokeColor('Nozzle', this.name),
      fillColor: xiloFillColor('Nozzle', this.name),
  	});

    if (this.isMetered) {
      var c = new surface.Path.Circle({radius: q*4/5});
      this.meter = new surface.Group({
        children: [
          c,
          new surface.Path({segments: [c.getPointAt(c.length*1/8), c.getPointAt(c.length*5/8)], strokeWidth: 1}),
          new surface.Path({segments: [c.getPointAt(c.length*3/8), c.getPointAt(c.length*7/8)], strokeWidth: 1}),
        ],
        strokeColor: xiloStrokeColor('Nozzle.Meter', this.name),
        fillColor: xiloFillColor('Nozzle.Meter', this.name),
      });
      this.paintRoot.addChild(this.meter);
      this.paintRoot.onFrame = this.onFrame.bind(this);
    }

    this.paintRoot.position = this.pos;
  }
  this.paintRoot.name = this.name;
  surface.project.activeLayer.addChild(this.paintRoot);
}

Nozzle.prototype.onFrame = function(event){
  if (this.flowRate != 0) {
    this.meter.rotate(this.flowRate);
  }
}
