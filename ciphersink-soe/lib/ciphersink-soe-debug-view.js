'use babel';

import * as util from './ciphersink-soe-util';

export default class CiphersinkSoeDebugView {

  constructor() {
    // Create root element
    this.element = document.createElement('div');
    this.element.classList.add('ciphersink-soe');

    // Create header
    const header = document.createElement('h1');
    header.textContent = 'Ciphersink-SOE Debug Information';
    header.classList.add('ciphersink-soe-debug-header');
    this.element.appendChild(header);

    // Create status element
    this.statusContainer = document.createElement('div');
    this.statusContainer.textContent = 'Status: Good';
    this.statusContainer.classList.add('icon','icon-question');
    this.element.appendChild(this.statusContainer);

    // Create paths element
    this.pathsContainer = document.createElement('div');
    this.pathsContainer.textContent = 'Path: <>';
    this.element.appendChild(this.pathsContainer);
  }

  // Returns an object that can be retrieved when package is activated
  serialize() {}

  // Tear down any state and detach
  destroy() {
    this.element.remove();
    //dont need to remove elements children
  }

  getElement() {
    return this.element;
  }

  update(controller){
    this.pathsContainer.textContent = 'Path: ' + util.getProjectPath();

    this.packageConfig = util.readJsonInProjectRoot("package.json");
    this.statusContainer.textContent = "" + this.packageConfig.name + ' - (v' + this.packageConfig.version + ") (Atom dev mode = " + atom.inDevMode() + ")";
  }

}
