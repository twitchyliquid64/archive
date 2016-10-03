package isnap

import (
    "image"
    "image/color"
    "image/draw"
    "math"
    "fmt"
)


//credit: http://rosettacode.org/wiki/Hough_transform#Go
func Hough(im image.Image, ntx, mry int) image.Image {//probably not working
    nimx := im.Bounds().Max.X
    mimy := im.Bounds().Max.Y
    mry = int(mry/2) * 2
    him := image.NewGray(image.Rect(0, 0, ntx, mry))
    draw.Draw(him, him.Bounds(), image.NewUniform(color.White),
        image.ZP, draw.Src)
 
    rmax := math.Hypot(float64(nimx), float64(mimy))
    dr := rmax / float64(mry/2)
    dth := math.Pi / float64(ntx)
 
    for jx := 0; jx < nimx; jx++ {
        for iy := 0; iy < mimy; iy++ {
            col := color.GrayModel.Convert(im.At(jx, iy)).(color.Gray)
            if col.Y == 255 {
                continue
            }
            for jtx := 0; jtx < ntx; jtx++ {
                th := dth * float64(jtx)
                r := float64(jx)*math.Cos(th) + float64(iy)*math.Sin(th)
                iry := mry/2 - int(math.Floor(r/dr+.5))
                col = him.At(jtx, iry).(color.Gray)
                if col.Y > 0 {
                    col.Y--
                    him.SetGray(jtx, iry, col)
                }
            }
        }
    }
    
    //average, st_dev := calcStats(him)
    //Threshold(him, uint8(average-int(st_dev/5)), 255, 0)//set the threshold as average+(standard deviation/5)
    Threshold(him, Otsu(him), 255, 0)
    
    //convert the co-ordinates back to cartesians
    //him = cartesian(him)
    
    return him
}


func calcStats(m *image.Gray)(av int, std float64){
	count := len(m.Pix)
	average := uint64(0)
	for i := 0; i < count; i++ {
		average += uint64(m.Pix[i])
	}
	average = average/uint64(count)
	
	
	//second iteration to calc culmative variance squared
	variance := uint64(0)
	for i := 0; i < count; i++ {
		variance += (average-uint64(m.Pix[i])) * (average-uint64(m.Pix[i]));
	}
	std = math.Sqrt(float64(variance/uint64(count)))
	
	
	av = int(average)
	return
}


// All color values greater than threshold are set to fgColor, the other to bgColor.
func Threshold(m *image.Gray, threshold, bgColor, fgColor uint8) {
	count := len(m.Pix)
	for i := 0; i < count; i++ {
		//dist := math.Abs(float64(m.Pix[i]-threshold))
		//m.Pix[i] = uint8(float64(fgColor)/(dist+1) + (float64(bgColor)*dist))
		
		if m.Pix[i] < threshold{
			m.Pix[i] = fgColor
		}else{
			m.Pix[i] = bgColor
		}
	}
}


// Otsu determines a threshold value using Otsu's method on grayscale images.
//
// "A threshold selection method from gray-level histograms" (Otsu, 1979).
func Otsu(m *image.Gray) uint8 {
	hist := Histogram(m)
	sum := 0
	for i, v := range hist {
		sum += i * v
	}
	wB, wF := 0, len(m.Pix)
	sumB, sumF := 0, sum
	maxVariance := 0.0
	threshold := uint8(0)
	for t := 0; t < 256; t++ {
		wB += hist[t]
		wF -= hist[t]
		if wB == 0 {
			continue
		}
		if wF == 0 {
			return threshold
		}
		sumB += t * hist[t]
		sumF = sum - sumB
		mB := float64(sumB) / float64(wB)
		mF := float64(sumF) / float64(wF)
		betweenVariance := float64(wB*wF) * (mB - mF) * (mB - mF)
		if betweenVariance > maxVariance {
			maxVariance = betweenVariance
			threshold = uint8(t)
		}
	}
	return threshold
}

// Histogram creates a histogram of a grayscale image.
func Histogram(m *image.Gray) []int {
	hist := make([]int, 256)
	count := len(m.Pix)
	for i := 0; i < count; i++ {
		hist[m.Pix[i]]++
	}
	return hist
}


func cartesian (m *image.Gray)*image.Gray{//probably not working
	b := m.Bounds()
	
    him := image.NewGray(b)
    draw.Draw(him, him.Bounds(), image.NewUniform(color.White),
        image.ZP, draw.Src)
	

	for y := b.Min.Y; y < b.Max.Y; y++ {
	 for x := b.Min.X; x < b.Max.X; x++ {
	  val := m.At(x, y)
	  xc := int(float64(y) * math.Cos(float64(x)))
	  yc := int(float64(y) * math.Sin(float64(x)))
		if xc>0 && yc>0{
		  him.SetGray(xc, yc, val.(color.Gray))
		  fmt.Println(xc, yc, val)
		}
	 }
	}
	
	return him
}
