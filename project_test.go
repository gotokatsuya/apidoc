package apidoc

import "testing"

func TestAppendAPI(t *testing.T) {
	p := Project{
		APIs: make([]API, 0),
	}
	a1 := NewAPI()
	a1.RequestMethod = "GET"
	a1.RequestPath = "/users"
	a1.RequestBody = "body1"
	p.appendAPI(a1)
	if len(p.APIs) != 1 {
		t.Fatal("API len is not 1")
	}

	a2 := NewAPI()
	a2.RequestMethod = "GET"
	a2.RequestPath = "/users"
	a2.RequestBody = "body2"
	p.appendAPI(a2)
	if len(p.APIs) != 1 {
		t.Fatal("API len is not 1")
	}

	a3 := NewAPI()
	a3.RequestMethod = "POST"
	a3.RequestPath = "/users"
	a3.RequestBody = "body3"
	p.appendAPI(a3)
	if len(p.APIs) != 2 {
		t.Fatal("API len is not 2")
	}
}
