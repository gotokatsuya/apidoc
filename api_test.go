package apidoc

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestSuppressedRequestHeaders(t *testing.T) {
	api := NewAPI()
	api.SuppressedRequestHeaders("Cache-Control")
	if !api.RequestSuppressedHeaders["Cache-Control"] {
		t.Fatal("Cache-Control is not set")
	}
}

func TestReadRequestHeader(t *testing.T) {
	api := NewAPI()
	var header http.Header = make(map[string][]string)
	header.Set("X-Name", "gotokatsuya")
	if err := api.ReadRequestHeader(header); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(api.RequestHeaders["X-Name"]) != "gotokatsuya" {
		t.Fatal("X-Name is not equal")
	}
}

func TestReadRequestURLParams(t *testing.T) {
	api := NewAPI()
	uri := "http://localhost:8080/hello?key=world"
	if err := api.ReadRequestURLParams(uri); err != nil {
		t.Fatal(err)
	}
	if api.RequestURLParams["key"] != "world" {
		t.Fatal("key is not equal")
	}
}

func TestSuppressedResponseHeaders(t *testing.T) {
	api := NewAPI()
	api.SuppressedResponseHeaders("Cache-Control")
	if !api.ResponseSuppressedHeaders["Cache-Control"] {
		t.Fatal("Cache-Control is not set")
	}
}

func TestReadResponseHeader(t *testing.T) {
	api := NewAPI()
	var header http.Header = make(map[string][]string)
	header.Set("X-Name", "gotokatsuya")
	if err := api.ReadResponseHeader(header); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(api.ResponseHeaders["X-Name"]) != "gotokatsuya" {
		t.Fatal("X-Name is not equal")
	}
}

func TestWrapResponseBody(t *testing.T) {
	api := NewAPI()

	type res struct {
		Name string `json:"name"`
	}
	r := res{
		Name: "gotokatsuya",
	}
	in, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	api.ResponseHeaders["Content-Type"] = "application/json"
	if err := api.WrapResponseBody(in); err != nil {
		t.Fatal(err)
	}
	t.Log(api.ResponseBody)
}
