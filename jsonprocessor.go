package negotiator

import (
	"encoding/json"
	"net/http"
	"strings"
)

type jsonProcessor struct {
	dense          bool
	prefix, indent string
}

func NewJSON() ResponseProcessor {
	return &jsonProcessor{true, "", ""}
}

func NewJSONIndent(prefix, index string) ResponseProcessor {
	return &jsonProcessor{false, prefix, index}
}

func NewJSONIndent2Spaces() ResponseProcessor {
	return NewJSONIndent("", "  ")
}

func (*jsonProcessor) CanProcess(mediaRange string) bool {
	return strings.EqualFold(mediaRange, "application/json") ||
		strings.HasPrefix(mediaRange, "application/json-") ||
		strings.HasSuffix(mediaRange, "+json")
}

func (p *jsonProcessor) Process(w http.ResponseWriter, model interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if p.dense {
		return json.NewEncoder(w).Encode(model)

	} else {
		js, err := json.MarshalIndent(model, p.prefix, p.indent)

		if err != nil {
			return err
		}

		_, err = w.Write(js)
		_, err = w.Write([]byte{'\n'})
		return err
	}
}
