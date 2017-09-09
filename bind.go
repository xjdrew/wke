package wke

import (
	"sync"
	"unsafe"

	"bitbucket.org/istoneio/golib/islog"
)

/*
#include <stdlib.h>
#include "wke.h"

extern wkeJSValue  gogate(wkeJSState* es);
*/
import "C"

type JSBindFunc func(e *JSState) JSValue

var nativeFunctions struct {
	sync.Mutex
	fs map[string]JSBindFunc
}

func JSBind(name string, fn JSBindFunc) {
	nativeFunctions.Lock()
	defer nativeFunctions.Unlock()
	if _, ok := nativeFunctions.fs[name]; ok {
		panic("repeated bind: " + name)
	}
	nativeFunctions.fs[name] = fn
}

func JSUnbind(name string, fn JSBindFunc) {
	nativeFunctions.Lock()
	defer nativeFunctions.Unlock()
	delete(nativeFunctions.fs, name)
}

//export goNativeCall
func goNativeCall(name *C.char, e *C.wkeJSState) C.wkeJSValue {
	s := C.GoString(name)
	fn := nativeFunctions.fs[s]
	state := &JSState{e}
	if fn == nil {
		islog.Errorf("goNativeCall: undefined function %s", s)
		return C.wkeJSValue(state.JSUndefined())
	}
	return C.wkeJSValue(fn(state))
}

func init() {
	nativeFunctions.fs = make(map[string]JSBindFunc)

	name := C.CString("gogate")
	defer C.free(unsafe.Pointer(name))
	C.wkeJSBindFunction(name, C.wkeJSNativeFunction(C.gogate), 1)
}
