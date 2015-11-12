//Package negotiator is a libary that handles content negotiation in web applications written in Go.
//
//For more infomation visit http://github.com/jchannon/negotiator
//
//	func getUser(w http.ResponseWriter, req *http.Request) {
//	    user := &User{"Joe","Bloggs"}
//	    negotiator.Negotiate(w, req, user)
//	}
//
package negotiator

import "net/http"

var processors = []ResponseProcessor{&jsonProcessor{}, &xmlProcessor{}}

//New sets up response processors. By default XML and JSON are created
func New(responseProcessors ...*ResponseProcessor) {
	for _, proc := range responseProcessors {
		processors = append(processors, *proc)
	}
}

//Negotiate your model based on HTTP Accept header
func Negotiate(w http.ResponseWriter, req *http.Request, model interface{}) {
	for _, processor := range processors {
		acceptHeader := req.Header.Get("Accept")
		if processor.CanProcess(acceptHeader) {
			processor.Process(w, model)
			return
		}
	}
}
