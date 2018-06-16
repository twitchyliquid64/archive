function Router(name, opts) {
  this.name = name || "router";
  this.pos = opts.pos || {x: 0, y: 0};
  this.downNodes = opts.downNodes || [];
  this.downLinks = opts.downLinks || [];
  this.size = opts.size || 3.5;

  this.paintRoot = null;
}

// state-aware API methods
Router.prototype.stateDescription = function(){
  return [];
}

// connection-aware API methods
Router.prototype.setUpstreamNode = function(up, upLink){
  this.upNode = up;
  this.upLink = upLink;
}
Router.prototype.setDownstreamNode = function(down, downLink){
  this.downNodes[this.downNodes.length] = down;
  this.downLinks[this.downLinks.length] = downLink || null;
}
Router.prototype.connectionDescription = function(){
  return [];
}

Router.prototype.bounds = function(){
  var bound = this.paintRoot.children[1].bounds.clone();
  bound.top += 10;
  return bound;
}

var routerSvg = '\
<svg height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg">\
    <path d="M20.2 5.9l.8-.8C19.6 3.7 17.8 3 16 3s-3.6.7-5 2.1l.8.8C13 4.8 14.5 4.2 16 4.2s3 .6 4.2 1.7zm-.9.8c-.9-.9-2.1-1.4-3.3-1.4s-2.4.5-3.3 1.4l.8.8c.7-.7 1.6-1 2.5-1 .9 0 1.8.3 2.5 1l.8-.8zM19 13h-2V9h-2v4H5c-1.1 0-2 .9-2 2v4c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-4c0-1.1-.9-2-2-2zM8 18H6v-2h2v2zm3.5 0h-2v-2h2v2zm3.5 0h-2v-2h2v2z"/>\
</svg>'

Router.prototype.addToScene = function(surface){
  if (!this.paintRoot) {
    this.paintRoot = surface.project.importSVG(routerSvg, {
      expandShapes: true,
      insert: false,
    });
    this.wifiSections = new surface.Group({children: this.paintRoot.children[1].removeChildren(0, 2)});
    this.paintRoot.addChild(this.wifiSections);

    this.paintRoot.scale(this.size);
    this.paintRoot.position = this.pos;
    this.paintRoot.name = this.name;
    this.paintRoot.strokeColor = xiloStrokeColor('Router', this.name);
    this.paintRoot.fillColor = xiloFillColor('Router', this.name);
    surface.project.activeLayer.addChild(this.paintRoot);
    this.wifiSections.fillColor = xiloFillColor('Router.Emission', this.name);
    this.wifiSections.strokeColor = xiloStrokeColor('Router.Emission', this.name);
  }
}
