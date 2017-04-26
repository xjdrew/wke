package main

import (
	"fmt"

	"github.com/xjdrew/wke"
)

func main() {
	wke.Initialize()
	defer wke.Finalize()

	// print wke version
	fmt.Println(wke.VersionString())

	webView := wke.NewWebView()
	es := webView.GlobalExec()

	// run js code
	v := webView.RunJS("5 + 5")
	fmt.Println("return: ", es.JSToInt(v))

	// go call js function
	webView.RunJS("function jsfunc(a) {return a;}")
	fn := es.JSGetGlobal("jsfunc")
	fmt.Println("IsFunction:", es.JSIsFunction(fn))
	fmt.Println("fn return:", es.JSToInt(es.JSCallGlobal(fn, []wke.JSValue{es.JSInt(10)})))

	// js call go function
	wke.JSBind("hello", func(e *wke.JSState) wke.JSValue {
		v := e.JSToString(e.JSArg(1))
		fmt.Println("hello, ", v)
		return es.JSUndefined()
	})

	wke.JSBind("mysum", func(e *wke.JSState) wke.JSValue {
		count := e.JSArgCount()
		var r int
		for i := 1; i < count; i++ {
			r += e.JSToInt(e.JSArg(i))
		}
		fmt.Println("in mysum:", r)
		return e.JSInt(r)
	})

	wke.JSBind("print", func(e *wke.JSState) wke.JSValue {
		s := e.JSToString(e.JSArg(1))
		fmt.Println("--- jsprint:", s)
		return e.JSUndefined()
	})

	webView.RunJS("gogate('hello', 'world')")
	webView.RunJS("gogate('print',gogate('mysum', 5, 6, 7).toString())")

	// fini
	webView.Destroy()
}
