package negotiator

import (
	"encoding/json"
	"net/http"
	"strings"
)

type JSONProcessor struct {
}

func (*JSONProcessor) CanProcess(mediaRange string) bool {
	return strings.EqualFold(mediaRange, "application/json") ||
		strings.HasPrefix(mediaRange, "application/json-") ||
		strings.EqualFold(mediaRange, "text/json") ||
		strings.HasSuffix(mediaRange, "+json")
}

func (*JSONProcessor) Process(w http.ResponseWriter, model interface{}) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(model)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}
