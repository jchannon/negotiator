# Negotiator [![GoDoc](https://godoc.org/github.com/jchannon/negotiator?status.png)](http://godoc.org/github.com/jchannon/negotiator) [![Build Status](https://travis-ci.org/jchannon/negotiator.svg?branch=master)](https://travis-ci.org/jchannon/negotiator)
This is a libary that handles content negotiation in web applications written in Go.

## Usage

### Simple
To return JSON/XML out of the box simple put this in your route handler:
```
func getUser(w http.ResponseWriter, req *http.Request) {
    user := &User{"Joe","Bloggs"}
    if err := negotiator.Negotiate(w, req, user); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```
### Custom

To add your own negotiator, for example you want to write a PDF with your model, do the following:


1) Create a type that conforms to the [ResponseProcessor](https://github.com/jchannon/negotiator/blob/master/responseprocessor.go) interface

2) Call `negotiator.New(responseProcessors ...ResponseProcessor)` and pass in a your custom processor. When your request handler calls `negotiator.Negotiate(w,req,model,errorHandler)` it will render a PDF if your Accept header defined it wanted a PDF response.
