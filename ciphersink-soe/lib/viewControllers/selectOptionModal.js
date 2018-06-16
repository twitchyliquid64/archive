'use babel';

import * as util from '../ciphersink-soe-util';
import { CompositeDisposable } from 'atom';

import {$$, View} from 'space-pen';
import {SelectListView} from 'atom-space-pen-views';

import jquery from 'jquery';
$ = jquery;

export default class SelectOptionModal extends SelectListView {

  constructor (items, selectionCB){
    super();
    this.addClass('overlay from-top');

    this.selectionCB = selectionCB;
    this.setItems(items);

    if (this.panel == null) {
      this.panel = atom.workspace.addModalPanel({
        item: this
      });
    }
    this.panel.show();

    this.subscriptions = new CompositeDisposable();

    this.focusFilterEditor();
  }

  viewForItem(item){
    return "<li><div class=\"icon file icon-file-text\">" + item + "</div></li>";
  }

  confirmed(item) {
    this.destroy();
    if (this.selectionCB)this.selectionCB(item);
  }

  cancelled() {
    this.destroy();
  }

  destroy() {
    this.panel.destroy();
    this.subscriptions.dispose();
  }


}
