package main

import (
  "github.com/0xAX/notificator"
  "github.com/kardianos/osext"
  "path"
)

var notify *notificator.Notificator
func initNotifications(){
  notify = notificator.New(notificator.Options{
    DefaultIcon: path.Join(iconsFolder(), "ic_build_black_24dp_2x.png"),
    AppName:     "CNC Dev Client",
  })
}


func iconsFolder()string{
  exeFolder, err := osext.ExecutableFolder()
  if err != nil{
    panic(err)
  }
  return path.Join(exeFolder, "resources", "icons")
}

func sendBuildNotification(msg string){
  notify.Push("CNC Dev Client", msg, path.Join(iconsFolder(), "ic_build_black_24dp_2x.png"), notificator.UR_NORMAL)
}

func sendTransmissionNotification(msg string){
  notify.Push("CNC Dev Client", msg, path.Join(iconsFolder(), "ic_cloud_done_black_24dp_2x.png"), notificator.UR_NORMAL)
}
