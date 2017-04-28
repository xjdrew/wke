package iui

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/xjdrew/win"
)

var (
	defaultWndProcPtr = syscall.NewCallback(DefaultWndProc)
)

type Window interface {
	Handle() win.HWND
	GetClientRect() win.RECT
	SetWndProc(f func(win.HWND, uint32, uintptr, uintptr) uintptr)
	ShowWindow(visible bool)
	SetStyle(style int32)
}

type WindowBase struct {
	hWnd win.HWND
}

func (w *WindowBase) Handle() win.HWND {
	return w.hWnd
}

func (w *WindowBase) GetClientRect() win.RECT {
	var r win.RECT
	if !win.GetClientRect(w.hWnd, &r) {
		// TODO: process error
		// lastError("GetClientRect")
	}
	return r
}

func (w *WindowBase) SetStyle(style int32) {
	win.SetWindowLong(w.hWnd, win.GWL_STYLE, int32(style))
}

func (w *WindowBase) SetWndProc(f func(win.HWND, uint32, uintptr, uintptr) uintptr) {
	win.SetWindowLongPtr(w.hWnd, win.GWLP_WNDPROC, syscall.NewCallback(f))
}

func (w *WindowBase) ShowWindow(visible bool) {
	var cmd int32
	if visible {
		cmd = win.SW_SHOW
	} else {
		cmd = win.SW_HIDE
	}
	win.ShowWindow(w.hWnd, cmd)
}

func DefaultWndProc(hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	// log.Printf("------ msg: %04x", msg)
	switch msg {
	case win.WM_CLOSE:
		log.Println("------ WM_CLOSE")
		win.PostQuitMessage(0)
	case win.WM_DESTROY:
		log.Println("------ WM_DESTROY")
	}
	return win.DefWindowProc(hWnd, msg, wParam, lParam)
}

func MustRegisterWindowClass(className string) {
	hInst := win.GetModuleHandle(nil)
	if hInst == 0 {
		panic("GetModuleHandle")
	}

	hIcon := win.LoadIcon(0, (*uint16)(unsafe.Pointer(uintptr(win.IDI_APPLICATION))))
	if hIcon == 0 {
		panic("LoadIcon")
	}

	hCursor := win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW))))
	if hCursor == 0 {
		panic("LoadCursor")
	}

	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.Style = win.CS_HREDRAW | win.CS_VREDRAW
	wc.LpfnWndProc = defaultWndProcPtr
	wc.HInstance = hInst
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = win.COLOR_BTNFACE + 1
	wc.LpszClassName = syscall.StringToUTF16Ptr(className)

	if atom := win.RegisterClassEx(&wc); atom == 0 {
		panic("RegisterClassEx")
	}
}

func Run() error {
	var msg win.MSG
	for {
		switch win.GetMessage(&msg, 0, 0, 0) {
		case 0: // WM_QUIT
			log.Println("------ WM_QUIT")
			return nil
		case -1: // error
			return syscall.GetLastError()
		default:
			win.TranslateMessage(&msg)
			win.DispatchMessage(&msg)
		}
	}
}

func NewWindow(className string) Window {
	return NewWindowEx(className, nil, nil)
}

func NewWindowEx(className string, parent Window, lpParam unsafe.Pointer) Window {
	// MustRegisterWindowClass(className)
	var hwndParent win.HWND
	if parent != nil {
		hwndParent = parent.Handle()
	}

	hWnd := win.CreateWindowEx(
		0,
		syscall.StringToUTF16Ptr(className),
		nil,
		win.WS_OVERLAPPEDWINDOW,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		hwndParent,
		0,
		0,
		lpParam)
	return &WindowBase{
		hWnd: hWnd,
	}
}
