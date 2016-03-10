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

func (p *jsonProcessor) Process(w http.ResponseWriter, dataModel interface{}) error {
	if dataModel == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil

	} else {
		w.Header().Set("Content-Type", "application/json")
		if p.dense {
			return json.NewEncoder(w).Encode(dataModel)

		} else {
			js, err := json.MarshalIndent(dataModel, p.prefix, p.indent)

			if err != nil {
				return err
			}

			return writeWithNewline(w, js)
		}
	}
}
