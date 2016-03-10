package negotiator

import (
	"encoding"
	"fmt"
	"net/http"
	"strings"
)

const defaultTxtContentType = "text/plain"

type txtProcessor struct {
	contentType string
}

// NewTXT creates an output processor that serialises strings in text/plain form.
// Model values should be one of the following:
//
// * string
//
// * fmt.Stringer
//
// * encoding.TextMarshaler
func NewTXT() ResponseProcessor {
	return &txtProcessor{defaultTxtContentType}
}

// Implements ContentTypeSettable for this type.
func (p *txtProcessor) SetContentType(contentType string) ResponseProcessor {
	p.contentType = contentType
	return p
}

func (*txtProcessor) CanProcess(mediaRange string) bool {
	return strings.EqualFold(mediaRange, "text/plain") || strings.EqualFold(mediaRange, "text/*")
}

func (p *txtProcessor) Process(w http.ResponseWriter, dataModel interface{}) error {
	if dataModel == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	w.Header().Set("Content-Type", p.contentType)
	return p.process(w, dataModel)
}

func (p *txtProcessor) process(w http.ResponseWriter, dataModel interface{}) error {
	s, ok := dataModel.(string)
	if ok {
		writeWithNewline(w, []byte(s))
		return nil
	}

	st, ok := dataModel.(fmt.Stringer)
	if ok {
		writeWithNewline(w, []byte(st.String()))
		return nil
	}

	tm, ok := dataModel.(encoding.TextMarshaler)
	if ok {
		b, err := tm.MarshalText()
		if err != nil {
			return err
		}
		writeWithNewline(w, b)
		return nil
	}

	return fmt.Errorf("Unsupported type for TXT: %T", dataModel)
}
