package main

import (
  "flag"
  "time"
  "fmt"
)

func main() {

  srvAddressPtr := flag.String("addr", "", "The address of the server formatted as <hostname> or <hostname>:<port>")
  userPtr := flag.String("user", "", "The username to use when logging in to CNC. The user must be an admin.")
  passPtr := flag.String("pass", "", "The password of the user specified in 'pass'.")
  pluginNamePtr := flag.String("plugin", "", "The plugin which to operate on.")
  dirPtr := flag.String("dir", "", "The directory in which the resources of a plugin will be saved in.")
  flag.Parse()

  if (*srvAddressPtr == ""){
    fmt.Println("Must specify a server: --addr ...")
    return
  }
  if (*userPtr == ""){
    fmt.Println("Must specify a username: --user ...")
    return
  }
  if (*passPtr == ""){
    fmt.Println("Must specify a password: --pass ...")
    return
  }

  initNotifications()

  fmt.Print("Connecting...")
  err := connect(*srvAddressPtr, *userPtr, *passPtr, *pluginNamePtr)
  if err == nil{
    fmt.Println("OK.")
  } else {
    fmt.Println("Error.")
    fmt.Println(err)
    return
  }
  defer CONN.Close()

  fmt.Print("Authenticating...")
  if didAuthenticate := <- authenticated; didAuthenticate {
    fmt.Println("OK.")
  } else {
    fmt.Println("Error.")
    fmt.Println(*remoteError)
    return
  }

  //wait until ready
  if isReady := <-ready; !isReady {
    fmt.Println("Remote [ERROR]: " + *remoteError)
    return
  }


  if *pluginNamePtr == "" {//we are getting a list of plugins and displaying them
    if gotList := <-gotPluginData; !gotList {
      fmt.Println("Remote [ERROR]: " + *remoteError)
      return
    }
    fmt.Println("Enabled Plugins:")
    for _, p := range plugins {
      fmt.Println(p.ID, "\t", p.Name)
    }
  } else { //normal operation
    fmt.Print("Downloading Plugin data...")
    CONN.Write(newPacket(&DataRequest{DataType: REQUEST_PLUGININFO}).Serialize())
    if didDownloadPlugin := <- hasGotPluginData; didDownloadPlugin {
      fmt.Println("OK.")
    } else {
      fmt.Println("Error.")
      fmt.Println(*remoteError)
      return
    }

    setupFiles(*dirPtr)
    defer watcher.Close()

    sendTransmissionNotification("Watcher service running.")
    initUI()
    runUILoop()
  }
  time.Sleep(time.Millisecond * 400)
}
