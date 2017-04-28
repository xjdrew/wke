package main

import (
	"log"

	"github.com/xjdrew/wke/iui"
)

func main() {
	wnd := iui.NewWindow("MyMainWindow")
	log.Println("--- create window")
	wnd.ShowWindow(true)
	if err := iui.Run(); err != nil {
		log.Println(err)
	}
}
