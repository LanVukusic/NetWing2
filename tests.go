package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	fs := http.FileServer(http.Dir("web/view"))
	http.Handle("/s/", http.StripPrefix("/s/", fs))

	http.ListenAndServe(":80", nil)
}
