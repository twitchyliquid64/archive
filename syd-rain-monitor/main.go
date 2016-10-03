package main

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/go-martini/martini"
	"html/template"
	"net/http"
	"image/png"
	"strings"
	"image"
	"time"
	"fmt"
)

const sydLegend =
"http://ws.cdn.bom.gov.au/products/radar_transparencies/IDR714.locations.png"
const homeX = 257
const homeY = 333

var web *martini.ClassicMartini
var prog IProgEngine

func getImage(url string)(a image.Image, b error){
	res, err := http.Get(url)
	if err != nil {
		return image.NewAlpha(image.ZR), err
	}

	img, err := png.Decode(res.Body)
	if err != nil{
		return image.NewAlpha(image.ZR), err
	}
	return img, nil
}

func run(){
	d, err := getRadarData()
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	fmt.Println("Map Resolution:", d["km"]+"km")
	fmt.Println("Center latitude:", d["lat"])
	fmt.Println("Center longitude:", d["lon"])

	updated = strings.TrimPrefix(d["img6"], "http://ws.cdn.bom.gov.au/radar/")[17:21]

	img6, err := getImage(d["img6"])
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	evaluateImage(img6)
	prog.Feed(zones)

	outimg := MakeRGBA(img6)
	gc := draw2dimg.NewGraphicContext(&outimg)

	drawBase(gc)

	//drawZones(gc)

	lastImage = outimg
	draw2dimg.SaveToPngFile("test.png", &outimg)
}

func servImg(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "image/png")
	enc := png.Encoder{CompressionLevel: png.BestSpeed}
	enc.Encode(res, &lastImage)
}

func servPage(res http.ResponseWriter, req *http.Request){
	t, err := template.ParseFiles("page.html")
	err = t.ExecuteTemplate(res, "page.html", struct{
		Zones []*Zone
		Updated string
		Prognosis IProgEngine
	}{
		Zones: zones,
		Updated: updated,
		Prognosis: prog,
	})
	if err != nil{
		fmt.Println(err)
	}
}

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

func main(){
	fmt.Println("Sydney rain monitor\nv0.1 Alpha - twitchyliquid64\n")

	prog = &ZoningPrognosis{}

	web = martini.Classic()
	web.Get("/img", servImg)
	web.Get("/", servPage)

	run()
	go doEvery(time.Minute * 3, run)
	web.Run()
}
