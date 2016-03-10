package negotiator

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONShouldProcessAcceptHeader(t *testing.T) {
	var acceptTests = []struct {
		acceptheader string
		expected     bool
	}{
		{"application/json", true},
		{"application/json-", true},
		{"application/CEA", false},
		{"+json", true},
	}

	processor := NewJSON()

	for _, tt := range acceptTests {
		result := processor.CanProcess(tt.acceptheader)
		assert.Equal(t, tt.expected, result, "Should process "+tt.acceptheader)
	}
}

func TestJSONShouldReturnNoContentIfNil(t *testing.T) {
	recorder := httptest.NewRecorder()

	processor := NewJSON()

	processor.Process(recorder, nil)

	assert.Equal(t, 204, recorder.Code)
}

func TestJSONShouldSetDefaultContentTypeHeader(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := struct {
		Name string
	}{
		"Joe Bloggs",
	}

	processor := NewJSON()

	processor.Process(recorder, model)

	assert.Equal(t, "application/json", recorder.HeaderMap.Get("Content-Type"))
}

func TestJSONShouldSetContentTypeHeader(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := struct {
		Name string
	}{
		"Joe Bloggs",
	}

	processor := NewJSON().SetContentType("application/calendar+json")

	processor.Process(recorder, model)

	assert.Equal(t, "application/calendar+json", recorder.HeaderMap.Get("Content-Type"))
}

func TestJSONShouldSetResponseBody(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := struct {
		Name string
	}{
		"Joe Bloggs",
	}

	processor := NewJSON()

	processor.Process(recorder, model)

	assert.Equal(t, "{\"Name\":\"Joe Bloggs\"}\n", recorder.Body.String())
}

func TestJSONShouldSetResponseBodyWithIndentation(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := struct {
		Name string
	}{
		"Joe Bloggs",
	}

	processor := NewJSONIndent2Spaces()

	processor.Process(recorder, model)

	assert.Equal(t, "{\n  \"Name\": \"Joe Bloggs\"\n}\n", recorder.Body.String())
}

func TestJSONShouldReturnErrorOnError(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := &User{
		"Joe Bloggs",
	}

	processor := NewJSON()

	err := processor.Process(recorder, model)

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
