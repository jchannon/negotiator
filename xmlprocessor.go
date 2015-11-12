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

func (*xmlProcessor) Process(w http.ResponseWriter, model interface{}) {
	x, err := xml.MarshalIndent(model, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}
