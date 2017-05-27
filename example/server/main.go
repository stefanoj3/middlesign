package main

import (
	"github.com/justinas/alice"
	"github.com/stefanoj3/middlesign"
	"net/http"
)

func myErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
}

func myApp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func main() {
	middlesignConfig := middlesign.DefaultConfig("my_super_secret")

	signedRequestMiddlewhare := func(h http.Handler) http.Handler {
		return middlesign.NewSignedRequestMiddleware(h, http.HandlerFunc(myErrorHandler), middlesignConfig)
	}

	myHandler := http.HandlerFunc(myApp)

	chain := alice.New(signedRequestMiddlewhare).Then(myHandler)
	http.ListenAndServe(":8000", chain)
}
