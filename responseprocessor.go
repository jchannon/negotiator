package negotiator

import "net/http"

// ResponseProcessor interface creates the contract for custom content negotiation.
type ResponseProcessor interface {
	CanProcess(mediaRange string) bool
	Process(w http.ResponseWriter, dataModel interface{}) error
}

// ContentTypeSettable interface provides for those response processors that allow the
// response Content-Type to be set explicitly.
type ContentTypeSettable interface {
	SetContentType(contentType string) ResponseProcessor
}

// AjaxResponseProcessor interface allows content negotiation to be biased when
// Ajax requests are handled.
type AjaxResponseProcessor interface {
	IsAjaxResponder() bool
}
