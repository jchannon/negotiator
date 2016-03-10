// Package negotiator is a library that handles content negotiation in web applications written in Go.
// Content negotiation is specified by RFC (http://tools.ietf.org/html/rfc7231) and, less formally, by
// Ajax (https://en.wikipedia.org/wiki/XMLHttpRequest).
//
// A Negotiator contains a list of ResponseProcessor. For each call to Negotiate, the best matching
// response processor is chosen and given the task of sending the response.
//
// For more information visit http://github.com/jchannon/negotiator
//
//	func getUser(w http.ResponseWriter, req *http.Request) {
//	    user := &User{"Joe", "Bloggs"}
//	    negotiator.Negotiate(w, req, user)
//	}
//
package negotiator

import (
	"net/http"
	"strings"
)

const (
	xRequestedWith = "X-Requested-With"
	xmlHttpRequest = "XMLHttpRequest"
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

// Firstly, all Ajax requests are processed by the first available Ajax processor.
// Otherwise, standard content negotiation kicks in.
//
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
	if IsAjax(req) {
		for _, processor := range processors {
			ajax, doesAjax := processor.(AjaxResponseProcessor)
			if doesAjax && ajax.IsAjaxResponder() {
				return processor.Process(w, req, dataModel)
			}
		}
	}

	accept := accept(req.Header.Get("Accept"))

	if accept == "" {
		return processors[0].Process(w, req, dataModel)
	}

	for _, mr := range accept.ParseMediaRanges() {
		if len(mr.Value) == 0 {
			continue
		}

		if strings.EqualFold(mr.Value, "*/*") {
			return processors[0].Process(w, req, dataModel)
		}

		for _, processor := range processors {
			if processor.CanProcess(mr.Value) {
				return processor.Process(w, req, dataModel)
			}
		}
	}

	http.Error(w, "", http.StatusNotAcceptable)
	return nil
}

// IsAjax tests whether a request has the Ajax header.
func IsAjax(req *http.Request) bool {
	xRequestedWith, ok := req.Header[xRequestedWith]
	return ok && len(xRequestedWith) == 1 && xRequestedWith[0] == xmlHttpRequest
}
