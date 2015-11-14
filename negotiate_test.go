package negotiator

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldAddCustomResponseProcessors(t *testing.T) {

	var fakeResponseProcessor = &fakeProcessor{}
	negotiator := New(fakeResponseProcessor)

	assert.Equal(t, 3, len(negotiator.processors))
}

func TestShouldAddCustomResponseProcessorsToBeginning(t *testing.T) {

	var fakeResponseProcessor = &fakeProcessor{}
	negotiator := New(fakeResponseProcessor)

	firstProcessor := negotiator.processors[0]
	processorName := reflect.TypeOf(firstProcessor).String()

	assert.Equal(t, "*negotiator.fakeProcessor", processorName)
}

func TestShouldReturn406IfNoAcceptHeader(t *testing.T) {
	var fakeResponseProcessor = &fakeProcessor{}
	negotiator := New(fakeResponseProcessor)

	req, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	negotiator.Negotiate(recorder, req, nil)

	assert.Equal(t, 406, recorder.Code)
}

func TestShouldReturn406IfNoAcceptHeaderWithoutCustomerResponseProcessor(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	Negotiate(recorder, req, nil)

	assert.Equal(t, 406, recorder.Code)

}

type fakeProcessor struct {
}

func (*fakeProcessor) CanProcess(mediaRange string) bool {
	return true
}

func (*fakeProcessor) Process(w http.ResponseWriter, model interface{}) {

}
