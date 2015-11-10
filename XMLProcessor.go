package negotiator

import (
	"encoding/xml"
	"net/http"
	"strings"
)

type XMLProcessor struct {
}

func (XMLProcessor) CanProcess(mediaRange string) bool {
	return strings.HasSuffix(mediaRange, "xml")
}

func (XMLProcessor) Process(w http.ResponseWriter, model interface{}) {
	x, err := xml.MarshalIndent(model, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}
