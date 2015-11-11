package negotiator

import (
	"encoding/json"
	"net/http"
	"strings"
)

type jsonProcessor struct {
}

func (*jsonProcessor) CanProcess(mediaRange string) bool {
	return strings.EqualFold(mediaRange, "application/json") ||
		strings.HasPrefix(mediaRange, "application/json-") ||
		strings.EqualFold(mediaRange, "text/json") ||
		strings.HasSuffix(mediaRange, "+json")
}

func (*jsonProcessor) Process(w http.ResponseWriter, model interface{}) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(model)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}
