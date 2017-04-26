package wnd

import (
	"image"
	"log"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/win"
	"github.com/xjdrew/wke"
)

const wkeWkeWndWindowClass = `\o/ WKE_WkeWnd_Class \o/`

func init() {
	walk.MustRegisterWindowClass(wkeWkeWndWindowClass)
}

type WkeWnd struct {
	walk.WidgetBase
	webView             *wke.WebView
	urlChangedPublisher walk.EventPublisher
	url                 string
	done                chan struct{}
	pixels              []byte
}

func NewWkeWnd(parent walk.Container) (*WkeWnd, error) {
	ww := new(WkeWnd)
	if err := walk.InitWidget(
		ww,
		parent,
		wkeWkeWndWindowClass,
		win.WS_VISIBLE,
		0); err != nil {
		return nil, err
	}

	ww.webView = wke.NewWebView()
	ww.done = make(chan struct{})

	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		for {
			select {
			case <-ww.done:
			case <-ticker.C:
				if ww.webView.IsDirty() {
					ww.Invalidate()
				}
			}
		}
	}()
	ww.MustRegisterProperty("URL", walk.NewProperty(
		func() interface{} {
			url := ww.URL()
			log.Println("--- get url:", url)
			return url
		},
		func(v interface{}) error {
			log.Println("--- set url:", v.(string))
			return ww.SetURL(v.(string))
		},
		ww.urlChangedPublisher.Event()))
	return ww, nil
}

func (ww *WkeWnd) Dispose() {
	close(ww.done)
	if ww.webView != nil {
		ww.webView.Destroy()
		ww.webView = nil
	}
	ww.WidgetBase.Dispose()
}

func (*WkeWnd) LayoutFlags() walk.LayoutFlags {
	return walk.ShrinkableHorz | walk.ShrinkableVert | walk.GrowableHorz | walk.GrowableVert | walk.GreedyHorz | walk.GreedyVert
}

func (*WkeWnd) SizeHint() walk.Size {
	return walk.Size{100, 100}
}

func (ww *WkeWnd) URL() string {
	return ww.url
}

func (ww *WkeWnd) SetURL(url string) error {
	ww.url = url
	ww.webView.LoadURL(url)
	ww.urlChangedPublisher.Publish()
	return nil
}

func (ww *WkeWnd) URLChanged() *walk.Event {
	return ww.urlChangedPublisher.Event()
}

func (ww *WkeWnd) WebView() *wke.WebView {
	return ww.webView
}

func (ww *WkeWnd) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_SIZE, win.WM_SIZING:
		ww.Invalidate()
	case win.WM_PAINT:
		var ps win.PAINTSTRUCT
		hwnd := ww.Handle()
		win.BeginPaint(hwnd, &ps)
		defer win.EndPaint(hwnd, &ps)
		canvas, err := ww.CreateCanvas()
		if err != nil {
			panic(err)
		}
		defer canvas.Dispose()

		r := &ps.RcPaint
		w := int(r.Right - r.Left)
		h := int(r.Bottom - r.Top)

		ww.webView.Resize(w, h)
		pixels := ww.webView.PaintNRGBA(nil)
		img := &image.NRGBA{
			Pix:    pixels,
			Stride: w * 4,
			Rect: image.Rectangle{
				Min: image.Point{0, 0},
				Max: image.Point{w, h},
			},
		}
		if bitmap, err := walk.NewBitmapFromImage(img); err != nil {
			panic(err)
		} else {
			canvas.DrawImage(bitmap, walk.Point{X: int(r.Left), Y: int(r.Top)})
		}
	}

	return ww.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
