package main

import (
  ui "github.com/gizak/termui"
  "strings"
  "fmt"
)

var topText *ui.Par
var statusText *ui.Par
var resourceList *ui.List
var memoryUsage *ui.Gauge
var hooksUsage *ui.Sparkline
var msgText *ui.Par

func initUI() {
  if err := ui.Init(); err != nil {
    panic(err)
  }

  //top bar which says the plugin
  topText = ui.NewPar(pluginDir)
  topText.TextFgColor = ui.ColorWhite
  topText.BorderLabel = pluginInfo.Name
  topText.BorderFg = ui.ColorCyan
  topText.Height = 3

  //top bar which says the plugin status
  statusText = ui.NewPar("Ready")
  statusText.TextFgColor = ui.ColorWhite
  statusText.BorderLabel = "Status"
  statusText.BorderFg = ui.ColorCyan
  statusText.Height = 3

  //display which shows all the resources and their status
  resourceList = ui.NewList()
  resourceList.ItemFgColor = ui.ColorYellow
	resourceList.BorderLabel = "Resources"
	resourceList.Height = 14
  resourceList.Items = getResourceList("")

  //server memory usage guage
  memoryUsage = ui.NewGauge()
  memoryUsage.Percent = 0
  memoryUsage.Height = 5
  memoryUsage.BorderLabel = "Server Memory Usage"
	memoryUsage.BarColor = ui.ColorYellow
	memoryUsage.BorderFg = ui.ColorWhite

  data := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
  hooksUsage := ui.NewSparkline()
  hooksUsage.Data = data[3:]
  hooksUsage.LineColor = ui.ColorMagenta
  hooksUsage.Height = 6
  sp := ui.NewSparklines(hooksUsage)
	sp.Height = 9
	sp.BorderLabel = "Serverside Hooks"

  msgText = ui.NewPar("Press R to reload the plugin on the server.\nSystem ready.\n")
  msgText.TextFgColor = ui.ColorWhite
  msgText.BorderLabel = "Messages"
  msgText.BorderFg = ui.ColorCyan
  msgText.Height = 6

  ui.Body.AddRows(
      ui.NewRow(
        ui.NewCol(7, 0, topText),
        ui.NewCol(5, 0, statusText)),
      ui.NewRow(
        ui.NewCol(7, 0, resourceList),
        ui.NewCol(5, 0, memoryUsage, sp)),
      ui.NewRow(
        ui.NewCol(12, 0, msgText)))


  // calculate layout
  ui.Body.Align()
  ui.Render(ui.Body)
}


func writeUiMessage(msg string){
  msgs := strings.Split(msgText.Text, "\n")
  msgs = append(msgs, msg)

  for len(msgs) > 3 {
    msgs = msgs[1:]
  }
  msgText.Text = strings.Join(msgs, "\n")
  //ui.Render(ui.Body)
}

func runUILoop(){
  ui.Handle("/sys/kbd/r", func(ui.Event) {
        writeUiMessage("Rebuilding plugin")
        statusText.Text = "Rebuilding Plugin"
        ui.Render(ui.Body)
        CONN.Write(newPacket(&DataRequest{DataType: REQUEST_RESTART}).Serialize())
        if successful := <-updateSuccessful; !successful {
          ui.StopLoop()
          return
        }
        statusText.Text = "Ready"
        ui.Render(ui.Body)
  })


  ui.Handle("/sys/kbd/q", func(ui.Event) {
        ui.StopLoop()
  })

  ui.Handle("/resource/update", func(e ui.Event) {
        ce := e.Data.(ChangeEvent)
        statusText.Text = "Pushing Changes"
        resourceList.Items = getResourceList(ce.FileName)
        ui.Render(ui.Body)

        r := getResourceFromID(ce.ResourceID)
        writeUiMessage("Updating: " + r.Name)
        r.Data = ce.Data
        CONN.Write(newPacket(&ResourceUpdate{R: r}).Serialize())

        if successful := <-updateSuccessful; !successful {
          ui.StopLoop()
          return
        }

        statusText.Text = "Ready, needs rebuild"
        resourceList.Items = getResourceList("")
        ui.Render(ui.Body)
  })

  ui.Handle("/remote/log", func(e ui.Event){
    lm := e.Data.(_LogMessage)
    writeUiMessage("(REMOTE)(" + lm.Type + ")(" + lm.Component + "): " + lm.Message)
  })

  ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Render(ui.Body)
	})

  defer func(){
    ui.Close()
    if (remoteError != nil && *remoteError != "") {
      fmt.Println(*remoteError)
    }
  }()
  ui.Loop()
}
