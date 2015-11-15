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

import "net/http"

//Negotiator is responsible for content negotiation when using custom response processors
type Negotiator struct{ processors []ResponseProcessor }

//New allows users to pass custom response processors. By default XML and JSON are already created
func New(responseProcessors ...ResponseProcessor) *Negotiator {
	processors := []ResponseProcessor{&jsonProcessor{}, &xmlProcessor{}}
	processors = append(responseProcessors, processors...)
	return &Negotiator{
		processors,
	}
}

//Negotiate your model based on HTTP Accept header
func (n *Negotiator) Negotiate(w http.ResponseWriter, req *http.Request, model interface{}) {
	negotiateHeader(n.processors, w, req, model)
}

//Negotiate your model based on HTTP Accept header. By default XML and JSON are handled
func Negotiate(w http.ResponseWriter, req *http.Request, model interface{}) {
	processors := []ResponseProcessor{&jsonProcessor{}, &xmlProcessor{}}
	negotiateHeader(processors, w, req, model)
}

func negotiateHeader(processors []ResponseProcessor, w http.ResponseWriter, req *http.Request, model interface{}) {
	accept := new(accept)

	accept.Header = req.Header.Get("Accept")

	for _, mr := range accept.ParseMediaRanges() {
		if len(mr.Value) == 0 {
			continue
		}
		for _, processor := range processors {
			if processor.CanProcess(mr.Value) {
				processor.Process(w, model)
				return
			}
		}
	}

	//rfc2616-sec14.1
	//If an Accept header field is present, and if the
	//server cannot send a response which is acceptable according to the combined
	//Accept field value, then the server SHOULD send a 406 (not acceptable)
	//response.
	http.Error(w, "", http.StatusNotAcceptable)
}
