package main

import (
	"log"

	"github.com/xjdrew/win"
	"github.com/xjdrew/wke"
	"github.com/xjdrew/wke/iui"
	"github.com/xjdrew/wke/webwindow"
)

type MyWebWindow struct {
	*webwindow.WebWindow
}

func (ww *MyWebWindow) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_CLOSE:
		log.Println("------ WM_CLOSE")
		ww.ShowWindow(false)
		return 0
	}
	return ww.WebWindow.WndProc(hwnd, msg, wParam, lParam)
}

func NewMyWebWindow(parent iui.Window) (*MyWebWindow, error) {
	w, _ := webwindow.NewWebWindow(parent)
	ww := &MyWebWindow{
		WebWindow: w,
	}
	ww.SetWndProc(ww.WndProc)
	ww.SetStyle(0)

	lb := &win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: win.COLORREF(win.COLOR_ACTIVEBORDER)}
	hBrush := win.CreateBrushIndirect(lb)
	win.SetClassLongPtr(ww.Handle(), -10, uintptr(hBrush))
	return ww, nil
}

func main() {
	ww, _ := webwindow.NewWebWindow(nil)

	wwSub, _ := NewMyWebWindow(ww)
	wwSub.WebView().LoadFile("E:\\gospace\\src\\github.com\\xjdrew\\wke\\demo\\webwindow\\example.html")
	//ww.SetURL("baidu.com")

	wv := ww.WebView()
	wv.LoadFile("E:\\gospace\\src\\github.com\\xjdrew\\wke\\demo\\webwindow\\main.html")

	// gs := wv.GlobalExec()
	wke.JSBind("openWindow", func(s *wke.JSState) wke.JSValue {
		log.Println("---- in openWindow")
		wwSub.ShowWindow(true)
		return s.JSUndefined()
	})

	wke.JSBind("closeWindow", func(s *wke.JSState) wke.JSValue {
		log.Println("---- in closeWindow")
		wwSub.ShowWindow(false)
		return s.JSUndefined()
	})

	log.Println("--- create window")
	ww.ShowWindow(true)
	if err := iui.Run(); err != nil {
		log.Println(err)
	}
}
