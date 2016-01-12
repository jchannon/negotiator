package negotiator

import (
	"encoding/xml"
	"net/http"
	"strings"
)

type xmlProcessor struct {
}

func (*xmlProcessor) CanProcess(mediaRange string) bool {
	return strings.HasSuffix(mediaRange, "xml")
}

func (*xmlProcessor) Process(w http.ResponseWriter, model interface{}) error {
	x, err := xml.MarshalIndent(model, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/xml")
	_, err = w.Write(x)
	return err
}
