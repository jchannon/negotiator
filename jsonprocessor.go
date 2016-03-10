package negotiator

import (
	"encoding/json"
	"net/http"
	"strings"
)

const defaultJSONContentType = "application/json"

type jsonProcessor struct {
	dense          bool
	prefix, indent string
	contentType    string
}

// NewJSON creates a new processor for XML withOUT indentation.
func NewJSON() ResponseProcessor {
	return &jsonProcessor{true, "", "", defaultJSONContentType}
}

// NewJSONIndent creates a new processor for XML with specified indentation.
func NewJSONIndent(prefix, index string) ResponseProcessor {
	return &jsonProcessor{false, prefix, index, defaultJSONContentType}
}

// NewJSONIndent2Spaces creates a new processor for XML with 2-space indentation.
func NewJSONIndent2Spaces() ResponseProcessor {
	return NewJSONIndent("", "  ")
}

func (p *jsonProcessor) SetContentType(contentType string) ResponseProcessor {
	p.contentType = contentType
	return p
}

func (*jsonProcessor) CanProcess(mediaRange string) bool {
	return strings.EqualFold(mediaRange, "application/json") ||
		strings.HasPrefix(mediaRange, "application/json-") ||
		strings.HasSuffix(mediaRange, "+json")
}

func (p *jsonProcessor) Process(w http.ResponseWriter, dataModel interface{}) error {
	if dataModel == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	w.Header().Set("Content-Type", p.contentType)
	if p.dense {
		return json.NewEncoder(w).Encode(dataModel)
	}

	js, err := json.MarshalIndent(dataModel, p.prefix, p.indent)

	if err != nil {
		return err
	}

	return writeWithNewline(w, js)
}
