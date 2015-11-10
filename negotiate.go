package negotiator

import "net/http"

var processors = []ResponseProcessor{&JSONProcessor{}, &XMLProcessor{}}

//New sets up response processors. By default XML and JSON are created
func New(responseProcessors []ResponseProcessor) {
	processors = append(processors, responseProcessors...)
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
