package main

import (
  "github.com/howeyc/fsnotify"
  ui "github.com/gizak/termui"
  "io/ioutil"
  "strconv"
  "strings"
  "path"
  "time"
  "sort"
  "fmt"
  "os"
)

var pluginDir string
var filenamesToResourceIDs map[string]int
var watcher *fsnotify.Watcher

type ChangeEvent struct {
  FileName string
  ResourceID int
  Data []byte
}

func setupFiles(dir string){
  if dir == ""{
    aPluginID := strconv.Itoa(int(pluginInfo.ID))
    cwd, _ := os.Getwd()
    pluginDir = path.Join(cwd, aPluginID)
  } else {
    pluginDir = dir
  }
  fmt.Println("Setting up plugin files in:", pluginDir)

  e := os.Mkdir(pluginDir, 0777)
  var err *os.PathError
  err, _ = e.(*os.PathError)
  if err != nil{
    if err.Err.Error() == "file exists" {
      fmt.Println("Directory already exists, using it")
    } else{
      fmt.Println("Error creating directory:", err)
      fmt.Println("Proceeding anyway")
    }
  }

  writeOutResources()
  setupWatcher()
}

func setupWatcher(){

  var err error
  watcher, err = fsnotify.NewWatcher()
  if err != nil{
    fmt.Println("Error creating watcher:", err)
  }

  err = watcher.Watch(pluginDir)
  if err != nil {
      fmt.Println("Error watching resource directory:", err)
  }

  go func() {
    lastEventFile := ""
    wasModify := false
    lastEventTime := time.Now()
    for {
      select {
      case ev := <-watcher.Event:
          if (ev.Name == lastEventFile) && (ev.IsModify() == wasModify) && (time.Now().Sub(lastEventTime) < (time.Millisecond*600)) {
            //must be duplicate
          } else if ev.IsModify() {
            lastEventFile = ev.Name
            wasModify = ev.IsModify()
            lastEventTime = time.Now()

            resID, ok := filenamesToResourceIDs[path.Base(ev.Name)]
            if ok{
              data, _ := ioutil.ReadFile(path.Join(pluginDir, path.Base(ev.Name)))
              e := ChangeEvent{
                FileName: path.Base(ev.Name),
                ResourceID: resID,
                Data: data,
              }
              ui.SendCustomEvt("/resource/update", e)
              //fmt.Println("Generated update event:", e.FileName)

            } else {
              fmt.Println("Recieved event for file not currently under tracking:", ev.Name)
            }
          }
      case err := <-watcher.Error:
          fmt.Println("watcher error:", err)
      }
    }
  }()
}

func writeOutResources() {
  filenamesToResourceIDs = map[string]int{}
  for _, resource := range pluginInfo.Resources {
    name := sanitizeName(resource.Name)
    if resource.IsExecutable {
      name = name + ".js"
    }

    filenamesToResourceIDs[name] = resource.ID
    fmt.Println("Writing:", path.Join(pluginDir, name))
    w, err := os.Create(path.Join(pluginDir, name))
    if err != nil{
      fmt.Println("Error!:", err, "Proceeding anyway.")
    }
    w.Write(resource.Data)
    w.Close()
  }
}

func sanitizeName(in string)string{
  return strings.Replace(strings.Replace(in, "/", "", -1), "*", "", -1)
}

func getResourceList(highlight string)[]string{
  var output []string
  for fname, _ := range filenamesToResourceIDs {
    out := fname
    if fname == highlight {
      out += " (*)"
    }
    output = append(output, out)
  }
  sort.Strings(output)
  return output
}

func getResourceFromID(id int)Resource{
  for _, res := range pluginInfo.Resources {
    if res.ID == id {
      return res
    }
  }
  return  Resource{}
}
