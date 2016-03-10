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

// Negotiator is responsible for content negotiation when using custom response processors.
type Negotiator struct{ processors []ResponseProcessor }

// NewWithJSONAndXML allows users to pass custom response processors. By default, processors
// for XML and JSON are already created.
func NewWithJSONAndXML(responseProcessors ...ResponseProcessor) *Negotiator {
	return New(append(responseProcessors, NewJSON(), NewXML())...)
}

//New allows users to pass custom response processors.
func New(responseProcessors ...ResponseProcessor) *Negotiator {
	return &Negotiator{
		responseProcessors,
	}
}

// Add more response processors. A new Negotiator is returned with the original processors plus
// the extra processors.
func (n *Negotiator) Add(responseProcessors ...ResponseProcessor) *Negotiator {
	return &Negotiator{
		append(n.processors, responseProcessors...),
	}
}

// Negotiate your model based on the HTTP Accept header.
func (n *Negotiator) Negotiate(w http.ResponseWriter, req *http.Request, dataModel interface{}) error {
	return negotiateHeader(n.processors, w, req, dataModel)
}

// Negotiate your model based on the HTTP Accept header. Only XML and JSON are handled.
func Negotiate(w http.ResponseWriter, req *http.Request, model interface{}) error {
	processors := []ResponseProcessor{NewJSON(), NewXML()}
	return negotiateHeader(processors, w, req, model)
}

// A request without any Accept header field implies that the user agent
// will accept any media type in response.
//
// If the header field is present in a request and none of the available
// representations for the response have a media type that is listed as
// acceptable, the origin server can either honour the header field by
// sending a 406 (Not Acceptable) response or disregard the header field
// by treating the response as if it is not subject to content negotiation.
// This implementation prefers the former.
//
// See rfc7231-sec5.3.2:
// http://tools.ietf.org/html/rfc7231#section-5.3.2
func negotiateHeader(processors []ResponseProcessor, w http.ResponseWriter, req *http.Request, dataModel interface{}) error {
	accept := new(accept)

	accept.Header = req.Header.Get("Accept")

	if accept.Header == "" {
		return processors[0].Process(w, dataModel)
	}

	for _, mr := range accept.ParseMediaRanges() {
		if len(mr.Value) == 0 {
			continue
		}

		if strings.EqualFold(mr.Value, "*/*") {
			return processors[0].Process(w, dataModel)
		}

		for _, processor := range processors {
			if processor.CanProcess(mr.Value) {
				return processor.Process(w, dataModel)
			}
		}
	}

	http.Error(w, "", http.StatusNotAcceptable)
	return nil
}
