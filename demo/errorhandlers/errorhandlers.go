package errorhandlers

import "net/http"

func GlobalErrorHandler(w http.ResponseWriter, err error ) {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
}
