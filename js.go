package wke

import (
	"unsafe"
)

// #include <stdlib.h>
// #include "wke.h"
import "C"

type JSType int

const (
	JSTYPE_NUMBER JSType = iota
	JSTYPE_STRING
	JSTYPE_BOOLEAN
	JSTYPE_OBJECT
	JSTYPE_FUNCTION
	JSTYPE_UNDEFINED
)

type JSValue C.wkeJSValue

type JSState struct {
	s *C.wkeJSState
}

func (e *JSState) JSArgCount() int {
	return int(C.wkeJSParamCount(e.s))
}

func (e *JSState) JSArgType(argIdx int) JSType {
	return JSType(C.wkeJSParamType(e.s, C.int(argIdx)))
}

func (e *JSState) JSArg(argIdx int) JSValue {
	return JSValue(C.wkeJSParam(e.s, C.int(argIdx)))
}

// JSTypeOf returns value js type
func (e *JSState) JSTypeOf(v JSValue) JSType {
	return JSType(C.wkeJSTypeOf(e.s, C.wkeJSValue(v)))
}

func (e *JSState) JSIsNumber(v JSValue) bool {
	return GoBool(C.wkeJSIsNumber(e.s, C.wkeJSValue(v)))
}
func (e *JSState) JSIsString(v JSValue) bool {
	return GoBool(C.wkeJSIsString(e.s, C.wkeJSValue(v)))
}
func (e *JSState) JSIsBool(v JSValue) bool {
	return GoBool(C.wkeJSIsBool(e.s, C.wkeJSValue(v)))
}
func (e *JSState) JSIsObject(v JSValue) bool {
	return GoBool(C.wkeJSIsObject(e.s, C.wkeJSValue(v)))
}
func (e *JSState) JSIsFunction(v JSValue) bool {
	return GoBool(C.wkeJSIsFunction(e.s, C.wkeJSValue(v)))
}
func (e *JSState) JSIsUndefined(v JSValue) bool {
	return GoBool(C.wkeJSIsUndefined(e.s, C.wkeJSValue(v)))
}
func (e *JSState) JSIsNull(v JSValue) bool {
	return GoBool(C.wkeJSIsNull(e.s, C.wkeJSValue(v)))
}
func (e *JSState) JSIsArray(v JSValue) bool {
	return GoBool(C.wkeJSIsArray(e.s, C.wkeJSValue(v)))
}
func (e *JSState) JSIsTrue(v JSValue) bool {
	return GoBool(C.wkeJSIsTrue(e.s, C.wkeJSValue(v)))
}
func (e *JSState) JSIsFalse(v JSValue) bool {
	return GoBool(C.wkeJSIsFalse(e.s, C.wkeJSValue(v)))
}

func (e *JSState) JSToInt(v JSValue) int {
	return int(C.wkeJSToInt(e.s, C.wkeJSValue(v)))
}

func (e *JSState) JSToDouble(v JSValue) float64 {
	return float64(C.wkeJSToDouble(e.s, C.wkeJSValue(v)))
}

func (e *JSState) JSToBoolean(v JSValue) bool {
	return GoBool(C.wkeJSToBool(e.s, C.wkeJSValue(v)))
}

func (e *JSState) JSToString(v JSValue) string {
	return C.GoString((*C.char)(C.wkeJSToTempString(e.s, C.wkeJSValue(v))))
}

func (e *JSState) JSInt(n int) JSValue {
	return JSValue(C.wkeJSInt(e.s, C.int(n)))
}
func (e *JSState) JSDouble(f float64) JSValue {
	return JSValue(C.wkeJSDouble(e.s, C.double(f)))
}
func (e *JSState) JSBool(b bool) JSValue {
	return JSValue(C.wkeJSBool(e.s, CBool(b)))
}

