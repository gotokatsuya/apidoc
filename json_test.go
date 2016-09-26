package apidoc

import (
	"encoding/json"
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	type data struct {
		Name string `json:"name"`
	}

	d := data{
		Name: "gotokatsuya",
	}
	in, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	out, err := PrettyPrint(in)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(out))
}
