package negotiator

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTXTShouldProcessAcceptHeader(t *testing.T) {
	var acceptTests = []struct {
		acceptheader string
		expected     bool
	}{
		{"text/plain", true},
		{"text/*", true},
		{"text/csv", false},
	}

	processor := NewTXT()

	for _, tt := range acceptTests {
		result := processor.CanProcess(tt.acceptheader)
		assert.Equal(t, tt.expected, result, "Should process "+tt.acceptheader)
	}
}

func TestTXTShouldReturnNoContentIfNil(t *testing.T) {
	recorder := httptest.NewRecorder()

	processor := NewTXT()

	processor.Process(recorder, nil)

	assert.Equal(t, 204, recorder.Code)
}

func TestTXTShouldSetDefaultContentTypeHeader(t *testing.T) {
	recorder := httptest.NewRecorder()

	processor := NewTXT()

	processor.Process(recorder, "Joe Bloggs")

	assert.Equal(t, "text/plain", recorder.HeaderMap.Get("Content-Type"))
}

func TestTXTShouldSetContentTypeHeader(t *testing.T) {
	recorder := httptest.NewRecorder()

	processor := NewTXT().(ContentTypeSettable).SetContentType("text/rtf")

	processor.Process(recorder, "Joe Bloggs")

	assert.Equal(t, "text/rtf", recorder.HeaderMap.Get("Content-Type"))
}

func TestTXTShouldSetResponseBody(t *testing.T) {
	models := []struct {
		stuff    interface{}
		expected string
	}{
		{"Joe Bloggs", "Joe Bloggs\n"},
		{hidden{tt(2001, 10, 31)}, "(2001-10-31)\n"},
	}

	processor := NewTXT()

	for _, m := range models {
		recorder := httptest.NewRecorder()
		err := processor.Process(recorder, m.stuff)
		assert.NoError(t, err)
		assert.Equal(t, m.expected, recorder.Body.String())
	}
}

func TestTXTShouldReturnErrorOnError(t *testing.T) {
	recorder := httptest.NewRecorder()

	processor := NewTXT()

	err := processor.Process(recorder, make(chan int, 0))

	assert.Error(t, err)
}
