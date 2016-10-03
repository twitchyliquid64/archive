package main

import "github.com/twitchyliquid64/isnap"
import "image"
import "fmt"
import "runtime"

func main(){
	
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	img, err := isnap.Jpg("test.jpg")
	fmt.Println(err)
	
	colourFiltered := isnap.ColourFilter(img, 240, 0, 0, 45)
	isnap.Save(colourFiltered.(image.Image), "out_colorfiltered.jpg", "jpg")
	sharpened := isnap.Sharpen(colourFiltered, 3)
	isnap.Save(sharpened.(image.Image), "out_sharpened.jpg", "jpg")
	
	feature := isnap.SingleIsolate(sharpened)
	fmt.Println(feature)
	writeable := isnap.MakeNRGBA(sharpened)
	isnap.DrawFeature(feature, &writeable)
	
	var outImg image.Image = &writeable
	err = isnap.Save(outImg.(image.Image), "out.jpg", "jpg")
	fmt.Println(err)
}
