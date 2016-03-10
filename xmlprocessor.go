package negotiator

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

type xmlProcessor struct {
	dense          bool
	prefix, indent string
}

func NewXML() ResponseProcessor {
	return &xmlProcessor{true, "", ""}
}

func NewXMLIndent(prefix, index string) ResponseProcessor {
	return &xmlProcessor{false, prefix, index}
}

func NewXMLIndent2Spaces() ResponseProcessor {
	return NewXMLIndent("", "  ")
}

func (*xmlProcessor) CanProcess(mediaRange string) bool {
	return strings.HasSuffix(mediaRange, "xml")
}

func (p *xmlProcessor) Process(w http.ResponseWriter, dataModel interface{}) error {
	if dataModel == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil

	} else {
		w.Header().Set("Content-Type", "application/xml")
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
