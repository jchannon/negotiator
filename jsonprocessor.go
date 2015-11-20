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
		strings.HasSuffix(mediaRange, "+json")
}

func (*jsonProcessor) Process(w http.ResponseWriter, model interface{}, errorHandler func(w http.ResponseWriter, err error)) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(model)

	if err != nil {
		errorHandler(w, err)
		return
	}

	w.Write(js)
}
