package main

import (
	"net/http"

	"github.com/jchannon/negotiator"
	"github.com/jchannon/negotiator/demo/responseprocessors"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/jchannon/negotiator/demo/errorhandlers"
)

func main() {
	goji.Use(negotiatormw)
	goji.Get("/", homeHandler)
	goji.Get("/oneoffnegotiator", customHandler)
	goji.Get("/multinegotiator", multiNegotiatorHandler)
	goji.Get("/multinegotiatoragain", multiNegotiatorHandlerAgain)
	goji.Serve()
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	user := &user{"Joe", "Bloggs"}
	negotiator.Negotiate(w, req, user, errorhandlers.GlobalErrorHandler)
}

func customHandler(w http.ResponseWriter, req *http.Request) {
	user := &user{"Joe", "Bloggs"}
	//Creating the negotiator could be done for only required handlers or use middleware for all
	textplainNegotiator := negotiator.New(&responseprocessors.PlainTextResponseProcessor{})
	textplainNegotiator.ErrorHandler = errorhandlers.GlobalErrorHandler
	textplainNegotiator.Negotiate(w, req, user)
}

func multiNegotiatorHandler(c web.C, w http.ResponseWriter, req *http.Request) {
	user := &user{"Joe", "Bloggs"}
	mynegotiator := c.Env["negotiator"].(*negotiator.Negotiator)
	mynegotiator.Negotiate(w, req, user)
}

func multiNegotiatorHandlerAgain(c web.C, w http.ResponseWriter, req *http.Request) {
	user := &user{"John", "Doe"}
	mynegotiator := c.Env["negotiator"].(*negotiator.Negotiator)
	mynegotiator.Negotiate(w, req, user)
}

type user struct {
	Firstname string
	Lastname  string
}

func negotiatormw(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		n := negotiator.New(&responseprocessors.PlainTextResponseProcessor{})
		n.ErrorHandler = errorhandlers.GlobalErrorHandler
		c.Env["negotiator"] = n

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
