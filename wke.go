package wke

import (
	"sync"
	"unsafe"
)

// #cgo LDFLAGS: '--enable-stdcall-fixup'
// #cgo LDFLAGS: -L${SRCDIR} -lwke
// # include <stdlib.h>
// #include "wke.h"
// extern void titleChangedCallback(wkeWebView* webView, void* param, const wkeString* title);
// extern void urlChangedCallback(wkeWebView* webView, void* param, const wkeString* url);
import "C"

type Rect struct {
	X, Y, W, H int32
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

type TitleChangedCallback func(title string)
type URLChangedCallback func(title string)

var webViewCallbacks struct {
	sync.Mutex
	cbs map[*C.wkeWebView]map[string]interface{}
}

func setWebViewCallback(v *C.wkeWebView, name string, f interface{}) {
	webViewCallbacks.Lock()
	defer webViewCallbacks.Unlock()
	if cb, ok := webViewCallbacks.cbs[v]; ok {
		cb[name] = f
	} else {
		webViewCallbacks.cbs[v] = map[string]interface{}{
			name: f,
		}
	}
}

func getWebViewCallback(v *C.wkeWebView, name string) interface{} {
	webViewCallbacks.Lock()
	defer webViewCallbacks.Unlock()
	if cb, ok := webViewCallbacks.cbs[v]; ok {
		return cb[name]
	}
	return nil
}

type WebView struct {
	v *C.wkeWebView
}

func (w *WebView) Destroy() {
	C.wkeDestroyWebView(w.v)
}

func (w *WebView) Name() string {
	return C.GoString(C.wkeGetName(w.v))
}

func (w *WebView) SetName(name string) {
	s := C.CString(name)
	C.wkeSetName(w.v, s)
	C.free(unsafe.Pointer(s))
}

func (w *WebView) IsTransparent() bool {
	b := C.wkeIsTransparent(w.v)
	return GoBool(b)
}

func (w *WebView) SetTransparent(transparent bool) {
	C.wkeSetTransparent(w.v, CBool(transparent))
}

func (w *WebView) SetUserAgent(agent string) {
	s := C.CString(agent)
	C.wkeSetUserAgent(w.v, (*C.utf8)(s))
	C.free(unsafe.Pointer(s))
}

func (w *WebView) LoadURL(url string) {
	s := C.CString(url)
	C.wkeLoadURL(w.v, (*C.utf8)(s))
	C.free(unsafe.Pointer(s))
}

func (w *WebView) LoadHTML(html string) {
	s := C.CString(html)
	C.wkeLoadHTML(w.v, (*C.utf8)(s))
	C.free(unsafe.Pointer(s))
}

func (w *WebView) LoadFile(filename string) {
	s := C.CString(filename)
	C.wkeLoadHTML(w.v, (*C.utf8)(s))
	C.free(unsafe.Pointer(s))
}

func (w *WebView) Load(path string) {
	s := C.CString(path)
	C.wkeLoad(w.v, (*C.utf8)(s))
	C.free(unsafe.Pointer(s))
}

/*
func (w *WebView) IsLoading() bool {
	b := C.wkeIsLoading(w.v)
	return GoBool(b)
}
*/

func (w *WebView) IsLoadingSucceeded() bool {
	b := C.wkeIsLoadingSucceeded(w.v)
	return GoBool(b)
}

func (w *WebView) IsLoadingFailed() bool {
	b := C.wkeIsLoadingFailed(w.v)
	return GoBool(b)
}

func (w *WebView) IsLoadingCompleted() bool {
	b := C.wkeIsLoadingCompleted(w.v)
	return GoBool(b)
}

func (w *WebView) IsDocumentReady() bool {
	b := C.wkeIsDocumentReady(w.v)
	return GoBool(b)
}

func (w *WebView) StopLoading() {
	C.wkeStopLoading(w.v)
}

func (w *WebView) Reload() {
	C.wkeReload(w.v)
}

func (w *WebView) Title() string {
	return C.GoString((*C.char)(C.wkeGetTitle(w.v)))
}

func (w *WebView) Resize(width, height int) {
	C.wkeResize(w.v, C.int(width), C.int(height))
}

func (w *WebView) Width() int {
	return int(C.wkeGetWidth(w.v))
}

func (w *WebView) Height() int {
	return int(C.wkeGetHeight(w.v))
}

func (w *WebView) ContentsWidth() int {
	return int(C.wkeGetContentWidth(w.v))
}

func (w *WebView) ContentsHeight() int {
	return int(C.wkeGetContentHeight(w.v))
}

func (w *WebView) SetDirty(dirty bool) {
	C.wkeSetDirty(w.v, CBool(dirty))
}

func (w *WebView) IsDirty() bool {
	return GoBool(C.wkeIsDirty(w.v))
}

func (w *WebView) AddDirtyArea(x, y, width, height int) {
	C.wkeAddDirtyArea(w.v, C.int(x), C.int(y), C.int(width), C.int(height))
}

func (w *WebView) LayoutIfNeeded() {
	C.wkeLayoutIfNeeded(w.v)
}

// Paint paints view's content as a RGBA pixel block
func (w *WebView) Paint(b []byte) []byte {
	width := w.Width()
	height := w.Height()
	wanted := width * height * 4
	if len(b) < wanted {
		b = make([]byte, wanted)
	} else {
		b = b[:wanted]
	}

	if len(b) > 0 {
		C.wkePaint2(w.v, unsafe.Pointer(&b[0]), 0)
	}

	return b
}

// BMP images are stored in BGRA order rather than RGBA order.
func (w WebView) PaintNRGBA(b []byte) []byte {
	b = w.Paint(b)
	// convert from bgra to rgba
	width := w.Width()
	height := w.Height()
	stride := width * 4
	for y := 0; y != height; y += 1 {
		p := b[y*stride : y*stride+stride]
		for i := 0; i < len(p); i += 4 {
			p[i+0], p[i+2] = p[i+2], p[i+0]
		}
	}
	return b
}

// RepaintIfNeeded repaint webview to low-level hdc if needed
func (w *WebView) RepaintIfNeeded() bool {
	return GoBool(C.wkeRepaintIfNeeded(w.v))
}

func (w *WebView) GetViewDC() unsafe.Pointer {
	return C.wkeGetViewDC(w.v)
}

func (w *WebView) SetRepaintInterval(ms int) {
	C.wkeSetRepaintInterval(w.v, C.int(ms))
}

func (w *WebView) GetRepaintInterval() int {
	return int(C.wkeGetRepaintInterval(w.v))
}

func (w *WebView) RepaintIfNeededAfterInterval() bool {
	return GoBool(C.wkeRepaintIfNeededAfterInterval(w.v))
}

func (w *WebView) CanGoBack() bool {
	return GoBool(C.wkeCanGoBack(w.v))
}

func (w *WebView) GoBack() bool {
	return GoBool(C.wkeGoBack(w.v))
}

func (w *WebView) CanGoForward() bool {
	return GoBool(C.wkeCanGoForward(w.v))
}

func (w *WebView) GoForward() bool {
	return GoBool(C.wkeGoForward(w.v))
}

func (w *WebView) EditorSelectAll() {
	C.wkeEditorSelectAll(w.v)
}

func (w *WebView) EditorCopy() {
	C.wkeEditorCopy(w.v)
}

func (w *WebView) EditorCut() {
	C.wkeEditorCut(w.v)
}

func (w *WebView) EditorPaste() {
	C.wkeEditorPaste(w.v)
}

func (w *WebView) EditorDelete() {
	C.wkeEditorDelete(w.v)
}

func (w *WebView) GetCookie() string {
	return C.GoString((*C.char)(C.wkeGetCookie(w.v)))
}

func (w *WebView) SetCookieEnabled(enable bool) {
	C.wkeSetCookieEnabled(w.v, CBool(enable))
}

func (w *WebView) IsCookieEnabled() bool {
	return GoBool(C.wkeIsCookieEnabled(w.v))
}

func (w *WebView) SetMediaVolume(volume float32) {
	C.wkeSetMediaVolume(w.v, C.float(volume))
}

func (w *WebView) MediaVolume() float32 {
	return float32(C.wkeGetMediaVolume(w.v))
}

func (w *WebView) FireMouseEvent(message MouseMsg, x, y int, flags MouseFlags) bool {
	return GoBool(C.wkeFireMouseEvent(w.v, C.uint(message), C.int(x), C.int(y), C.uint(flags)))
}

func (w *WebView) FireContextMenuEvent(x, y int, flags MouseFlags) bool {
	return GoBool(C.wkeFireContextMenuEvent(w.v, C.int(x), C.int(y), C.uint(flags)))
}

func (w *WebView) FireMouseWheelEvent(x, y, delta int, flags MouseFlags) bool {
	return GoBool(C.wkeFireMouseWheelEvent(w.v, C.int(x), C.int(y), C.int(delta), C.uint(flags)))
}

func (w *WebView) FireKeyUpEvent(keyCode uint, flags KeyFlags, systemKey bool) bool {
	return GoBool(C.wkeFireKeyUpEvent(w.v, C.uint(keyCode), C.uint(flags), CBool(systemKey)))
}

func (w *WebView) FireKeyDownEvent(keyCode uint, flags KeyFlags, systemKey bool) bool {
	return GoBool(C.wkeFireKeyDownEvent(w.v, C.uint(keyCode), C.uint(flags), CBool(systemKey)))
}

func (w *WebView) FireKeyPressEvent(keyCode uint, flags KeyFlags, systemKey bool) bool {
	return GoBool(C.wkeFireKeyPressEvent(w.v, C.uint(keyCode), C.uint(flags), CBool(systemKey)))
}

func (w *WebView) SetFocus() {
	C.wkeSetFocus(w.v)
}

func (w *WebView) KillFocus() {
	C.wkeKillFocus(w.v)
}

func (w *WebView) GetCaretRect() Rect {
	rect := C.wkeGetCaretRect(w.v)
	return Rect{
		int32(rect.x),
		int32(rect.y),
		int32(rect.w),
		int32(rect.h),
	}
}

func (w *WebView) RunJS(script string) JSValue {
	s := C.CString(script)
	v := C.wkeRunJS(w.v, (*C.utf8)(s))
	C.free(unsafe.Pointer(s))
	return JSValue(v)
}

func (w *WebView) GlobalExec() *JSState {
	return &JSState{
		s: C.wkeGlobalExec(w.v),
	}
}

func (w *WebView) Sleep() {
	C.wkeSleep(w.v)
}

func (w *WebView) Awaken() {
	C.wkeWake(w.v)
}

func (w *WebView) IsAwake() bool {
	return GoBool(C.wkeIsAwake(w.v))
}

func (w *WebView) SetZoomFactor(factor float32) {
	C.wkeSetZoomFactor(w.v, C.float(factor))
}

func (w *WebView) ZoomFactor() float32 {
	return float32(C.wkeGetZoomFactor(w.v))
}

func (w *WebView) SetEditable(editable bool) {
	C.wkeSetEditable(w.v, CBool(editable))
}

//export goTitleChanged
func goTitleChanged(v *C.wkeWebView, title *C.char) {
	f := getWebViewCallback(v, "TitleChanged").(TitleChangedCallback)
	if f != nil {
		f(C.GoString(title))
	}
}

//export goURLChanged
func goURLChanged(v *C.wkeWebView, url *C.char) {
	f := getWebViewCallback(v, "URLChanged").(URLChangedCallback)
	if f != nil {
		f(C.GoString(url))
	}
}

func (w *WebView) SetTitleChanged(f TitleChangedCallback) {
	setWebViewCallback(w.v, "TitleChanged", f)
	C.wkeOnTitleChanged(w.v, C.wkeTitleChangedCallback(C.titleChangedCallback), nil)
}

func (w *WebView) SetURLChanged(f URLChangedCallback) {
	setWebViewCallback(w.v, "URLChanged", f)
	C.wkeOnURLChanged(w.v, C.wkeURLChangedCallback(C.urlChangedCallback), nil)
}

// NewWebView create a new webview
func NewWebView() *WebView {
	v := C.wkeCreateWebView()
	return &WebView{v: v}
}

// GetWebView find webview by name
func GetWebView(name string) *WebView {
	s := C.CString(name)
	v := C.wkeGetWebView(s)
	C.free(unsafe.Pointer(s))
	return &WebView{v: v}
}

// init wke environment
func Initialize() {
	C.wkeInitialize()
}

func Finalize() {
	C.wkeFinalize()
}

func Update() {
	C.wkeUpdate()
}

func Version() uint {
	return uint(C.wkeGetVersion())
}

func VersionString() string {
	return C.GoString((*C.char)(C.wkeGetVersionString()))
}

func RepaintAllNeeded() bool {
	return GoBool(C.wkeRepaintAllNeeded())
}

func RunMessageLoop(b bool) int {
	cb := CBool(b)
	return int(C.wkeRunMessageLoop((*C.bool)(&cb)))
}

// init wke
func init() {
	Initialize()
	webViewCallbacks.cbs = make(map[*C.wkeWebView]map[string]interface{})
}
