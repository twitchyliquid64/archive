package isnap

import (
	"os"
	"math"
	"image"
	"errors"
	"image/jpeg"
	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
)

func Jpg(file string)(img image.Image, err error){
	r, err := os.Open(file)
	if err != nil{
		return
	}
	defer r.Close()
	img, err = jpeg.Decode(r)
	return
}


func Save(img image.Image, fname, format string)error{
	if format == "jpg"{
		w, err := os.OpenFile(fname, os.O_WRONLY | os.O_CREATE, 0777)
		if err != nil{
			return err
		}
		defer w.Close()
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 100,})
	}
	return errors.New("Unexpected save format")
}


func ColourFilter(src image.Image, rCentre, gCentre, bCentre, tolerance float32)image.Image{
	g := gift.New(
		gift.ColorFunc(
			func(r0, g0, b0, a0 float32) (r, g, b, a float32) {
				
				r0 = r0*255;
				g0 = g0*255;
				b0 = b0*255;
				
				if ((r0 > (rCentre-tolerance)) && (r0 < (rCentre+tolerance))) && ((g0 > (gCentre-tolerance)) && (g0 < (gCentre+tolerance))) && ((b0 > (bCentre-tolerance)) && (b0 < (bCentre+tolerance))){
					r = 1;
					g = 1;
					b = 1;
				}else{
					r = 0;
					g = 0;
					b = 0;
				}

				a = a0       // preserve the alpha channel
				return
			},
		),
	)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	return dst
}


func Edge(src image.Image)image.Image{
	g := gift.New(
			gift.Convolution(
			[]float32{
				-1, -1, -1,
				-1, 8, -1,
				-1, -1, -1,
			},
			false, false, false, 0.0,
		),
	)
	
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	return dst
}


func Blur(src image.Image, sigma float64)image.Image{
	return imaging.Blur(src, sigma)
}

func Sharpen(src image.Image, sigma float64)image.Image{
	return imaging.Sharpen(src, sigma)
}


type Parallelogram struct {
	Bounds 		image.Rectangle			//bounding box of the object.
	Angle		float64
	Center 		image.Point
	BottomLeft	image.Point
	TopLeft		image.Point
	BottomRight	image.Point
	TopRight	image.Point
}

func SingleIsolate(src image.Image)(feature Parallelogram){
	b := src.Bounds()
	var bottomLeft, topLeft,  bottomRight, topRight image.Point
	
	
	minX := b.Max.X
	minY := b.Max.Y
	maxX := 0
	maxY := 0
	

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r,g,b,_ := src.At(x, y).RGBA()
			
			if (r>200) && (g>200) && (b>200){
				if x<minX{minX=x}
				if y<minY{minY=y}
				
				if x>maxX{maxX=x}
				if y>maxY{maxY=y}
			}
			
		}
	}
	
	//find the bottom left point
	for y := b.Min.Y; y < b.Max.Y; y++ {
		r,g,b,_ := src.At(minX, y).RGBA()
		if (r>200) && (g>200) && (b>200){
			bottomLeft = image.Pt(minX, y)
			break
		}
	}
	//find the top right point
	for y := b.Min.Y; y < b.Max.Y; y++ {
		r,g,b,_ := src.At(maxX, y).RGBA()
		if (r>200) && (g>200) && (b>200){
			topRight = image.Pt(maxX, y)
			break
		}
	}
	
	//find the top left point
	for x := b.Min.X; x < b.Max.X; x++ {
		r,g,b,_ := src.At(x, minY).RGBA()
		if (r>200) && (g>200) && (b>200){
			topLeft = image.Pt(x, minY)
			break
		}
	}
	//find the bottom right point
	for x := b.Min.X; x < b.Max.X; x++ {
		r,g,b,_ := src.At(x, maxY).RGBA()
		if (r>200) && (g>200) && (b>200){
			bottomRight = image.Pt(x, maxY)
			break
		}
	}
	
	
	opp := bottomLeft.Y-topLeft.Y
	adj := topLeft.X-bottomLeft.X
	ang := float64(90)-math.Atan2(float64(opp),float64(adj))/math.Pi*180
	
	return Parallelogram{
		Bounds: image.Rect(minX, minY, maxX, maxY),
		Angle: ang,
		Center: image.Pt((minX+maxX)/2, (maxY+minY)/2),
		BottomLeft: bottomLeft,
		TopLeft: topLeft,
		BottomRight: bottomRight,
		TopRight: topRight,
	}
}
