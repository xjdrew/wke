package wnd

import (
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

	ww.MustRegisterProperty("URL", walk.NewProperty(
		func() interface{} {
			url := ww.URL()
			return url
		},
		func(v interface{}) error {
			return ww.SetURL(v.(string))
		},
		ww.urlChangedPublisher.Event()))

	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		for {
			select {
			case <-ww.done:
			case <-ticker.C:
				if ww.webView.IsDirty() {
					// trigger WM_PAINT and don't erase background
					win.InvalidateRect(ww.Handle(), nil, false)
				}
			}
		}
	}()

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
		wbcb := ww.WidgetBase.ClientBounds()
		ww.webView.Resize(wbcb.Width, wbcb.Height)
	case win.WM_PAINT:
		var ps win.PAINTSTRUCT
		win.BeginPaint(hwnd, &ps)
		win.EndPaint(hwnd, &ps)

		ww.WebView().RepaintIfNeeded()
		hMemDC := (win.HDC)(ww.webView.GetViewDC())

		hdc := win.GetDC(hwnd)
		defer win.ReleaseDC(hwnd, hdc)

		r := &ps.RcPaint
		win.BitBlt(hdc, r.Left, r.Top, r.Right-r.Left, r.Bottom-r.Top,
			hMemDC, 0, 0, win.SRCCOPY)

	}

	return ww.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
