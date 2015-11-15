package main

import (
	"net/http"

	"github.com/jchannon/negotiator"
	"github.com/jchannon/negotiator/demo/responseprocessors"
)

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/custom", customHandler)
	http.ListenAndServe(":9001", nil)
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	user := &user{"Joe", "Bloggs"}
	negotiator.Negotiate(w, req, user)
}

func customHandler(w http.ResponseWriter, req *http.Request) {
	user := &user{"Joe", "Bloggs"}
	//Creating the negotiator could be put in middleware so you don't have to do this in every handler
	textplainNegotiator := negotiator.New(&responseprocessors.PlainTextResponseProcessor{})
	textplainNegotiator.Negotiate(w, req, user)
}

type user struct {
	Firstname string
	Lastname  string
}
