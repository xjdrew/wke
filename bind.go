package wke

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

/*
#include <stdlib.h>
#include "wke.h"

extern wkeJSValue  gogate(wkeJSState* es);
*/
import "C"

type JSHandler interface {
	Handle(*JSState) JSValue
}

type JSHandlerFunc func(e *JSState) JSValue

// Handle calls f(e)
func (f JSHandlerFunc) Handle(e *JSState) JSValue {
	return f(e)
}

type JSRawHandler struct {
	fn reflect.Value
}

func (h *JSRawHandler) Handle(e *JSState) JSValue {
	typ := h.fn.Type()
	numIn := typ.NumIn()

	funcName := e.JSToString(e.JSArg(0))

	// arg0 is funcation name
	if typ.NumIn() != e.JSArgCount()-1 {
		fmt.Printf("gogate: call <%s> with unmatch input arguments: %d != %d", funcName, numIn, e.JSArgCount()-1)
		return e.JSUndefined()
	}

	var in []reflect.Value
	for i := 0; i < typ.NumIn(); i++ {
		kind := typ.In(i).Kind()
		v := e.JSArg(i + 1)
		switch kind {
		case reflect.Bool:
			if e.JSIsBool(v) {
				in = append(in, reflect.ValueOf(e.JSToBoolean(v)))
				continue
			}
		case reflect.Int:
			if e.JSIsNumber(v) {
				in = append(in, reflect.ValueOf(e.JSToInt(v)))
				continue
			}
		case reflect.Float64:
			if e.JSIsNumber(v) {
				in = append(in, reflect.ValueOf(e.JSToDouble(v)))
				continue
			}
		case reflect.String:
			if e.JSIsString(v) {
				in = append(in, reflect.ValueOf(e.JSToString(v)))
				continue
			}
		}
		fmt.Printf("gogate: funcation <%s> argument %d should be %s", funcName, i, kind)
		return e.JSUndefined()
	}

	out := h.fn.Call(in)
	if len(out) == 0 {
		return e.JSUndefined()
	} else {
		ov := out[0]
		switch ov.Type().Kind() {
		case reflect.Bool:
			return e.JSBool(ov.Bool())
		case reflect.Int:
			return e.JSInt(int(ov.Int()))
		case reflect.Float64:
			return e.JSDouble(ov.Float())
		case reflect.String:
			return e.JSString(ov.String())
		default:
			fmt.Printf("gogate: funcation <%s> output type<%s> is unsupported", funcName, ov.Type().Kind())
			return e.JSUndefined()
		}
	}
}

// panic if failed
func verifyJSRawHandler(fn interface{}) {
	typ := reflect.TypeOf(fn)
	if typ.Kind() != reflect.Func {
		panic("WrapJSFunc: fn is not a funcation")
	}

	if typ.NumOut() > 1 {
		panic("WrapJSFunc: fn's output parameter count is more than 1")
	}

	for i := 0; i < typ.NumOut(); i++ {
		out := typ.Out(i)
		kind := out.Kind()
		if kind != reflect.Bool && kind != reflect.Int && kind != reflect.Float64 && kind != reflect.String {
			panic(fmt.Sprintf("WrapJSFunc: output parameter %d's type<%s> is unsupported", i, kind))
		}
	}

	for i := 0; i < typ.NumIn(); i++ {
		in := typ.In(i)
		kind := in.Kind()
		if kind != reflect.Bool && kind != reflect.Int && kind != reflect.Float64 && kind != reflect.String {
			panic(fmt.Sprintf("WrapJSFunc: argument %d's type<%s> is unsupported", i, kind))
		}
	}
}

// fn rules:
// * must be funcation
// * 0 or 1 output parameter
// * input and output parameter type must be in [bool, int, float64, string]
// panic if violate rules
func NewJSRawHandler(fn interface{}) *JSRawHandler {
	verifyJSRawHandler(fn)
	return &JSRawHandler{
		fn: reflect.ValueOf(fn),
	}
}

var nativeFunctions struct {
	sync.Mutex
	handlers map[string]JSHandler
}

func jsbind(name string, handler JSHandler) {
	nativeFunctions.Lock()
	defer nativeFunctions.Unlock()
	if _, ok := nativeFunctions.handlers[name]; ok {
		panic("repeated bind: " + name)
	}
	nativeFunctions.handlers[name] = handler
}

func JSBindFunc(name string, fn func(e *JSState) JSValue) {
	jsbind(name, JSHandlerFunc(fn))
}

func JSBindRaw(name string, fn interface{}) {
	jsbind(name, NewJSRawHandler(fn))
}

func JSBind(name string, handler JSHandler) {
	jsbind(name, handler)
}

func JSUnbindFunc(name string) {
	nativeFunctions.Lock()
	defer nativeFunctions.Unlock()
	delete(nativeFunctions.handlers, name)
}

//export goNativeCall
func goNativeCall(name *C.char, e *C.wkeJSState) C.wkeJSValue {
	s := C.GoString(name)

	nativeFunctions.Lock()
	handler := nativeFunctions.handlers[s]
	nativeFunctions.Unlock()

	state := &JSState{e}
	if handler == nil {
		fmt.Printf("goNativeCall: undefined function %s", s)
		return C.wkeJSValue(state.JSUndefined())
	}
	return C.wkeJSValue(handler.Handle(state))
}

func init() {
	nativeFunctions.handlers = make(map[string]JSHandler)

	name := C.CString("gogate")
	defer C.free(unsafe.Pointer(name))
	C.wkeJSBindFunction(name, C.wkeJSNativeFunction(C.gogate), 1)
}
