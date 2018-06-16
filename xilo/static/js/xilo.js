function Xilo(canvas) {
  this.canvas = canvas;
  this.paperSurface = new paper.PaperScope();
  this.paperSurface.setup(canvas);
  this.paperSurface.settings.insertItems = false;
  this.paperSurface.project.view.onFrame = this.onFrame;
  this.paperSurface.project.view.onMouseDown = this.onMouseDown.bind(this);

  this.itemsByName = {};
  this.items = [];
  this.links = [];
}

// adds a non-pipe component to the graph.
Xilo.prototype.add = function(item) {
  this.items.push(item);
  this.itemsByName[item.name] = item;
  item.addToScene(this.paperSurface);
}

Xilo.prototype.addLink = function(link) {
  link.draw(this.paperSurface);
}

Xilo.prototype.onFrame = function(event){
}

Xilo.prototype.onMouseDown = function(event){
  this.paperSurface.project.activeLayer.selected = false;
    var hitResult = this.paperSurface.project.hitTest(event.point, {
  	segments: true,
  	stroke: true,
  	fill: true,
  	tolerance: 5
  });
	if (!hitResult || !hitResult.item.parent.name)
		return;

  var n = hitResult.item.parent.name;
  console.log(n, ": ", this.itemsByName[n]);
  var attrs = this.itemsByName[n].stateDescription()
  for (var i = 0; i < attrs.length; i++) {
      console.log("\t", attrs[i].name, attrs[i].get());
  }
  hitResult.item.parent.selected = true;
}


