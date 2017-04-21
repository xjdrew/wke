package main

import (
	"fmt"

	"github.com/xjdrew/wke"
)

func main() {
	// print wke version
	fmt.Println(wke.VersionString())

	webView := wke.NewWebView()
	es := webView.GlobalExec()

	// run js code
	v := webView.RunJS("5 + 5")
	fmt.Println("return: ", es.ToInt(v))

	// go call js function
	webView.RunJS("function jsfunc(a) {return a;}")
	fn := es.JsGetGlobal("jsfunc")
	fmt.Println("IsFunction:", fn.IsFunction())
	fmt.Println("fn return:", es.ToInt(es.CallGlobal(fn, []wke.JsValue{wke.JsInt(10)})))

	// js call go function
	wke.JsBind("hello", func(e wke.JsExecState) wke.JsValue {
		v := e.ToString(e.JsArg(1))
		fmt.Println("hello, ", v)
		return wke.JsUndefined()
	})

	wke.JsBind("mysum", func(e wke.JsExecState) wke.JsValue {
		count := e.JsArgCount()
		var r int
		for i := 1; i < count; i++ {
			r += e.ToInt(e.JsArg(i))
		}
		fmt.Println("in mysum:", r)
		return wke.JsInt(r)
	})

	wke.JsBind("print", func(e wke.JsExecState) wke.JsValue {
		s := e.ToString(e.JsArg(1))
		fmt.Println("--- jsprint:", s)
		return wke.JsUndefined()
	})

	webView.RunJS("gogate('hello', 'world')")
	webView.RunJS("gogate('print',gogate('mysum', 5, 6, 7).toString())")

	// fini
	webView.Destroy()
	wke.Shutdown()
}
