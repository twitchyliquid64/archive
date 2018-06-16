function Sensor(name, opts) {
  this.name = name || "sensor";
  this.pos = opts.pos || {x: 0, y: 0};
  this.size = opts.size || 2;

  this.upNode = opts.upNode;
  this.upLink = opts.upLink;

  this.showWifi = opts.showWifi || false;
  this.showButton = opts.showButton || true;

  this.paintRoot = null;
}

// state-aware API methods
Sensor.prototype.stateDescription = function(){
  return [];
}

// connection-aware API methods
Sensor.prototype.setUpstreamNode = function(up, upLink){
  this.upNode = up;
  this.upLink = upLink;
}
Sensor.prototype.setDownstreamNode = function(down, downLink){
}
Sensor.prototype.connectionDescription = function(){
  return [];
}

Sensor.prototype.bounds = function(){
  var bound = this.paintRoot.children[1].bounds.clone();
  return bound;
}

var sensorSvg = '\
<svg height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg">\
    <path d="M15 9H9c-.55 0-1 .45-1 1v12c0 .55.45 1 1 1h6c.55 0 1-.45 1-1V10c0-.55-.45-1-1-1zm-3 6c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zM7.05 6.05l1.41 1.41C9.37 6.56 10.62 6 12 6s2.63.56 3.54 1.46l1.41-1.41C15.68 4.78 13.93 4 12 4s-3.68.78-4.95 2.05zM12 0C8.96 0 6.21 1.23 4.22 3.22l1.41 1.41C7.26 3.01 9.51 2 12 2s4.74 1.01 6.36 2.64l1.41-1.41C17.79 1.23 15.04 0 12 0z"/>\
</svg>\
'

Sensor.prototype.addToScene = function(surface){
  if (!this.paintRoot) {
    this.paintRoot = surface.project.importSVG(sensorSvg, {
      expandShapes: true,
      insert: false,
    });

    this.wifiSections = new surface.Group({children: this.paintRoot.children[1].removeChildren(2, 4)});
    if (this.showWifi) {
      this.paintRoot.addChild(this.wifiSections);
    }
    this.buttonSections = new surface.Group({children: this.paintRoot.children[1].removeChildren(1, 2)});
    if (this.showButton) {
      this.paintRoot.addChild(this.buttonSections);
    }

    this.paintRoot.scale(this.size);
    this.paintRoot.position = this.pos;
    this.paintRoot.name = this.name;
    this.paintRoot.strokeColor = xiloStrokeColor('Sensor', this.name);
    this.paintRoot.fillColor = xiloFillColor('Sensor', this.name);
    surface.project.activeLayer.addChild(this.paintRoot);
    this.wifiSections.fillColor = xiloFillColor('Sensor.Emission', this.name);
    this.wifiSections.strokeColor = xiloStrokeColor('Sensor.Emission', this.name);
    this.buttonSections.fillColor = xiloFillColor('Sensor.Button', this.name);
    this.buttonSections.strokeColor = xiloStrokeColor('Sensor.Button', this.name);
  }
}
