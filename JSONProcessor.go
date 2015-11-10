package negotiator

import "net/http"

type JSONProcessor struct {
}

func (JSONProcessor) CanProcess(mediaRange string) bool {
	return true
}

func (JSONProcessor) Process(w http.ResponseWriter, model interface{}) {

}
