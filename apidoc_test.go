package apidoc

import "testing"

func TestInit(t *testing.T) {
	if err := Init(Project{
		DocumentTitle: "apidoc-test",
		DocumentPath:  "apidoc-test.html",
	}); err != nil {
		t.Fatal(err)
	}
	if p.DocumentTitle != "apidoc-test" {
		t.Fatal("DocumentTitle is not equal")
	}
	if p.DocumentPath != "apidoc-test.html" {
		t.Fatal("DocumentPath is not equal")
	}
	if err := p.deleteDocumentFile(); err != nil {
		t.Fatal(err)
	}
}
