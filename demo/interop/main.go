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
	wke.JSBindFunc("hello", func(e *wke.JSState) wke.JSValue {
		v := e.JSToString(e.JSArg(1))
		fmt.Println("hello, ", v)
		return es.JSUndefined()
	})

	wke.JSBindRaw("testWrap", func(b bool, i int, d float64, s string) string {
		fmt.Println("testWrap:", b, i, d, s)
		return "ok"
	})

	wke.JSBindFunc("mysum", func(e *wke.JSState) wke.JSValue {
		count := e.JSArgCount()
		var r int
		for i := 1; i < count; i++ {
			r += e.JSToInt(e.JSArg(i))
		}
		fmt.Println("in mysum:", r)
		return e.JSInt(r)
	})

	wke.JSBindFunc("print", func(e *wke.JSState) wke.JSValue {
		s := e.JSToString(e.JSArg(1))
		fmt.Println("--- jsprint:", s)
		return e.JSUndefined()
	})

	jst := webView.GlobalExec()
	webView.RunJS("gogate('hello', 'world')")
	webView.RunJS("gogate('print',gogate('mysum', 5, 6, 7).toString())")
	fmt.Println("testWrap:", jst.JSToString(webView.RunJS("gogate('testWrap', true, 10, 3.1415, 'world')")))

	jsv := webView.RunJS("gogate('alert', 'world')")
	fmt.Println(jst.JSIsUndefined(jsv))

	// fini
	webView.Destroy()
}
