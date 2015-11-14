package main

import (
	"net/http"

	"github.com/jchannon/negotiator"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.ListenAndServe(":9001", nil)
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	user := &user{"Joe", "Bloggs"}
	negotiator.Negotiate(w, req, user)
}

type user struct {
	Firstname string
	Lastname  string
}
