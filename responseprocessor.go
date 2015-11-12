package negotiator

import "net/http"

//ResponseProcessor interface that creates contract for custom content negotiation
type ResponseProcessor interface {
	CanProcess(mediaRange string) bool
	Process(w http.ResponseWriter, model interface{})
}
