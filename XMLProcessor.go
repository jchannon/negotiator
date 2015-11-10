package negotiator

import "net/http"

type XMLProcessor struct {
}

func (XMLProcessor) CanProcess(mediaRange string) bool {
	return true
}

func (XMLProcessor) Process(w http.ResponseWriter, model interface{}) {

}
