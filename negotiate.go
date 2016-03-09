//Package negotiator is a library that handles content negotiation in web applications written in Go.
//
//For more information visit http://github.com/jchannon/negotiator
//
//	func getUser(w http.ResponseWriter, req *http.Request) {
//	    user := &User{"Joe","Bloggs"}
//	    negotiator.Negotiate(w, req, user)
//	}
//
package negotiator

import (
	"net/http"
	"strings"
)

//Negotiator is responsible for content negotiation when using custom response processors.
type Negotiator struct{ processors []ResponseProcessor }

//New allows users to pass custom response processors. By default XML and JSON are already created.
func NewWithJsonAndXml(responseProcessors ...ResponseProcessor) *Negotiator {
	return New(append(responseProcessors, NewJSON(), NewXML())...)
}

//New allows users to pass custom response processors.
func New(responseProcessors ...ResponseProcessor) *Negotiator {
	return &Negotiator{
		responseProcessors,
	}
}

//Negotiate your model based on the HTTP Accept header.
func (n *Negotiator) Negotiate(w http.ResponseWriter, req *http.Request, model interface{}) error {
	return negotiateHeader(n.processors, w, req, model)
}

//Negotiate your model based on the HTTP Accept header. Only XML and JSON are handled.
func Negotiate(w http.ResponseWriter, req *http.Request, model interface{}) error {
	processors := []ResponseProcessor{NewJSON(), NewXML()}
	return negotiateHeader(processors, w, req, model)
}

func negotiateHeader(processors []ResponseProcessor, w http.ResponseWriter, req *http.Request, model interface{}) error {
	accept := new(accept)

	accept.Header = req.Header.Get("Accept")

	// http://tools.ietf.org/html/rfc7231#section-5.3.2
	// rfc7231-sec5.3.2:
	// A request without any Accept header field implies that the user agent
	// will accept any media type in response.
	if accept.Header == "" {
		return processors[0].Process(w, model)
	}

	for _, mr := range accept.ParseMediaRanges() {
		if len(mr.Value) == 0 {
			continue
		}

		if strings.EqualFold(mr.Value, "*/*") {
			return processors[0].Process(w, model)
		}

		for _, processor := range processors {
			if processor.CanProcess(mr.Value) {
				return processor.Process(w, model)
			}
		}
	}

	//rfc2616-sec14.1
	//If an Accept header field is present, and if the
	//server cannot send a response which is acceptable according to the combined
	//Accept field value, then the server SHOULD send a 406 (not acceptable)
	//response.
	http.Error(w, "", http.StatusNotAcceptable)
	return nil
}
