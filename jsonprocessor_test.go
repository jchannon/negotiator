package negotiator

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldProcessJSONAcceptHeader(t *testing.T) {
	var acceptTests = []struct {
		acceptheader string
	}{
		{"application/json"},
		{"application/json-"},
		{"+json"},
	}

	jsonProcessor := &jsonProcessor{}

	for _, tt := range acceptTests {
		result := jsonProcessor.CanProcess(tt.acceptheader)
		assert.True(t, result, "Should process "+tt.acceptheader)
	}
}

func TestShouldSetContentTypeHeader(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := struct {
		Name string
	}{
		"Joe Bloggs",
	}

	jsonProcessor := &jsonProcessor{}

	jsonProcessor.Process(recorder, model)

	assert.Equal(t, "application/json", recorder.HeaderMap.Get("Content-Type"))
}

func TestShouldSetResponseBody(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := struct {
		Name string
	}{
		"Joe Bloggs",
	}

	jsonProcessor := &jsonProcessor{}

	jsonProcessor.Process(recorder, model)

	assert.Equal(t, "{\"Name\":\"Joe Bloggs\"}", recorder.Body.String())
}

func TestShouldReturnErrorOnJsonError(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := &User{
		"Joe Bloggs",
	}

	jsonProcessor := &jsonProcessor{}

	err := jsonProcessor.Process(recorder, model)

	assert.Error(t, err)
}

type User struct {
	Name string
}

func (u *User) MarshalJSON() ([]byte, error) {
	return nil, errors.New("oops")
}

func jsontestErrorHandler(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}
