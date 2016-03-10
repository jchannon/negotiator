package negotiator

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldAddCustomResponseProcessors(t *testing.T) {

	var fakeResponseProcessor = &fakeProcessor{}
	negotiator := NewWithJSONAndXML(fakeResponseProcessor)

	assert.Equal(t, 3, len(negotiator.processors))
}

func TestShouldAddCustomResponseProcessors2(t *testing.T) {

	var fakeResponseProcessor = &fakeProcessor{}
	negotiator := New().Add(NewJSON(), NewXML()).Add(fakeResponseProcessor)

	assert.Equal(t, 3, len(negotiator.processors))
}

func TestShouldAddCustomResponseProcessorsToBeginning(t *testing.T) {

	var fakeResponseProcessor = &fakeProcessor{}
	negotiator := NewWithJSONAndXML(fakeResponseProcessor)

	firstProcessor := negotiator.processors[0]
	processorName := reflect.TypeOf(firstProcessor).String()

	assert.Equal(t, "*negotiator.fakeProcessor", processorName)
}

func TestShouldReturnDefaultProcessorIfNoAcceptHeader(t *testing.T) {
	var fakeResponseProcessor = &fakeProcessor{}
	negotiator := New(fakeResponseProcessor)

	req, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	negotiator.Negotiate(recorder, req, nil)

	assert.Equal(t, "boo ya!", recorder.Body.String())
}

func TestShouldGiveJSONResponseForAjaxRequests(t *testing.T) {
	negotiator := NewWithJSONAndXML()

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add(xRequestedWith, xmlHttpRequest)
	recorder := httptest.NewRecorder()

	model := &ValidXMLUser{
		"Joe Bloggs",
	}
	negotiator.Negotiate(recorder, req, model)

	assert.Equal(t, "{\"Name\":\"Joe Bloggs\"}\n", recorder.Body.String())
}

func TestShouldReturn406IfNoMatchingAcceptHeader(t *testing.T) {
	var fakeResponseProcessor = &fakeProcessor{}
	negotiator := New(fakeResponseProcessor)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept", "image/png")
	recorder := httptest.NewRecorder()

	negotiator.Negotiate(recorder, req, nil)

	assert.Equal(t, http.StatusNotAcceptable, recorder.Code)
}

func TestShouldReturn406IfNoAcceptHeaderWithoutCustomResponseProcessor(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	Negotiate(recorder, req, "foo")

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestShouldNegotiateAndWriteToResponseBody(t *testing.T) {
	var fakeResponseProcessor = &fakeProcessor{}
	negotiator := New(fakeResponseProcessor)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept", "application/negotiatortesting")
	recorder := httptest.NewRecorder()

	negotiator.Negotiate(recorder, req, nil)

	assert.Equal(t, "boo ya!", recorder.Body.String())
}

func TestShouldNegotiateADefaultProcessor(t *testing.T) {
	var fakeResponseProcessor = &fakeProcessor{}
	negotiator := New(fakeResponseProcessor)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept", "*/*")
	recorder := httptest.NewRecorder()

	negotiator.Negotiate(recorder, req, nil)

	assert.Equal(t, "boo ya!", recorder.Body.String())
}

type fakeProcessor struct{}

func (*fakeProcessor) CanProcess(mediaRange string) bool {
	return strings.EqualFold(mediaRange, "application/negotiatortesting")
}

func (*fakeProcessor) Process(w http.ResponseWriter, model interface{}) error {
	w.Write([]byte("boo ya!"))
	return nil
}