func (e *JSState) JSUndefined() JSValue {
	return JSValue(C.wkeJSUndefined(e.s))
}
func (e *JSState) JSNull() JSValue {
	return JSValue(C.wkeJSNull(e.s))
}
func (e *JSState) JSTrue() JSValue {
	return JSValue(C.wkeJSTrue(e.s))
}
func (e *JSState) JSFalse() JSValue {
	return JSValue(C.wkeJSFalse(e.s))
}

func (e *JSState) JSString(s string) JSValue {
	cs := C.CString(s)
	v := C.wkeJSString(e.s, (*C.utf8)(cs))
	C.free(unsafe.Pointer(cs))
	return JSValue(v)
}

func (e *JSState) JSEmptyObject() JSValue {
	return JSValue(C.wkeJSEmptyObject(e.s))
}

func (e *JSState) JSEmptyArray() JSValue {
	return JSValue(C.wkeJSEmptyArray(e.s))
}

func (e *JSState) JSGet(object JSValue, prop string) JSValue {
	cs := C.CString(prop)
	v := C.wkeJSGet(e.s, C.wkeJSValue(object), cs)
	C.free(unsafe.Pointer(cs))
	return JSValue(v)
}

func (e *JSState) JSSet(object JSValue, prop string, v JSValue) {
	cs := C.CString(prop)
	C.wkeJSSet(e.s, C.wkeJSValue(object), cs, C.wkeJSValue(v))
	C.free(unsafe.Pointer(cs))
}

func (e *JSState) JSGetAt(object JSValue, index int) JSValue {
	v := C.wkeJSGetAt(e.s, C.wkeJSValue(object), C.int(index))
	return JSValue(v)
}

func (e *JSState) JSSetAt(object JSValue, index int, v JSValue) {
	C.wkeJSSetAt(e.s, C.wkeJSValue(object), C.int(index), C.wkeJSValue(v))
}

func (e *JSState) JSGetLength(object JSValue) int {
	return int(C.wkeJSGetLength(e.s, C.wkeJSValue(object)))
}

func (e *JSState) JSSetLength(object JSValue, length int) {
	C.wkeJSSetLength(e.s, C.wkeJSValue(object), C.int(length))
}

// JSGlobalObject returns window object
func (e *JSState) JSGlobalObject() JSValue {
	return JSValue(C.wkeJSGlobalObject(e.s))
}

// GetWebView returns WebView associated with this JSExecState
func (e *JSState) GetWebView() *WebView {
	return &WebView{
		v: C.wkeJSGetWebView(e.s),
	}
}

func (e *JSState) Eval(s string) JSValue {
	cs := C.CString(s)
	v := C.wkeJSEval(e.s, (*C.utf8)(cs))
	C.free(unsafe.Pointer(cs))
	return JSValue(v)
}

func (e *JSState) JSCall(f JSValue, thisObject JSValue, args []JSValue) JSValue {
	v := C.wkeJSCall(e.s, C.wkeJSValue(f), C.wkeJSValue(thisObject), (*C.wkeJSValue)(&args[0]), C.int(len(args)))
	return JSValue(v)
}

func (e *JSState) JSCallGlobal(f JSValue, args []JSValue) JSValue {
	var v C.wkeJSValue
	if len(args) == 0 {
		v = C.wkeJSCallGlobal(e.s, C.wkeJSValue(f), (*C.wkeJSValue)(nil), 0)
	} else {
		v = C.wkeJSCallGlobal(e.s, C.wkeJSValue(f), (*C.wkeJSValue)(&args[0]), C.int(len(args)))
	}
	return JSValue(v)
}

func (e *JSState) JSGetGlobal(prop string) JSValue {
	cs := C.CString(prop)
	v := C.wkeJSGetGlobal(e.s, cs)
	C.free(unsafe.Pointer(cs))
	return JSValue(v)
}

func (e *JSState) JSSetGlobal(prop string, v JSValue) {
	cs := C.CString(prop)
	C.wkeJSSetGlobal(e.s, cs, C.wkeJSValue(v))
	C.free(unsafe.Pointer(cs))
}

// JSGC triggers js garbage collect
func JSCollectGarbge() {
	C.wkeJSCollectGarbge()
}
