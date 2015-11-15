package responseprocessors

import (
	"net/http"
	"reflect"
	"strings"
)

type PlainTextResponseProcessor struct {
}

func (*PlainTextResponseProcessor) CanProcess(mediaRange string) bool {
	return strings.EqualFold(mediaRange, "text/plain")
}

func (*PlainTextResponseProcessor) Process(w http.ResponseWriter, model interface{}) {

	w.Header().Set("Content-Type", "text/plain")

	val := reflect.ValueOf(model).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i).String()
		typeField := val.Type().Field(i)

		w.Write([]byte(typeField.Name + " : " + valueField + " "))
	}

}
