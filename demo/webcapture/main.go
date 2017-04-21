package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/xjdrew/wke"
)

func savePng(file string, img image.Image) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

func main() {
	flag.Parse()
	var url string
	if flag.NArg() > 0 {
		url = flag.Arg(0)
	} else {
		url = "http://example.com"
	}

	// print wke version
	fmt.Println(wke.VersionString())

	webView := wke.NewWebView()
	webView.Resize(1024, 768)

	fmt.Printf("loading url %s ...\n", url)
	webView.LoadURL(url)
	for {
		wke.Update()
		if webView.IsLoadComplete() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// hidden scrollbar
	webView.RunJS("document.body.style.overflow='hidden'")

	w := webView.ContentsWidth()
	h := webView.ContentsHeight()
	webView.Resize(w, h)

	pixels := webView.Paint(nil)
	img := &image.RGBA{
		Pix:    pixels,
		Stride: w * 4,
		Rect: image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{w, h},
		},
	}
	savePng(webView.Title()+".png", img)

	// fini
	webView.Destroy()
	wke.Shutdown()
}
