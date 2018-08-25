package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "HTTP Request: %s %s\n", r.Method, r.URL.Path)
	})
	http.ListenAndServe(":80", nil)
}
