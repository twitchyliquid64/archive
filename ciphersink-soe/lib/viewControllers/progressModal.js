'use babel';

import * as util from '../ciphersink-soe-util';

export default class ProgressModal {

  constructor(serializedState) {
    // Create root element
    this.element = document.createElement('div');
    this.element.classList.add('progress-modal');

    // Create header element
    this.header = document.createElement('h2');
    this.header.textContent = "Deployment";
    this.header.classList.add('align-left');
    this.element.appendChild(this.header);

    // Create progress container element
    this.progressContainer = document.createElement('div');
    this.progressContainer.classList.add('block')

    // Create spinner and text elements
    this.progressSpinner = document.createElement('div');
    this.progressSpinner.classList.add('inline-block','spinner', 'progress-spinner')
    this.progressText = document.createElement('span');
    this.progressText.classList.add('inline-block')
    this.progressText.textContent = "Loading ... Please Wait.";

    this.progressContainer.appendChild(this.progressSpinner);
    this.progressContainer.appendChild(this.progressText);
    this.element.appendChild(this.progressContainer);

    this._createModal();
  }

  _createModal() {
    this.modal = atom.workspace.addModalPanel({
      item: this.getElement(),
      visible: false
    });
  }

  show() {
    this.modal.show();
  }

  hide() {
    this.modal.hide();
  }

  isVisible() {
    return this.modal.isVisible();
  }

  // Returns an object that can be retrieved when package is activated
  serialize() {}

  // Tear down any state and detach
  destroy() {
    this.element.remove();
    this.modal.destroy();
    //dont need to remove elements children
  }

  getElement() {
    return this.element;
  }
}
