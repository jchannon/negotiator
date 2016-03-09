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

	xmlProcessor := NewXML()

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

	xmlProcessor := NewXML()

	xmlProcessor.Process(recorder, model)

	assert.Equal(t, "application/xml", recorder.HeaderMap.Get("Content-Type"))
}

func TestShouldSetXmlResponseBody(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := &ValidXMLUser{
		"Joe Bloggs",
	}

	xmlProcessor := NewXML()

	xmlProcessor.Process(recorder, model)

	assert.Equal(t, "<ValidXMLUser><Name>Joe Bloggs</Name></ValidXMLUser>", recorder.Body.String())
}

func TestShouldSetXmlResponseBodyWithIndentation(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := &ValidXMLUser{
		"Joe Bloggs",
	}

	xmlProcessor := NewXMLIndent2Spaces()

	xmlProcessor.Process(recorder, model)

	assert.Equal(t, "<ValidXMLUser>\n  <Name>Joe Bloggs</Name>\n</ValidXMLUser>\n", recorder.Body.String())
}

func TestShouldReturnErrorOnXmlError(t *testing.T) {
	recorder := httptest.NewRecorder()

	model := &XMLUser{
		"Joe Bloggs",
	}

	xmlProcessor := NewXMLIndent2Spaces()

	err := xmlProcessor.Process(recorder, model)

	assert.Error(t, err)
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
