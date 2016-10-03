package main

import (
        "github.com/llgcode/draw2d"
        "github.com/llgcode/draw2d/draw2dkit"
	"strconv"
	"image/color"
	"image"
	"fmt"
)

var lastImage image.RGBA
var updated string

func MakeRGBA(src image.Image)image.RGBA{
        b := src.Bounds()
        dst := image.NewRGBA(b)

        for y := b.Min.Y; y < b.Max.Y; y++ {
                for x := b.Min.X; x < b.Max.X; x++ {
                        dst.SetRGBA(x,y, color.RGBAModel.Convert(src.At(x, y)).(color.RGBA))
                }
        }

        return *dst
}

func scoreColor(c color.Color)int{
	r,g,b,a := c.RGBA()
	b = uint32(float64(b)/(0xffff) * 255) * 1
	g = uint32(float64(g)/(0xffff) * 255) * 4
	r = uint32(float64(r)/(0xffff) * 255) * 6

	maxPossible := (255*6) + (255*4) + 255
  if a == 0 {
    return 0
  }
	//fmt.Println(c)
	//fmt.Println(int(uint32(maxPossible)-r-g-b)-50)
	return int(uint32(maxPossible)-r-g-b)-50
}

func scoreHorizontalLine(src image.Image, startX, endX, Y int)int{
	score := 0
	for i := startX; i < endX; i++ {
		score += scoreColor(src.At(i,Y))
	}
	return score/(endX-startX)
}

func scoreVerticalLine(src image.Image, startY, endY, X int)int{
	score := 0
	for i := startY; i < endY; i++ {
		score += scoreColor(src.At(X,i))
	}
	return score/(endY-startY)
}


func scoreZone(src image.Image, distance int)int{
	score := scoreHorizontalLine(src, homeX-distance, homeX+distance, homeY-distance)
	score += scoreHorizontalLine(src, homeX-distance, homeX+distance, homeY+distance)
	score += scoreVerticalLine(src, homeY-distance, homeY+distance, homeX-distance)
	score += scoreVerticalLine(src, homeY-distance, homeY+distance, homeX+distance)
	return score / 4
}

func evaluateImage(src image.Image){
	for _, zone := range zones {
		score := scoreZone(src, zone.Offset)
		fmt.Println("Zone " + zone.Name + "(" + strconv.Itoa(zone.Offset) + "):", score)
		zone.Score = score

		if score > zone.Threshold {
			fmt.Println("Alert for zone " + zone.Name + "!: " + strconv.Itoa(score) + ">" + strconv.Itoa(zone.Threshold))
		}
	}
}

func drawZones(gc draw2d.GraphicContext){
	for _, zone := range zones {
		offset := zone.Offset

	        gc.SetLineWidth(2)
        	gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0x00})
		if zone.Triggered() {
		        gc.SetStrokeColor(color.RGBA{0xFF, 0x33, 0x00, 0xff})
		} else {
		        gc.SetStrokeColor(color.RGBA{0x33, 0xFF, 0x00, 0xff})
		}

        	draw2dkit.Rectangle(gc, float64(homeX-offset), float64(homeY-offset), float64(homeX+offset), float64(homeY+offset))
	        gc.FillStroke()
	}
}

func drawBase(gc draw2d.GraphicContext){
        legend, err := getImage(sydLegend)
        if err != nil {
                fmt.Println("ERROR:", err)
                return
        }

        gc.SetLineWidth(2)

	//legend
        gc.DrawImage(legend)

	//crosshairs
        gc.SetStrokeColor(color.RGBA{0x77, 0x00, 0x00, 0xff})
        gc.SetLineWidth(1)
	gc.MoveTo(0, homeY)
	gc.LineTo(512, homeY)
	gc.Close()
	gc.FillStroke()
	gc.MoveTo(homeX, 0)
	gc.LineTo(homeX, 512)
	gc.Close()
	gc.FillStroke()

	//top box
        gc.SetFillColor(color.RGBA{0xCC, 0xCC, 0xCC, 0xff})
        gc.SetStrokeColor(color.RGBA{0xCC, 0xCC, 0xCC, 0xff})
	draw2dkit.Rectangle(gc, 0, 0, 512, 20)
	gc.FillStroke()

	// Set the font luximbi.ttf
	draw2d.SetFontFolder("font")
	gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold })
	// Set the fill text color to black
	gc.SetFillColor(image.Black)
	gc.SetFontSize(14)
	// Display Topbar message
	gc.FillStringAt("Sydney Rain Monitor v0.1        twitchyliquid64", 5, 17)

	//home circle
        gc.SetLineWidth(2)
        gc.SetFillColor(color.RGBA{0xFF, 0x44, 0x44, 0xff})
        gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
        draw2dkit.Circle(gc, homeX, homeY, 3)
        gc.FillStroke()

}
