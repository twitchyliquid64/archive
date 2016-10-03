package main

import (
  "encoding/json"
  "time"
  "fmt"
)

type Packet struct {
  Type string
  Subdata []byte
}

func (p *Packet)Serialize()[]byte{
  d, err := json.Marshal(p)
  if err != nil{
    fmt.Println(err)
    return []byte("")
  }
  return d
}

type Subpacket interface{
  Typ()string
}

func newPacket(pkt Subpacket)*Packet{
  d, err := json.Marshal(pkt)
  if err != nil{
    fmt.Println(err)
    return nil
  }

  return &Packet{
    Type: pkt.Typ(),
    Subdata: d,
  }
}


type TextMsg struct{
  Fatal bool
  Message string
}
func (m *TextMsg)Typ()string{
  return "txtmsg"
}
func decodeTextMsg(data []byte)*TextMsg{
  var t TextMsg
  err := json.Unmarshal(data, &t)
  if err != nil{
    return nil
  }
  return &t
}


type FatalError struct{
  Error string
}
func (m *FatalError)Typ()string{
  return "ferror"
}
func decodeFatalError(data []byte)*FatalError{
  var t FatalError
  err := json.Unmarshal(data, &t)
  if err != nil{
    return nil
  }
  return &t
}


const (
  STATUS_AUTHENTICATED string = "AUTH OK"
  STATUS_READY string = "READY"
  STATUS_SAVE_SUCCESSFUL string = "SAVE GOOD"

)


type Status struct{
  Status string
}
func (m *Status)Typ()string{
  return "status"
}

func decodeStatusMsg(data []byte)*Status{
  var t Status
  err := json.Unmarshal(data, &t)
  if err != nil{
    return nil
  }
  return &t
}







type Plugin struct {
    ID        uint `gorm:"primary_key"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time

    Name string `sql:"not null;unique;index"`
    Icon string
    Description string
    Enabled bool

    HasCrashed bool
    ErrorStr string

    Resources []Resource
}


type Resource struct {
  ID int      `gorm:"primary_key"`
  PluginID int `sql:"index"`
  Name string `sql:"index"`
  Data []byte
  IsExecutable bool
  IsTemplate bool
  JSONData string `sql:"-"` //only used for JSON deserialisation - not a DB field
}



type PluginList struct{
  Plugins []Plugin
}
func (m *PluginList)Typ()string{
  return "plist"
}

func decodePluginList(data []byte)*PluginList{
  var t PluginList
  err := json.Unmarshal(data, &t)
  if err != nil{
    return nil
  }
  return &t
}




const (
  REQUEST_PLUGININFO string = "plugininfo"
  REQUEST_RESTART string = "pluginRestart"
)
type DataRequest struct{
  DataType string
  ID int
}
func (m *DataRequest)Typ()string{
  return "dataRequest"
}
func decodeDataRequest(data []byte)*DataRequest{
  var t DataRequest
  err := json.Unmarshal(data, &t)
  if err != nil{
    return nil
  }
  return &t
}


type PluginInfo struct{
  P Plugin
}
func (m *PluginInfo)Typ()string{
  return "plugininfo"
}
func decodePluginInfo(data []byte)*PluginInfo{
  var t PluginInfo
  err := json.Unmarshal(data, &t)
  if err != nil{
    return nil
  }
  return &t
}




type ResourceUpdate struct{
  R Resource
}
func (m *ResourceUpdate)Typ()string{
  return "resourceUpdate"
}
func decodeResourceUpdate(data []byte)*ResourceUpdate{
  var t ResourceUpdate
  err := json.Unmarshal(data, &t)
  if err != nil{
    return nil
  }
  return &t
}



type _LogMessage struct {
  Component string
  Type string
  Message string
  Created int64
}


type LogMessage struct{
  Msg _LogMessage
}
func (m *LogMessage)Typ()string{
  return "logMessage"
}
func decodeLogMessage(data []byte)*LogMessage{
  var t LogMessage
  err := json.Unmarshal(data, &t)
  if err != nil{
    return nil
  }
  return &t
}
