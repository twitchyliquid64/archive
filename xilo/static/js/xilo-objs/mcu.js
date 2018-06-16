// MCU is a xilo object which represents a control unit / decision maker.
// Takes a name and an object with options.
// Unlike other modules, the MCU displays its name on the system diagram.
// options:
//  - position {x,y}
//  - size - in both width and height
//  - downNodes - array of downstream nodes it is connected to
function MCU(name, opts) {
  this.name = name || "door";
  this.pos = opts.pos || {x: 0, y: 0};
  this.downNodes = opts.downNodes || [];
  this.downLinks = opts.downLinks || [];
  this.size = opts.size || 100;

  this.paintRoot = null;
}

// state-aware API methods
MCU.prototype.stateDescription = function(){
  return [];
}

// connection-aware API methods
MCU.prototype.setUpstreamNode = function(up, upLink){
  this.upNode = up;
  this.upLink = upLink;
}
MCU.prototype.setDownstreamNode = function(down, downLink){
  this.downNodes[this.downNodes.length] = down;
  this.downLinks[this.downLinks.length] = downLink || null;
}
MCU.prototype.connectionDescription = function(){
  return [];
}

MCU.prototype.bounds = function(){
  return this.paintRoot.bounds;
}

MCU.prototype.addToScene = function(surface){
  if (!this.paintRoot) {
    var half = this.size/2;
    var bound2 = this.size*2.5/3;
    this.outline = new surface.Path.Rectangle({point: [0, 0], size: [this.size, this.size], radius: this.size/20, fillColor: xiloFillColor('MCU', this.name)});
    this.inside = new surface.Path.Rectangle({center: [half, half], size: [bound2, bound2], radius: this.size/400});
    this.left = new surface.Path.Rectangle({point: [0, half-(this.size-bound2)/4], size: [(this.size-bound2)/2, (this.size-bound2)/2]});
    this.right = new surface.Path.Rectangle({point: [this.size-(this.size-bound2)/2, half-(this.size-bound2)/4], size: [(this.size-bound2)/2, (this.size-bound2)/2]});
    this.top = new surface.Path.Rectangle({point: [this.size/2-(this.size-bound2)/4, 0], size: [(this.size-bound2)/2, (this.size-bound2)/2]});
    this.bottom = new surface.Path.Rectangle({point: [this.size/2-(this.size-bound2)/4, this.size-(this.size-bound2)/2], size: [(this.size-bound2)/2, (this.size-bound2)/2]});

    this.text = new surface.PointText(new surface.Point(this.size/2, this.size/2+(this.size-bound2)/4));
    this.text.justification = 'center';
    this.text.strokeColor = xiloStrokeColor('MCU', this.name);
    this.text.fillColor = xiloStrokeColor('MCU', this.name);
    this.text.strokeWidth = 0.1;
    this.text.content = this.name;
    this.text.scale(1.7);

    this.paintRoot = new surface.Group({
      children: [
        this.outline,
        this.inside,
        this.left,
        this.right,
        this.top,
        this.bottom,
        this.text,
      ],
      strokeColor: xiloStrokeColor('MCU', this.name),
    });
    this.paintRoot.position = this.pos;
  }
  this.paintRoot.name = this.name;
  surface.project.activeLayer.addChild(this.paintRoot);
}
