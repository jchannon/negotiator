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

func (*xmlProcessor) Process(w http.ResponseWriter, model interface{}, errorHandler func(w http.ResponseWriter, err error)) {
	x, err := xml.MarshalIndent(model, "", "  ")
	if err != nil {
		errorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}
