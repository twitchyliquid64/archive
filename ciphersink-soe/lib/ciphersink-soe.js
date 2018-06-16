'use babel';

import CiphersinkSoeDebugView from './ciphersink-soe-debug-view';
import ProgressModal from './viewControllers/progressModal'
import SelectOptionModal from './viewControllers/selectOptionModal'

import { CompositeDisposable } from 'atom';
import * as util from './ciphersink-soe-util';

export default {
  subscriptions: null,

  setupDebugView(){
    this.debugView = new CiphersinkSoeDebugView();
    this.debugViewPanel = atom.workspace.addModalPanel({
      item: this.debugView.getElement(),
      visible: false
    });
  },

  activate(state) {
    this.setupDebugView(state);

    // Events subscribed to in atom's system can be easily cleaned up with a CompositeDisposable
    this.subscriptions = new CompositeDisposable();

    // Register command that toggles this view
    this.subscriptions.add(atom.commands.add('atom-workspace', {
      'ciphersink-soe:toggle-dev-information': () => this.toggleDevInformation()
    }));
    this.subscriptions.add(atom.commands.add('atom-workspace', {
      'ciphersink-soe:deploy': () => this.deploy()
    }));
    this.subscriptions.add(atom.commands.add('atom-workspace', {
      'ciphersink-soe:test-select-modal': () => this.testSelectModal()
    }));
  },

  deactivate() {
    this.debugViewPanel.destroy();
    this.subscriptions.dispose();
    this.debugView.destroy();
    if (this.deployProgressModal){
      this.deployProgressModal.destroy();
    }
  },

  serialize() {
    return {};
  },

  toggleDevInformation() {
    console.log('Debug information was toggled!');
    this.debugView.update(this);
    this.debugViewPanel.isVisible() ? this.debugViewPanel.hide() : this.debugViewPanel.show();
  },

  deploy() {
    console.log("deploy()");
    if (!this.deployProgressModal){
      this.deployProgressModal = new ProgressModal();
      this.deployProgressModal.show();
    } else {
      this.deployProgressModal.hide();
      this.deployProgressModal.destroy();
      this.deployProgressModal = null;
    }
  },


  testSelectModal(){
    console.log("testOptions()");
    new SelectOptionModal(["hi", "hello there"], function(item){
      console.log(item);
    });
  }

};
