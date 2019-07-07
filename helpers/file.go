package helpers

import (
	"github.com/zserge/webview"
)

// Alert creates an "alert" windows on the UI
func Alert(w webview.WebView, text string) {
	w.Eval("alert('" + text + "')")
}
