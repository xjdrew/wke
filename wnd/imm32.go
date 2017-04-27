package wnd

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

const (
	WM_IME_STARTCOMPOSITION = 0x010D
)

const (
	CFS_DEFAULT        = 0x0000
	CFS_RECT           = 0x0001
	CFS_POINT          = 0x0002
	CFS_FORCE_POSITION = 0x0020
	CFS_CANDIDATEPOS   = 0x0040
	CFS_EXCLUDE        = 0x0080
)

type CANDIDATEFORM struct {
	Index      uint32
	Style      uint32
	CurrentPos win.POINT
	Area       win.RECT
}

type COMPOSITIONFORM struct {
	Style      uint32
	CurrentPos win.POINT
	Area       win.RECT
}

type HIMC uint32

var (
	immGetContext           uintptr
	immGetCandidateWindow   uintptr
	immSetCandidateWindow   uintptr
	immReleaseContext       uintptr
	immSetCompositionWindow uintptr
)

func init() {
	libimm32, err := syscall.LoadLibrary("imm32.dll")
	if err != nil {
		return
	}

	immGetContext = win.MustGetProcAddress(uintptr(libimm32), "ImmGetContext")
	immGetCandidateWindow = win.MustGetProcAddress(uintptr(libimm32), "ImmGetCandidateWindow")
	immSetCandidateWindow = win.MustGetProcAddress(uintptr(libimm32), "ImmSetCandidateWindow")
	immReleaseContext = win.MustGetProcAddress(uintptr(libimm32), "ImmReleaseContext")
	immSetCompositionWindow = win.MustGetProcAddress(uintptr(libimm32), "ImmSetCompositionWindow")
}

func ImmGetContext(hwnd win.HWND) HIMC {
	if immGetContext == 0 {
		return 0
	}
	ret, _, _ := syscall.Syscall(immGetContext, 1, uintptr(hwnd), 0, 0)
	return HIMC(ret)
}
func ImmGetCandidateWindow(himc HIMC, index uint32, form *CANDIDATEFORM) bool {
	if immGetCandidateWindow == 0 {
		return false
	}
	ret, _, err := syscall.Syscall(immGetCandidateWindow, 3, uintptr(himc), uintptr(index), uintptr(unsafe.Pointer(form)))
	log.Print("ImmGetCandidateWindow:", err, syscall.GetLastError())
	return ret != 0
}

func ImmSetCandidateWindow(himc HIMC, form *CANDIDATEFORM) bool {
	if immSetCandidateWindow == 0 {
		return false
	}
	ret, _, _ := syscall.Syscall(immSetCandidateWindow, 2, uintptr(himc), uintptr(unsafe.Pointer(form)), 0)
	return ret != 0
}

func ImmReleaseContext(hwnd win.HWND, himc HIMC) bool {
	if immReleaseContext == 0 {
		return false
	}
	ret, _, _ := syscall.Syscall(immReleaseContext, 1, uintptr(hwnd), uintptr(himc), 0)
	return ret != 0
}

func ImmSetCompositionWindow(himc HIMC, form *COMPOSITIONFORM) bool {
	if immSetCompositionWindow == 0 {
		return false
	}
	ret, _, _ := syscall.Syscall(immSetCompositionWindow, 2, uintptr(himc), uintptr(unsafe.Pointer(form)), 0)
	return ret != 0
}
