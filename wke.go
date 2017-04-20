package wke

import (
	"unsafe"
)

// #cgo LDFLAGS: -L${SRCDIR} -lwke
// # include <stdlib.h>
// #include "wke.h"
import "C"

type Rect struct {
	x, y, w, h int
}

// GoBool convert C.bool to bool
func GoBool(b C.bool) bool {
	if b == 0 {
		return false
	}
	return true
}

// CBool convert bool to C.bool
func CBool(b bool) C.bool {
	if b {
		return 1
	}
	return 0
}

type WebView struct {
	v C.wkeWebView
}

func (w WebView) Destroy() {
	C.wkeDestroyWebView(w.v)
}

func (w WebView) Name() string {
	return C.GoString(C.wkeWebViewName(w.v))
}

func (w WebView) SetName(name string) {
	s := C.CString(name)
	C.wkeSetWebViewName(w.v, s)
	C.free(unsafe.Pointer(s))
}

func (w WebView) IsTransparent() bool {
	b := C.wkeIsTransparent(w.v)
	return GoBool(b)
}

func (w WebView) SetTransparent(transparent bool) {
	C.wkeSetTransparent(w.v, CBool(transparent))
}

func (w WebView) LoadURL(url string) {
	s := C.CString(url)
	C.wkeLoadURL(w.v, (*C.utf8)(s))
	C.free(unsafe.Pointer(s))
}

func (w WebView) LoadHTML(html string) {
	s := C.CString(html)
	C.wkeLoadHTML(w.v, (*C.utf8)(s))
	C.free(unsafe.Pointer(s))
}

func (w WebView) LoadFile(filename string) {
	s := C.CString(filename)
	C.wkeLoadHTML(w.v, (*C.utf8)(s))
	C.free(unsafe.Pointer(s))
}

func (w WebView) IsLoaded() bool {
	b := C.wkeIsLoaded(w.v)
	return GoBool(b)
}

func (w WebView) IsLoadFailed() bool {
	b := C.wkeIsLoadFailed(w.v)
	return GoBool(b)
}

func (w WebView) IsLoadComplete() bool {
	b := C.wkeIsLoadComplete(w.v)
	return GoBool(b)
}

func (w WebView) IsDocumentReady() bool {
	b := C.wkeIsDocumentReady(w.v)
	return GoBool(b)
}

func (w WebView) StopLoading() {
	C.wkeStopLoading(w.v)
}

func (w WebView) Reload() {
	C.wkeReload(w.v)
}

func (w WebView) Title() string {
	return C.GoString((*C.char)(C.wkeTitle(w.v)))
}

func (w WebView) Resize(width, height int) {
	C.wkeResize(w.v, C.int(width), C.int(height))
}

func (w WebView) Width() int {
	return int(C.wkeWidth(w.v))
}

func (w WebView) Height() int {
	return int(C.wkeHeight(w.v))
}

func (w WebView) ContentsWidth() int {
	return int(C.wkeContentsWidth(w.v))
}

func (w WebView) ContentsHeight() int {
	return int(C.wkeContentsHeight(w.v))
}

func (w WebView) SetDirty(dirty bool) {
	C.wkeSetDirty(w.v, CBool(dirty))
}

func (w WebView) IsDirty() bool {
	return GoBool(C.wkeIsDirty(w.v))
}

func (w WebView) AddDirtyArea(x, y, width, height int) {
	C.wkeAddDirtyArea(w.v, C.int(x), C.int(y), C.int(width), C.int(height))
}

func (w WebView) LayoutIfNeeded() {
	C.wkeLayoutIfNeeded(w.v)
}

// Paint paints view's content as memory block
func (w WebView) Paint(b []byte) []byte {
	width := w.ContentsWidth()
	height := w.ContentsHeight()
	wanted := width * height * 4
	if len(b) < wanted {
		b = make([]byte, wanted)
	} else {
		b = b[:wanted]
	}

	C.wkePaint(w.v, unsafe.Pointer(&b[0]), 0)
	return b
}

func (w WebView) CanGoBack() bool {
	return GoBool(C.wkeCanGoBack(w.v))
}

func (w WebView) GoBack() bool {
	return GoBool(C.wkeGoBack(w.v))
}

func (w WebView) CanGoForward() bool {
	return GoBool(C.wkeCanGoForward(w.v))
}

func (w WebView) GoForward() bool {
	return GoBool(C.wkeGoForward(w.v))
}

