package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gotokatsuya/apidoc"
)

func apidocMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apidoc.IsDisabled() {
			handler.ServeHTTP(w, r)
			return
		}

		api := apidoc.NewAPI()
		api.SuppressedRequestHeaders("Cache-Control", "Content-Length", "X-Request-Id", "ETag", "Set-Cookie")
		api.ReadRequest(r, false)

		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, r)

		api.SuppressedResponseHeaders("Cache-Control", "Content-Length", "X-Request-Id", "X-Runtime", "X-XSS-Protection", "ETag")
		api.ReadResponseHeader(recorder.Header())
		api.WrapResponseBody(recorder.Body.Bytes())
		api.ResponseStatusCode = recorder.Code

		apidoc.Gen(api)

		for key, values := range recorder.Header() {
			for _, value := range values {
				w.Header().Set(key, value)
			}
		}
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
	})
}

func init() {
	d := flag.Bool("d", false, "disable api doc")
	flag.Parse()
	if *d {
		apidoc.Disable()
	}
	if apidoc.IsDisabled() {
		return
	}
	apidoc.Init(apidoc.Project{
		DocumentTitle: "example-basic",
		DocumentPath:  "basic-apidoc.html",
	})
}

func getHandler() http.Handler {
	mux := http.DefaultServeMux
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello")
	})
	return apidocMiddleware(mux)
}

func main() {
	log.Fatal(http.ListenAndServe(":8079", getHandler()))
}
