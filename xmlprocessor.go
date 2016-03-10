package negotiator

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

const defaultXmlContentType = "application/xml"

type xmlProcessor struct {
	dense          bool
	prefix, indent string
	contentType    string
}

func NewXML() ResponseProcessor {
	return &xmlProcessor{true, "", "", defaultXmlContentType}
}

func NewXMLIndent(prefix, index string) ResponseProcessor {
	return &xmlProcessor{false, prefix, index, defaultXmlContentType}
}

func NewXMLIndent2Spaces() ResponseProcessor {
	return NewXMLIndent("", "  ")
}

func (p *xmlProcessor) SetContentType(contentType string) ResponseProcessor {
	p.contentType = contentType
	return p
}

func (*xmlProcessor) CanProcess(mediaRange string) bool {
	return strings.Contains(mediaRange, "/xml") || strings.HasSuffix(mediaRange, "+xml")
}

func (p *xmlProcessor) Process(w http.ResponseWriter, dataModel interface{}) error {
	if dataModel == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil

	} else {
		w.Header().Set("Content-Type", p.contentType)
		if p.dense {
			return xml.NewEncoder(w).Encode(dataModel)

		} else {
			x, err := xml.MarshalIndent(dataModel, p.prefix, p.indent)
			if err != nil {
				return err
			}

			return writeWithNewline(w, x)
		}
	}
}

func writeWithNewline(w io.Writer, x []byte) error {
	_, err := w.Write(x)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte{'\n'})
	return err
}
