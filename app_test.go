package apidoc

import (
	"testing"
)

func TestAppendAPI(t *testing.T) {
	app := NewApp("", "", "")

	a1 := NewAPI()
	a1.RequestMethod = "GET"
	a1.RequestPath = "/users"
	a1.RequestBody = "body1"
	app.AppendAPI(a1)
	if len(app.APIList) != 1 {
		t.Fatalf("%v is invalid length", len(app.APIList))
		return
	}

	a2 := NewAPI()
	a2.RequestMethod = "GET"
	a2.RequestPath = "/users"
	a2.RequestBody = "body2"
	app.AppendAPI(a2)
	if len(app.APIList) != 1 {
		t.Fatalf("%v is invalid length", len(app.APIList))
		return
	}

	a3 := NewAPI()
	a3.RequestMethod = "POST"
	a3.RequestPath = "/users"
	a3.RequestBody = "body3"
	app.AppendAPI(a3)
	if len(app.APIList) != 2 {
		t.Fatalf("%v is invalid length", len(app.APIList))
		return
	}
}
