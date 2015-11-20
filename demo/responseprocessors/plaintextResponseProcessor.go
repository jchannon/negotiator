package responseprocessors

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type PlainTextResponseProcessor struct {
}

func (*PlainTextResponseProcessor) CanProcess(mediaRange string) bool {
	return strings.EqualFold(mediaRange, "text/plain")
}

func (*PlainTextResponseProcessor) Process(w http.ResponseWriter, model interface{}, errorHandler func(w http.ResponseWriter, err error)) {

	if currTime := time.Now(); currTime.Second()%2 == 0 {
		err := errors.New("This is a sample error showcasing how to use a error handler with negotiator")
		errorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	val := reflect.ValueOf(model).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i).String()
		typeField := val.Type().Field(i)

		w.Write([]byte(typeField.Name + " : " + valueField + " "))
	}

}
