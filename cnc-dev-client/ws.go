package main

import (
  ui "github.com/gizak/termui"
  "golang.org/x/net/websocket"
  "encoding/json"
  "crypto/tls"
  "net/url"
  "fmt"
)

var CONN *websocket.Conn

var authenticated chan bool
var ready chan bool
var chansInitialized bool
var gotPluginData chan bool
var hasGotPluginData chan bool
var updateSuccessful chan bool

var remoteError *string
var plugins []Plugin
var pluginInfo Plugin

func initChans(){
  authenticated = make(chan bool, 1)
  ready = make(chan bool, 1)
  gotPluginData = make(chan bool, 1)
  hasGotPluginData = make(chan bool, 1)
  updateSuccessful = make(chan bool, 1)
  chansInitialized = true
}

func killChans(){
  if chansInitialized{
    chansInitialized = false
    close(authenticated)
    close(ready)
    close(gotPluginData)
    close(hasGotPluginData)
    close(updateSuccessful)
  }
}

func connect(serv, user, pass, plugin string)error{
  initChans()

  origin := "http://" + serv + "/"
  urlS := "wss://" + serv + "/ws/devclient"

  config, err := websocket.NewConfig(urlS, origin)
  if err != nil{
    return err
  }

  config.Location.RawQuery = "user=" + user + "&pass=" + pass + "&plugin=" + url.QueryEscape(plugin)

  config.TlsConfig = &tls.Config{
    InsecureSkipVerify: true,
  }

  CONN, err = websocket.DialConfig(config)
  if err != nil{
    return err
  }
  go readRoutine()
  return nil
}


func processTextMessage(data []byte){
  decodeTextMsg(data)
  panic("not implemented")
}

func processError(data []byte){
  msg := decodeFatalError(data)
  remoteError = &msg.Error
  killChans()
}

func processPluginList(data []byte){
  msg := decodePluginList(data)
  plugins = msg.Plugins
  gotPluginData <- true
}

func processPluginInfo(data []byte){
  msg := decodePluginInfo(data)
  pluginInfo = msg.P
  hasGotPluginData <- true
}

func processLogMessage(data []byte){
  msg := decodeLogMessage(data)
  ui.SendCustomEvt("/remote/log", msg.Msg)
}

func processStatusMessage(data []byte) {
  defer func(){
    //recover()
  }()

  msg := decodeStatusMsg(data)
  switch msg.Status {
  case STATUS_AUTHENTICATED:
    select {
      case authenticated <- true:
      default:
    }

  case STATUS_SAVE_SUCCESSFUL:
    select {
    case updateSuccessful <- true:
    default:
    }

  case STATUS_READY:
    select {
      case ready <- true:
      default:
    }


  }
}

func processMessage(data []byte){
  var pkt Packet
  err := json.Unmarshal(data, &pkt)
  if err != nil{
    fmt.Println("JSON Error: ", err)
    return
  }

  switch pkt.Type {
  case "txtmsg":
    processTextMessage(pkt.Subdata)
  case "status":
    processStatusMessage(pkt.Subdata)
  case "ferror":
    processError(pkt.Subdata)
  case "plist":
    processPluginList(pkt.Subdata)
  case "plugininfo":
    processPluginInfo(pkt.Subdata)
  case "logMessage":
    processLogMessage(pkt.Subdata)
  default:
    fmt.Println("Unknown type: ", pkt.Type, " ---- ", string(pkt.Subdata))
  }
}

func readRoutine() {
  for {
    var data []byte
    err := websocket.Message.Receive(CONN, &data)
    if err != nil{
      killChans()
      CONN.Close()
      return
    }
    processMessage(data)
  }
}
