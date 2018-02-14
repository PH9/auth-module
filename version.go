package main

import (
	"fmt"
	"net/http"
)

func version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprint(w, `{"version":"1.0","build":"1"}`)
}
