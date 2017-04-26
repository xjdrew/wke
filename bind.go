package wke

import (
	"sync"
	"unsafe"
)

/*
#include <stdlib.h>
#include "wke.h"

extern wkeJSValue  gogate(wkeJSState* es);
*/
import "C"

type JSBindFunc func(e *JSState) JSValue

var mu sync.Mutex
var jsBindingFunctions map[string]JSBindFunc

//export jsNativeCall
func jsNativeCall(name *C.char, e *C.wkeJSState) C.wkeJSValue {
	s := C.GoString(name)
	fn := jsBindingFunctions[s]
	return C.wkeJSValue(fn(&JSState{e}))
}

func init() {
	jsBindingFunctions = make(map[string]JSBindFunc)

	name := C.CString("gogate")
	defer C.free(unsafe.Pointer(name))
	C.wkeJSBindFunction(name, C.wkeJSNativeFunction(C.gogate), 1)
}

func JSBind(name string, fn JSBindFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := jsBindingFunctions[name]; ok {
		panic("repeated bind: " + name)
	}
	jsBindingFunctions[name] = fn
}

func JSUnbind(name string, fn JSBindFunc) {
	mu.Lock()
	defer mu.Unlock()
	delete(jsBindingFunctions, name)
}
