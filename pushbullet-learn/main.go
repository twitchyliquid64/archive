    package main

    import (
        "fmt"

        "github.com/howeyc/pbullet"
    )

    func main() {
        pbullet.SetAPIKey("PUT YAH KEY IN HEAR") // https://www.pushbullet.com/settings

        devList, err := pbullet.GetDevices()
        if err != nil {
            //fmt.Println(err)
        }
        for _, dev := range devList.Devices {
	       fmt.Println(dev)
            _, pushErr := dev.PushNote("IOFlame", "LOLCAKES")
            if pushErr != nil {
                fmt.Println(err)
            }
        }
        fmt.Println("Done")
    }
