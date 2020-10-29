package web

import (
	webview2 "github.com/webview/webview"
)

func ListtyWeb() {

	wv := webview2.New(true)
	defer wv.Destroy()
	wv.SetTitle("My test")
	wv.SetSize(800, 600, webview2.HintNone)
	wv.Navigate("http://stelios.dev")
	wv.Run()

}
