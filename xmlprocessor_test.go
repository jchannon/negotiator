package negotiator

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldProcessXMLAcceptHeader(t *testing.T) {
	var acceptTests = []struct {
		acceptheader string
	}{
		{"application/xml"},
	}

	xmlProcessor := &xmlProcessor{}

	for _, tt := range acceptTests {
		result := xmlProcessor.CanProcess(tt.acceptheader)
		assert.True(t, result, "Should process "+tt.acceptheader)
	}
}

func TestShouldSetXmlContentTypeHeader(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := &ValidXMLUser{
		"Joe Bloggs",
	}

	xmlProcessor := &xmlProcessor{}

	xmlProcessor.Process(recorder, model, nil)

	assert.Equal(t, "application/xml", recorder.HeaderMap.Get("Content-Type"))
}

func TestShouldSetXmlResponseBody(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := &ValidXMLUser{
		"Joe Bloggs",
	}

	xmlProcessor := &xmlProcessor{}

	xmlProcessor.Process(recorder, model, nil)

	assert.Equal(t, "<ValidXMLUser>\n  <Name>Joe Bloggs</Name>\n</ValidXMLUser>", recorder.Body.String())
}

func TestShouldCallErrorHandlerOnXmlError(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := &XMLUser{
		"Joe Bloggs",
	}

	xmlProcessor := &xmlProcessor{}

	xmlProcessor.Process(recorder, model, xmltestErrorHandler)

	assert.Equal(t, 500, recorder.Code)
}

type ValidXMLUser struct {
	Name string
}

type XMLUser struct {
	Name string
}

func (u *XMLUser) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return errors.New("oops")
}

func xmltestErrorHandler(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}
