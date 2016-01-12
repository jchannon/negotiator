package main

import (
	"net/http"

	"github.com/jchannon/negotiator"
	"github.com/jchannon/negotiator/demo/responseprocessors"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func main() {
	goji.Use(negotiatormw)
	goji.Get("/", appHandler(homeHandler))
	goji.Get("/oneoffnegotiator", appHandler(customHandler))
	goji.Get("/multinegotiator", appHandler(multiNegotiatorHandler))
	goji.Get("/multinegotiatoragain", appHandler(multiNegotiatorHandlerAgain))
	goji.Serve()
}

func homeHandler(c web.C, w http.ResponseWriter, req *http.Request) error {
	user := &user{"Joe", "Bloggs"}
	return negotiator.Negotiate(w, req, user)
}

func customHandler(c web.C, w http.ResponseWriter, req *http.Request) error {
	user := &user{"Joe", "Bloggs"}
	//Creating the negotiator could be done for only required handlers or use middleware for all
	textplainNegotiator := negotiator.New(&responseprocessors.PlainTextResponseProcessor{})
	return textplainNegotiator.Negotiate(w, req, user)
}

func multiNegotiatorHandler(c web.C, w http.ResponseWriter, req *http.Request) error {
	user := &user{"Joe", "Bloggs"}
	mynegotiator := c.Env["negotiator"].(*negotiator.Negotiator)
	return mynegotiator.Negotiate(w, req, user)
}

func multiNegotiatorHandlerAgain(c web.C, w http.ResponseWriter, req *http.Request) error {
	user := &user{"John", "Doe"}
	mynegotiator := c.Env["negotiator"].(*negotiator.Negotiator)
	return mynegotiator.Negotiate(w, req, user)
}

type user struct {
	Firstname string
	Lastname  string
}

func negotiatormw(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c.Env["negotiator"] = negotiator.New(&responseprocessors.PlainTextResponseProcessor{})

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Application error handler
// Goji requires to implement both ServeHTTP and ServeHTTPC
type appHandler func(web.C, http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(c web.C, w http.ResponseWriter, r *http.Request) {
	if err := fn(c, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (fn appHandler) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	if err := fn(c, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
