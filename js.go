package wke

import (
	"unsafe"
)

// #include <stdlib.h>
// #include "wke.h"
import "C"

type JsType int

const (
	JSTYPE_NUMBER JsType = iota
	JSTYPE_STRING
	JSTYPE_BOOLEAN
	JSTYPE_OBJECT
	JSTYPE_FUNCTION
	JSTYPE_UNDEFINED
)

type JsValue int64

func (v JsValue) TypeOf() JsType {
	return JsType(C.jsTypeOf(C.jsValue(v)))
}

func (v JsValue) IsNumber() bool {
	return GoBool(C.jsIsNumber(C.jsValue(v)))
}
func (v JsValue) IsString() bool {
	return GoBool(C.jsIsString(C.jsValue(v)))
}
func (v JsValue) IsBoolean() bool {
	return GoBool(C.jsIsBoolean(C.jsValue(v)))
}
func (v JsValue) IsObject() bool {
	return GoBool(C.jsIsObject(C.jsValue(v)))
}
func (v JsValue) IsFunction() bool {
	return GoBool(C.jsIsFunction(C.jsValue(v)))
}
func (v JsValue) IsUndefined() bool {
	return GoBool(C.jsIsUndefined(C.jsValue(v)))
}
func (v JsValue) IsNull() bool {
	return GoBool(C.jsIsNull(C.jsValue(v)))
}
func (v JsValue) IsArray() bool {
	return GoBool(C.jsIsArray(C.jsValue(v)))
}
func (v JsValue) IsTrue() bool {
	return GoBool(C.jsIsTrue(C.jsValue(v)))
}
func (v JsValue) IsFalse() bool {
	return GoBool(C.jsIsFalse(C.jsValue(v)))
}

func JsInt(n int) JsValue {
	return JsValue(C.jsInt(C.int(n)))
}

func JsFloat(f float64) JsValue {
	return JsValue(C.jsFloat(C.float(f)))
}
func JsBoolean(b bool) JsValue {
	return JsValue(C.jsBoolean(CBool(b)))
}

func JsUndefined() JsValue {
	return JsValue(C.jsUndefined())
}
func JsNull() JsValue {
	return JsValue(C.jsNull())
}
func JsTrue() JsValue {
	return JsValue(C.jsTrue())
}
func JsFalse() JsValue {
	return JsValue(C.jsFalse())
}

type JsExecState struct {
	s C.jsExecState
}

func (e JsExecState) ToInt(v JsValue) int {
	return int(C.jsToInt(e.s, C.jsValue(v)))
}

func (e JsExecState) ToFloat(v JsValue) float64 {
	return float64(C.jsToFloat(e.s, C.jsValue(v)))
}

func (e JsExecState) ToBoolean(v JsValue) bool {
	return GoBool(C.jsToBoolean(e.s, C.jsValue(v)))
}

func (e JsExecState) ToString(v JsValue) string {
	return C.GoString((*C.char)(C.jsToString(e.s, C.jsValue(v))))
}

func (e JsExecState) JsString(s string) JsValue {
	cs := C.CString(s)
	v := C.jsString(e.s, (*C.utf8)(cs))
	C.free(unsafe.Pointer(cs))
	return JsValue(v)
}

func (e JsExecState) JsObject() JsValue {
	return JsValue(C.jsObject(e.s))
}

func (e JsExecState) JsArray() JsValue {
	return JsValue(C.jsArray(e.s))
}

func (e JsExecState) JsFunction(fn C.jsNativeFunction, argCount uint) JsValue {
	return JsValue(C.jsFunction(e.s, fn, C.uint(argCount)))
}

// JsGlobalObject returns window object
func (e JsExecState) JsGlobalObject() JsValue {
	return JsValue(C.jsGlobalObject(e.s))
}

func (e JsExecState) Eval(s string) JsValue {
	cs := C.CString(s)
	v := C.jsEval(e.s, (*C.utf8)(cs))
	C.free(unsafe.Pointer(cs))
	return JsValue(v)
}

func (e JsExecState) Call(f JsValue, thisObject JsValue, args []JsValue) JsValue {
	v := C.jsCall(e.s, C.jsValue(f), C.jsValue(thisObject), (*C.jsValue)(&args[0]), C.int(len(args)))
	return JsValue(v)
}

func (e JsExecState) CallGlobal(f JsValue, args []JsValue) JsValue {
	v := C.jsCallGlobal(e.s, C.jsValue(f), (*C.jsValue)(&args[0]), C.int(len(args)))
	return JsValue(v)
}

func (e JsExecState) JsGet(object JsValue, prop string) JsValue {
	cs := C.CString(prop)
	v := C.jsGet(e.s, C.jsValue(object), cs)
	C.free(unsafe.Pointer(cs))
	return JsValue(v)
}

func (e JsExecState) JsSet(object JsValue, prop string, v JsValue) {
	cs := C.CString(prop)
	C.jsSet(e.s, C.jsValue(object), cs, C.jsValue(v))
	C.free(unsafe.Pointer(cs))
}

func (e JsExecState) JsGetGlobal(prop string) JsValue {
	cs := C.CString(prop)
	v := C.jsGetGlobal(e.s, cs)
	C.free(unsafe.Pointer(cs))
	return JsValue(v)
}

func (e JsExecState) JsSetGlobal(prop string, v JsValue) {
	cs := C.CString(prop)
	C.jsSetGlobal(e.s, cs, C.jsValue(v))
	C.free(unsafe.Pointer(cs))
}

func (e JsExecState) JsGetAt(object JsValue, index int) JsValue {
	v := C.jsGetAt(e.s, C.jsValue(object), C.int(index))
	return JsValue(v)
}

func (e JsExecState) JsSetAt(object JsValue, index int, v JsValue) {
	C.jsSetAt(e.s, C.jsValue(object), C.int(index), C.jsValue(v))
}

func (e JsExecState) JsGetLength(object JsValue) int {
	return int(C.jsGetLength(e.s, C.jsValue(object)))
}

func (e JsExecState) JsSetLength(object JsValue, length int) {
	C.jsSetLength(e.s, C.jsValue(object), C.int(length))
}

func (e JsExecState) GetWebView() WebView {
	return WebView{
		v: C.jsGetWebView(e.s),
	}
}

// JsGC triggers js garbage collect
func JsGC() {
	C.jsGC()
}
