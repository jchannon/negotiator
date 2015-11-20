package errorhandlers

import "net/http"

//GlobalErrorHandler is a error handler example to show how to pass a function to negotiator
func GlobalErrorHandler(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
