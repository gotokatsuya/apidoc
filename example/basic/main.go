package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gotokatsuya/apidoc"
)

func apidocMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api := apidoc.NewAPI()
		api.SuppressedRequestHeaders("Cache-Control", "Content-Length", "X-Request-Id", "ETag", "Set-Cookie")
		api.ReadRequest(r, false)

		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, r)

		api.SuppressedResponseHeaders("Cache-Control", "Content-Length", "X-Request-Id", "X-Runtime", "X-XSS-Protection", "ETag")
		api.ReadResponseHeader(recorder.Header())
		api.WrapResponseBody(recorder.Body.Bytes())
		api.ResponseStatusCode = recorder.Code

		apidoc.GenerateDocument(api)

		for key, values := range recorder.Header() {
			for _, value := range values {
				w.Header().Set(key, value)
			}
		}
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
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

func init() {
	app := apidoc.NewApp("basic-example", "default.tpl.html", "basic-apidoc.html")
	if err := app.Init(); err != nil {
		panic(err)
	}
	apidoc.Setup(app)
}

func main() {
	log.Fatal(http.ListenAndServe(":8079", getHandler()))
}