function lerpColor(a, b, amount) {
    var ah = parseInt(a.replace(/#/g, ''), 16),
        ar = ah >> 16, ag = ah >> 8 & 0xff, ab = ah & 0xff,
        bh = parseInt(b.replace(/#/g, ''), 16),
        br = bh >> 16, bg = bh >> 8 & 0xff, bb = bh & 0xff,
        rr = ar + amount * (br - ar),
        rg = ag + amount * (bg - ag),
        rb = ab + amount * (bb - ab);
    return '#' + ((1 << 24) + (rr << 16) + (rg << 8) + rb | 0).toString(16).slice(1);
};

var linkIndex = 0;

function Link(upPos, downPos, upNode, downNode, disposition) {
  console.log("Making link", linkIndex, "From", upPos, "to", downPos);
  this.linkID = linkIndex++;
  this.upPos = upPos;
  this.downPos = downPos;
  this.upNode = upNode;
  this.downNode = downNode;
  this.disposition = disposition || 'VERT-HORIZONTAL';
  if (upNode && downNode) {
    downNode.setUpstreamNode(upNode, this);
    upNode.setDownstreamNode(downNode, this);
    this.name = upNode.name + '-' + downNode.name;
  }
}

Link.prototype.draw = function(paperSurface) {
  if ((this.upPos.y == this.downPos.y) || (this.upPos.x == this.downPos.x)) {
    var l = new paperSurface.Path({segments: [this.downPos, this.upPos], strokeWidth: 1, strokeColor: xiloStrokeColor('Link', this.name)});
    paperSurface.project.activeLayer.addChild(l);
    l.sendToBack();
  } else {
    switch (this.disposition) {
      case 'DIRECT':
        var l = new paperSurface.Path({segments: [this.downPos, this.upPos], strokeWidth: 1, strokeColor: xiloStrokeColor('Link', this.name)});
        paperSurface.project.activeLayer.addChild(l);
        l.sendToBack();
        break;
      case 'VERT-HORIZONTAL':
        var vertical = new paperSurface.Path({segments: [this.upPos, {x: this.upPos.x, y: this.downPos.y-5}], strokeWidth: 1, strokeColor: xiloStrokeColor('Link', this.name)});
        paperSurface.project.activeLayer.addChild(vertical);
        vertical.sendToBack();
        var horizontal = new paperSurface.Path({segments: [this.downPos, {x: this.upPos.x+5, y: this.downPos.y}], strokeWidth: 1, strokeColor: xiloStrokeColor('Link', this.name)});
        paperSurface.project.activeLayer.addChild(horizontal);
        horizontal.sendToBack();
        var curve = new paperSurface.Path.Arc({
            from:     {x: this.upPos.x, y: this.downPos.y-5},
            through:  {x: this.upPos.x+1, y: this.downPos.y-1} ,
            to:       {x: this.upPos.x+5, y: this.downPos.y},
            strokeColor: xiloStrokeColor('Link', this.name),
        });
        paperSurface.project.activeLayer.addChild(curve);
        curve.sendToBack();
        break;
      case 'HORI-VERTICAL':
        var horizontal = new paperSurface.Path({segments: [this.upPos, {x: this.downPos.x-5, y: this.upPos.y}], strokeWidth: 1, strokeColor: xiloStrokeColor('Link', this.name)});
        paperSurface.project.activeLayer.addChild(horizontal);
        horizontal.sendToBack();
        var vertical = new paperSurface.Path({segments: [this.downPos, {x: this.downPos.x, y: this.upPos.y+5}], strokeWidth: 1, strokeColor: xiloStrokeColor('Link', this.name)});
        paperSurface.project.activeLayer.addChild(vertical);
        vertical.sendToBack();
        var curve = new paperSurface.Path.Arc({
            from:     {x: this.downPos.x-5, y: this.upPos.y},
            through:  {x: this.downPos.x-1, y: this.upPos.y+1} ,
            to:       {x: this.downPos.x, y: this.upPos.y+5},
            strokeColor: xiloStrokeColor('Link', this.name)
        });
        paperSurface.project.activeLayer.addChild(curve);
        curve.sendToBack();
        break;
    }
  }
  return this;
}

function BoundLabel(parent, opts) {
  this.name = opts.name || 'unnamed label';
  this.content = opts.content || this.name;
  this.parent = parent;
  this.offset = opts.offset || {x: 0, y: 0};
  this.size = opts.size || 60;
  this.justification = opts.justification || 'center';
}

BoundLabel.prototype.addToScene = function(surface) {
  this.text = new surface.PointText(new surface.Point(this.size/2, this.size/2));
  var c = this.parent.bounds().center;
  this.text.position = {x: c.x + this.offset.x, y: c.y + this.offset.y};
  this.text.justification = this.justification;
  this.text.strokeColor = xiloStrokeColor('BoundLabel', this.name);
  this.text.fillColor = xiloFillColor('BoundLabel', this.name);
  this.text.strokeWidth = 0.1;
  this.text.content = this.content;
  this.text.scale(this.size/50);
  surface.project.activeLayer.addChild(this.text);
}

var xiloStyles = {
  strokeColor: 'black',
  fillColor: 'white',
  strokeWidth: 1,
  byName: {},
  byKind: {
    'Router.Emission': {fillColor: '#990000'},
    'BoundLabel': {fillColor: 'white', strokeColor: 'white'},
  },
}
function xiloSetStrokeColor(color) {
  xiloStyles.strokeColor = color;
}
function xiloSetFillColor(color) {
  xiloStyles.fillColor = color;
}
function xiloSetKindDetail(detail, kind, value) {
  if (!xiloStyles.byKind[kind])
    xiloStyles.byKind[kind] = {};
  xiloStyles.byKind[kind][detail] = value;
}
function xiloSetEntityDetail(detail, name, value) {
  if (!xiloStyles.byName[name])
    xiloStyles.byName[name] = {};
  xiloStyles.byName[name][detail] = value;
}

function xiloSetKindStrokeColor(kind, color) {
  return xiloSetKindDetail('strokeColor', kind, color);
}
function xiloSetEntityStrokeColor(name, color) {
  return xiloSetEntityDetail('strokeColor', name, color)
}

function xiloSetKindFillColor(kind, color) {
  return xiloSetKindDetail('fillColor', kind, color);
}
function xiloSetEntityFillColor(kind, name) {
  return xiloSetEntityDetail('fillColor', name, color)
}

function xiloStrokeColor(kind, name) {
  if (name && name in xiloStyles.byName && 'strokeColor' in xiloStyles.byName[name])
    return xiloStyles.byName[name].strokeColor;
  if (kind && kind in xiloStyles.byKind && 'strokeColor' in xiloStyles.byKind[kind])
    return xiloStyles.byKind[kind].strokeColor;
  return xiloStyles.strokeColor;
}
function xiloFillColor(kind, name) {
  if (name && name in xiloStyles.byName && 'fillColor' in xiloStyles.byName[name])
    return xiloStyles.byName[name].fillColor;
  if (kind && kind in xiloStyles.byKind && 'fillColor' in xiloStyles.byKind[kind])
    return xiloStyles.byKind[kind].fillColor;
  return xiloStyles.fillColor;
}
