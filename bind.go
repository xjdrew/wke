package wke

import (
	"sync"
)

/*
#include "wke.h"

extern jsValue gogate(jsExecState es);
*/
import "C"

type JsBindFunc func(e JsExecState) JsValue

var mu sync.Mutex
var jsBindingFunctions map[string]JsBindFunc

//export jsNativeCall
func jsNativeCall(name *C.char, e C.jsExecState) C.jsValue {
	s := C.GoString(name)
	fn := jsBindingFunctions[s]
	return C.jsValue(fn(JsExecState{e}))
}

func init() {
	jsBindingFunctions = make(map[string]JsBindFunc)
	C.jsBindFunction(C.CString("gogate"), C.jsNativeFunction(C.gogate), 0)
}

func JsBind(name string, fn JsBindFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := jsBindingFunctions[name]; ok {
		panic("repeated bind: " + name)
	}
	jsBindingFunctions[name] = fn
}

func JsUnbind(name string, fn JsBindFunc) {
	mu.Lock()
	defer mu.Unlock()
	delete(jsBindingFunctions, name)
}
