package isnap


import (
	"math"
	"image"
	"image/color"
)



func DrawFeature(f Parallelogram, img *image.NRGBA) {
	Drawpoint(f.Center, img, color.RGBA{255,40,40,255})
	Drawbox(f.Bounds, img, color.RGBA{40,40,255,255})
	
	Drawline(f.TopLeft.X, f.TopLeft.Y, f.TopRight.X, f.TopRight.Y, img, color.RGBA{0, 255,0,255})
	Drawline(f.BottomLeft.X, f.BottomLeft.Y, f.TopLeft.X, f.TopLeft.Y, img, color.RGBA{0, 255,0,255})
	Drawline(f.BottomLeft.X, f.BottomLeft.Y, f.BottomRight.X, f.BottomRight.Y, img, color.RGBA{0, 255,0,255})
	Drawline(f.BottomRight.X, f.BottomRight.Y, f.TopRight.X, f.TopRight.Y, img, color.RGBA{0, 255,0,255})
}

func Drawbox(bounds image.Rectangle, img *image.NRGBA, col color.Color) {
	Drawline(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y, img, col)
	Drawline(bounds.Max.X, bounds.Min.Y, bounds.Min.X, bounds.Max.Y, img, col)
	
	Drawline(bounds.Max.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y, img, col)
	Drawline(bounds.Max.X, bounds.Min.Y, bounds.Min.X, bounds.Min.Y, img, col)
	Drawline(bounds.Min.X, bounds.Min.Y, bounds.Min.X, bounds.Max.Y, img, col)
	Drawline(bounds.Max.X, bounds.Max.Y, bounds.Min.X, bounds.Max.Y, img, col)
}

func Drawpoint(p image.Point, img *image.NRGBA, col color.Color) {
	Drawline(p.X, p.Y-20, p.X, p.Y+20, img, col)
	Drawline(p.X-20, p.Y, p.X+20, p.Y, img, col)
}

func Drawline(x0, y0, x1, y1 int, img *image.NRGBA, col color.Color) {
	dx := math.Abs(float64(x1 - x0))
	dy := math.Abs(float64(y1 - y0))
	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy

	for {
		r, g, b, a := col.RGBA() 
		img.Pix[(y0-img.Rect.Min.Y)*img.Stride + (x0-img.Rect.Min.X)*4] = uint8(r/257)
		img.Pix[1 + (y0-img.Rect.Min.Y)*img.Stride + (x0-img.Rect.Min.X)*4] = uint8(g/257)
		img.Pix[2 + (y0-img.Rect.Min.Y)*img.Stride + (x0-img.Rect.Min.X)*4] = uint8(b/257)
		img.Pix[3 + (y0-img.Rect.Min.Y)*img.Stride + (x0-img.Rect.Min.X)*4] = uint8(a/257)
		
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}


func MakeNRGBA(src image.Image)image.NRGBA{
	b := src.Bounds()
	dst := image.NewNRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {

			dst.SetNRGBA(x,y, color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA))
		}
	}
	
	
	return *dst
}

