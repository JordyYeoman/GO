package main

import (
	"io"
	"net/http"
)

func main() {
	http.Handle("/files", http.FileServer(http.Dir("./assets")))
	http.HandleFunc("/spicy", spicy)
	http.ListenAndServe("localhost:8080", nil)
}

func spicy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `<img src="shakenbake.png" />`)
}
