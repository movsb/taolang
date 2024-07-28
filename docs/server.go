//go:build !wasm

package main

import (
	"net/http"
)

const listen = ":3826" // ".tao"

func main() {
	http.Handle("GET /", http.FileServer(http.Dir(`.`)))
	http.ListenAndServe(listen, nil)
}
