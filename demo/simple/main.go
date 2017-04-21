package main

import (
	"syscall"
	"time"
	"unsafe"

	"github.com/lxn/win"
	//"github.com/xjdrew/wke"
)

func main() {
	const minWindowClass = `\o/ WKE_MainWindow_Class \o/`
	MustRegisterWindowClass(minWindowClass)
	hMainWnd := win.CreateWindowEx(
		0,
		syscall.StringToUTF16Ptr(minWindowClass),
		syscall.StringToUTF16Ptr("main"),
		win.WS_OVERLAPPEDWINDOW,
		0,
		0,
		100,
		100,
		0,
		0,
		win.GetModuleHandle(nil),
		nil,
	)
	win.ShowWindow(hMainWnd, win.SW_SHOW)
	win.UpdateWindow(hMainWnd)
	time.Sleep(100 * time.Second)
}

func defaultWndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	return
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
	wc.LpfnWndProc = syscall.NewCallback(defaultWndProc)
	wc.HInstance = hInst
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = win.COLOR_BTNFACE + 1
	wc.LpszClassName = syscall.StringToUTF16Ptr(className)

	if atom := win.RegisterClassEx(&wc); atom == 0 {
		panic("RegisterClassEx")
	}
}
