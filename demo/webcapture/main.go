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

	wke.Initialize()
	defer wke.Finalize()
	// print wke version
	fmt.Println(wke.VersionString())

	webView := wke.NewWebView()
	defer webView.Destroy()
	webView.Resize(800, 600)
	webView.SetTransparent(false)

	fmt.Printf("loading url %s ...\n", url)
	webView.LoadURL(url)
	for {
		wke.Update()
		if webView.IsLoadingCompleted() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// hidden scrollbar
	webView.RunJS("document.body.style.overflow='hidden'")

	w := webView.Width()
	h := webView.Height()
	pixels := webView.PaintNRGBA(nil)
	img := &image.NRGBA{
		Pix:    pixels,
		Stride: w * 4,
		Rect: image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{w, h},
		},
	}
	savePng(webView.Title()+".png", img)
}
