'use babel';

import { File } from 'atom';

import * as path from 'path'
import { readFileSync } from 'fs';

export default {
  log(msg){
    console.log(msg);
  },

  // Returns the full path of the file currently in focus.
  getCurrentFilePath() {
    return atom.workspace.getActiveTextEditor().buffer.file.path;
  },

  // Returns the path to the base of the current project.
  getProjectPath() {
    let paths = atom.project.getPaths();
    let currentFilePath = this.getCurrentFilePath();
    for(var i = 0; i < paths.length; i++){ //returns first path that prefixes the path of the current file.
      if (currentFilePath.lastIndexOf(paths[i], 0) === 0){
        return paths[i];
      }
    }

    if (paths.length > 1){
      console.log("Warning: More than one project path, none of which could be matched to the current file. Implicitly picking the first project path.");
      console.log(paths);
    }
    return paths[0];
  },

  readJsonInProjectRoot(fName){
    try{
      return JSON.parse(readFileSync(path.join(this.getProjectPath(), fName)).toString());
    } catch (e) {
      if (e.code === 'ENOENT') {
        console.log(path.join(this.getProjectPath(), fName));
        console.log('util.readJsonInProjectRoot(): File not found!');
        return {"error": "File not found", "exception": e};
      } else {
        console.log(path.join(this.getProjectPath(), fName));
        console.log('util.readJsonInProjectRoot(): ' + e);
        throw e;
      }
    }
  },

  isAtomPackageInstalled(pname){
    if (atom.packages.loadedPackages[pname])
      return true;
    return false;
  }
}