func (w WebView) SelectAll() {
	C.wkeSelectAll(w.v)
}

func (w WebView) Copy() {
	C.wkeCopy(w.v)
}

func (w WebView) Cut() {
	C.wkeCut(w.v)
}

func (w WebView) Paste() {
	C.wkePaste(w.v)
}

func (w WebView) Delete() {
	C.wkeDelete(w.v)
}

func (w WebView) SetCookieEnabled(enable bool) {
	C.wkeSetCookieEnabled(w.v, CBool(enable))
}

func (w WebView) CookieEnabled() bool {
	return GoBool(C.wkeCookieEnabled(w.v))
}

func (w WebView) SetMediaVolume(volume float32) {
	C.wkeSetMediaVolume(w.v, C.float(volume))
}

func (w WebView) MediaVolume() float32 {
	return float32(C.wkeMediaVolume(w.v))
}

func (w WebView) MouseEvent(message uint, x, y int, flags uint) bool {
	return GoBool(C.wkeMouseEvent(w.v, C.uint(message), C.int(x), C.int(y), C.uint(flags)))
}

func (w WebView) ContextMenuEvent(x, y int, flags uint) bool {
	return GoBool(C.wkeContextMenuEvent(w.v, C.int(x), C.int(y), C.uint(flags)))
}

func (w WebView) MouseWheel(x, y, delta int, flags uint) bool {
	return GoBool(C.wkeMouseWheel(w.v, C.int(x), C.int(y), C.int(delta), C.uint(flags)))
}

func (w WebView) KeyUp(keyCode uint, flags uint, systemKey bool) bool {
	return GoBool(C.wkeKeyUp(w.v, C.uint(keyCode), C.uint(flags), CBool(systemKey)))
}

func (w WebView) KeyDown(keyCode uint, flags uint, systemKey bool) bool {
	return GoBool(C.wkeKeyDown(w.v, C.uint(keyCode), C.uint(flags), CBool(systemKey)))
}

func (w WebView) KeyPress(keyCode uint, flags uint, systemKey bool) bool {
	return GoBool(C.wkeKeyPress(w.v, C.uint(keyCode), C.uint(flags), CBool(systemKey)))
}

func (w WebView) Focus() {
	C.wkeFocus(w.v)
}

func (w WebView) Unfocus() {
	C.wkeUnfocus(w.v)
}

func (w WebView) GetCaret() Rect {
	rect := C.wkeGetCaret(w.v)
	return Rect{
		int(rect.x),
		int(rect.y),
		int(rect.w),
		int(rect.h),
	}
}

func (w WebView) RunJS(script string) JsValue {
	s := C.CString(script)
	v := C.wkeRunJS(w.v, (*C.utf8)(s))
	C.free(unsafe.Pointer(s))
	return JsValue(v)
}

func (w WebView) GlobalExec() JsExecState {
	return JsExecState{
		s: C.wkeGlobalExec(w.v),
	}
}

func (w WebView) Sleep() {
	C.wkeSleep(w.v)
}

func (w WebView) Awaken() {
	C.wkeAwaken(w.v)
}

func (w WebView) IsAwake() bool {
	return GoBool(C.wkeIsAwake(w.v))
}

func (w WebView) SetZoomFactor(factor float32) {
	C.wkeSetZoomFactor(w.v, C.float(factor))
}

func (w WebView) ZoomFactor() float32 {
	return float32(C.wkeZoomFactor(w.v))
}

func (w WebView) SetEditable(editable bool) {
	C.wkeSetEditable(w.v, CBool(editable))
}

func (w WebView) SetClientHandler(handler *C.wkeClientHandler) {
	C.wkeSetClientHandler(w.v, handler)
}

func (w WebView) GetClientHandler() *C.wkeClientHandler {
	return C.wkeGetClientHandler(w.v)
}

// NewWebView create a new webview
func NewWebView() WebView {
	return WebView{v: C.wkeCreateWebView()}
}

// GetWebView find webview by name
func GetWebView(name string) WebView {
	s := C.CString(name)
	v := C.wkeGetWebView(s)
	C.free(unsafe.Pointer(s))
	return WebView{v: v}
}

// init wke environment
func init() {
	C.wkeInit()
}

func Shutdown() {
	C.wkeShutdown()
}

func Update() {
	C.wkeUpdate()
}

func Version() uint {
	return uint(C.wkeVersion())
}

func VersionString() string {
	return C.GoString((*C.char)(C.wkeVersionString()))
}
