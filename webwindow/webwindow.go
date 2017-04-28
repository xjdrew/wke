package webwindow

import (
	"time"

	"github.com/xjdrew/win"
	"github.com/xjdrew/wke"
	"github.com/xjdrew/wke/iui"
)

const WebWindowClass = `\o/ WKE_WebWindow_Class \o/`

func init() {
	iui.MustRegisterWindowClass(WebWindowClass)
}

type WebWindow struct {
	*iui.WindowBase
	webView *wke.WebView
	url     string
	title   string
	done    chan struct{}
}

func NewWebWindow(parent iui.Window) (*WebWindow, error) {
	window := iui.NewWindowEx(WebWindowClass, parent, nil)
	webView := wke.NewWebView()
	done := make(chan struct{})

	ww := &WebWindow{
		WindowBase: window.(*iui.WindowBase),
		webView:    webView,
		done:       done,
	}

	ww.SetWndProc(ww.WndProc)
	ww.webView.SetTitleChanged(func(title string) {
		ww.title = title
		// ww.titleChangedPublisher.Publish()
	})

	ww.webView.SetURLChanged(func(url string) {
		ww.url = url
		// ww.urlChangedPublisher.Publish()
	})

	go func() {
		// 每秒50帧的频率刷新
		ticker := time.NewTicker(20 * time.Millisecond)
		for {
			select {
			case <-ww.done:
				break
			case <-ticker.C:
				if ww.webView != nil && ww.webView.IsDirty() {
					// trigger WM_PAINT and don't erase background
					win.InvalidateRect(ww.Handle(), nil, false)
				}
			}
		}
	}()

	return ww, nil
}

func (ww *WebWindow) Dispose() {
	close(ww.done)
	if ww.webView != nil {
		ww.webView.Destroy()
		ww.webView = nil
	}
}

func (ww *WebWindow) URL() string {
	return ww.url
}

func (ww *WebWindow) SetURL(url string) {
	ww.webView.Load(url)
}

func (ww *WebWindow) Title() string {
	return ww.title
}

/*
func (ww *WkeWnd) URLChanged() *walk.Event {
	return ww.urlChangedPublisher.Event()
}

func (ww *WkeWnd) TitleChanged() *walk.Event {
	return ww.titleChangedPublisher.Event()
}
*/

func (ww *WebWindow) WebView() *wke.WebView {
	return ww.webView
}

func (ww *WebWindow) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_SIZE, win.WM_SIZING:
		r := ww.GetClientRect()
		ww.webView.Resize(int(r.Right-r.Left), int(r.Bottom-r.Top))
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

	case win.WM_IME_STARTCOMPOSITION:
		caret := ww.webView.GetCaretRect()
		var form win.COMPOSITIONFORM
		form.Style = win.CFS_POINT
		form.CurrentPos.X = caret.X
		form.CurrentPos.Y = caret.Y + caret.H
		form.Area.Top = caret.Y
		form.Area.Bottom = caret.Y + caret.H
		form.Area.Left = caret.X
		form.Area.Right = caret.X + caret.W

		hIMC := win.ImmGetContext(hwnd)
		win.ImmSetCompositionWindow(hIMC, &form)
		win.ImmReleaseContext(hwnd, hIMC)
		return 0
	case win.WM_SETFOCUS:
		ww.webView.SetFocus()
	case win.WM_KILLFOCUS:
		ww.webView.KillFocus()
	}

	return iui.DefaultWndProc(hwnd, msg, wParam, lParam)
}
