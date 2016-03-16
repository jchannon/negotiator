package negotiator

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

const defaultXMLContentType = "application/xml"

type xmlProcessor struct {
	dense          bool
	prefix, indent string
	contentType    string
}

// NewXML creates a new processor for XML without indentation.
func NewXML() ResponseProcessor {
	return &xmlProcessor{true, "", "", defaultXMLContentType}
}

// NewXMLIndent creates a new processor for XML with a specified indentation.
func NewXMLIndent(prefix, index string) ResponseProcessor {
	return &xmlProcessor{false, prefix, index, defaultXMLContentType}
}

// NewXMLIndent2Spaces creates a new processor for XML with 2-space indentation.
func NewXMLIndent2Spaces() ResponseProcessor {
	return NewXMLIndent("", "  ")
}

// Implements ContentTypeSettable for this type.
func (p *xmlProcessor) SetContentType(contentType string) ResponseProcessor {
	p.contentType = contentType
	return p
}

func (*xmlProcessor) CanProcess(mediaRange string) bool {
	return strings.Contains(mediaRange, "/xml") || strings.HasSuffix(mediaRange, "+xml")
}

func (p *xmlProcessor) Process(w http.ResponseWriter, req *http.Request, dataModel interface{}, context ...interface{}) error {
	if dataModel == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	w.Header().Set("Content-Type", p.contentType)
	if p.dense {
		return xml.NewEncoder(w).Encode(dataModel)
	}

	x, err := xml.MarshalIndent(dataModel, p.prefix, p.indent)
	if err != nil {
		return err
	}

	return writeWithNewline(w, x)
}

func writeWithNewline(w io.Writer, x []byte) error {
	_, err := w.Write(x)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte{'\n'})
	return err
}
