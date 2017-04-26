// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/xjdrew/wke/wnd"
	"github.com/xjdrew/wke/wnd/declarative"
)

func main() {
	var le *walk.LineEdit
	var ww *wnd.WkeWnd

	MainWindow{
		Icon:    Bind("'../img/' + icon(ww.URL) + '.ico'"),
		Title:   "WkeWnd Example",
		MinSize: Size{800, 600},
		Layout:  VBox{MarginsZero: true},
		Children: []Widget{
			LineEdit{
				AssignTo: &le,
				Text:     Bind("ww.URL"),
				OnKeyDown: func(key walk.Key) {
					if key == walk.KeyReturn {
						ww.SetURL(le.Text())
					}
				},
			},
			declarative.WkeWnd{
				AssignTo: &ww,
				Name:     "ww",
				URL:      "https://github.com/lxn/walk",
			},
		},
		Functions: map[string]func(args ...interface{}) (interface{}, error){
			"icon": func(args ...interface{}) (interface{}, error) {
				if strings.HasPrefix(args[0].(string), "https") {
					return "check", nil
				}

				return "stop", nil
			},
		},
	}.Run()
}
