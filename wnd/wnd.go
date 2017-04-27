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
	webView               *wke.WebView
	url                   string
	title                 string
	urlChangedPublisher   walk.EventPublisher
	titleChangedPublisher walk.EventPublisher
	done                  chan struct{}
}

func NewWkeWnd(parent walk.Container) (*WkeWnd, error) {
	ww := new(WkeWnd)

	if err := walk.InitWidget(
		ww,
		parent,
		wkeWkeWndWindowClass,
		win.WS_CHILD|win.WS_VISIBLE|win.WS_TABSTOP,
		0); err != nil {
		return nil, err
	}

	ww.webView = wke.NewWebView()
	ww.done = make(chan struct{})
	ww.webView.SetTitleChanged(func(title string) {
		ww.title = title
		ww.titleChangedPublisher.Publish()
	})
	ww.webView.SetURLChanged(func(url string) {
		ww.url = url
		ww.urlChangedPublisher.Publish()
	})

	ww.MustRegisterProperty("URL", walk.NewProperty(
		func() interface{} {
			url := ww.URL()
			return url
		},
		func(v interface{}) error {
			ww.SetURL(v.(string))
			return nil
		},
		ww.urlChangedPublisher.Event()))

	ww.MustRegisterProperty("Title", walk.NewProperty(
		func() interface{} {
			url := ww.Title()
			return url
		},
		nil,
		ww.titleChangedPublisher.Event()))

	go func() {
		// 每秒50帧的频率刷新
		ticker := time.NewTicker(20 * time.Millisecond)
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
	return walk.Size{Width: 100, Height: 100}
}

func (ww *WkeWnd) URL() string {
	return ww.url
}

func (ww *WkeWnd) SetURL(url string) {
	ww.webView.Load(url)
}

func (ww *WkeWnd) Title() string {
	return ww.title
}

func (ww *WkeWnd) URLChanged() *walk.Event {
	return ww.urlChangedPublisher.Event()
}

func (ww *WkeWnd) TitleChanged() *walk.Event {
	return ww.titleChangedPublisher.Event()
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
			hMemDC, r.Left, r.Top, win.SRCCOPY)
	case win.WM_GETDLGCODE:
		// for WM_CHAR
		// ugly hack: form call win.IsDialogMessage, filter char input
		return win.DLGC_WANTALLKEYS
	case win.WM_KEYDOWN:
		ww.webView.FireKeyDownEvent(uint(wParam), lParamToKeyFlags(lParam), false)
		return 0
	case win.WM_KEYUP:
		ww.webView.FireKeyUpEvent(uint(wParam), lParamToKeyFlags(lParam), false)
	case win.WM_CHAR:
		ww.webView.FireKeyPressEvent(uint(wParam), lParamToKeyFlags(lParam), false)
	case win.WM_LBUTTONDOWN:
		fallthrough
	case win.WM_MBUTTONDOWN:
		fallthrough
	case win.WM_RBUTTONDOWN:
		fallthrough
	case win.WM_LBUTTONDBLCLK:
		fallthrough
	case win.WM_MBUTTONDBLCLK:
		fallthrough
	case win.WM_RBUTTONDBLCLK:
		fallthrough
	case win.WM_LBUTTONUP:
		fallthrough
	case win.WM_MBUTTONUP:
		fallthrough
	case win.WM_RBUTTONUP:
		fallthrough
	case win.WM_MOUSEMOVE:
		if msg == win.WM_LBUTTONDOWN || msg == win.WM_MBUTTONDOWN || msg == win.WM_RBUTTONDOWN {
			win.SetFocus(hwnd)
		}
		ww.webView.FireMouseEvent(wke.MouseMsg(msg), int(win.GET_X_LPARAM(lParam)), int(win.GET_Y_LPARAM(lParam)), wParamToMouseFlags(wParam))

	case win.WM_MOUSEWHEEL:
		p := win.POINT{X: win.GET_X_LPARAM(lParam), Y: win.GET_Y_LPARAM(lParam)}
		win.ScreenToClient(hwnd, &p)
		ww.webView.FireMouseWheelEvent(int(p.X), int(p.Y), wParamToWheelDelta(wParam), wParamToMouseFlags(wParam))

	case WM_IME_STARTCOMPOSITION:
		caret := ww.webView.GetCaretRect()
		var form COMPOSITIONFORM
		form.Style = CFS_POINT
		form.CurrentPos.X = caret.X
		form.CurrentPos.Y = caret.Y + caret.H
		form.Area.Top = caret.Y
		form.Area.Bottom = caret.Y + caret.H
		form.Area.Left = caret.X
		form.Area.Right = caret.X + caret.W

		hIMC := ImmGetContext(hwnd)
		ImmSetCompositionWindow(hIMC, &form)
		ImmReleaseContext(hwnd, hIMC)
		return 0
	case win.WM_SETFOCUS:
		ww.webView.SetFocus()
	case win.WM_KILLFOCUS:
		ww.webView.KillFocus()
	}

	return ww.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
