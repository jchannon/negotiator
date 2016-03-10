package negotiator

import "net/http"

// ResponseProcessor interface that creates contract for custom content negotiation.
type ResponseProcessor interface {
	SetContentType(contentType string) ResponseProcessor
	CanProcess(mediaRange string) bool
	Process(w http.ResponseWriter, dataModel interface{}) error
}
