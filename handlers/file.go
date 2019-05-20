package handlers

import (
	"github.com/zserge/webview"
)

// W is a global variable used to get webview components
var W webview.WebView

// Must is an error handler
func Must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
