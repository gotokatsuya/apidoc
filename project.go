package apidoc

import (
	"encoding/json"
	"go/build"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
)

// Project has project setting
type Project struct {
	DocumentTitle string
	DocumentPath  string
	TemplatePath  string

	APIs []API
}

func (p *Project) hasDocumentPath() bool {
	return p.DocumentPath != ""
}

func (p *Project) getDocumentPath() string {
	if p.hasDocumentPath() {
		return p.DocumentPath
	}
	return "apidoc.html"
}

func (p *Project) getDocumentJSONPath() string {
	return p.getDocumentPath() + ".json"
}

func (p *Project) hasTemplatePath() bool {
	return p.TemplatePath != ""
}

func findAppPath() string {
	const appName = "github.com/gotokatsuya/apidoc"
	appPkg, err := build.Import(appName, "", build.FindOnly)
	if err != nil {
		return ""
	}
	return path.Join(appPkg.SrcRoot, appName)
}

func (p *Project) getTemplatePath() string {
	if p.hasTemplatePath() {
		return p.TemplatePath
	}
	return path.Join(findAppPath(), "default.tpl.html")
}

func (p *Project) openDocumentJSONFile() (*os.File, error) {
	filePath, err := filepath.Abs(p.getDocumentJSONPath())
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (p *Project) loadDocumentJSONFile() error {
	file, err := p.openDocumentJSONFile()
	defer file.Close()
	if err != nil {
		return err
	}
	if err := json.NewDecoder(io.Reader(file)).Decode(&p.APIs); err != nil {
		return err
	}
	return nil
}

func (p *Project) createDocumentJSONFile() (*os.File, error) {
	filePath, err := filepath.Abs(p.getDocumentJSONPath())
	if err != nil {
		return nil, err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (p *Project) deleteDocumentJSONFile() error {
	filePath, err := filepath.Abs(p.getDocumentJSONPath())
	if err != nil {
		return err
	}
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

func (p *Project) writeDocumentJSONFile() error {
	file, err := p.createDocumentJSONFile()
	defer file.Close()
	if err != nil {
		return err
	}
	b, err := json.Marshal(p.APIs)
	if err != nil {
		return err
	}
	out, err := PrettyPrint(b)
	if err != nil {
		return err
	}
	if _, err := file.Write(out); err != nil {
		return err
	}
	return nil
}

func (p *Project) createDocumentFile() (*os.File, error) {
	filePath, err := filepath.Abs(p.getDocumentPath())
	if err != nil {
		return nil, err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (p *Project) deleteDocumentFile() error {
	filePath, err := filepath.Abs(p.getDocumentPath())
	if err != nil {
		return err
	}
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

func (p *Project) writeDocumentFile() error {
	t := template.Must(template.ParseFiles(p.getTemplatePath()))
	file, err := p.createDocumentFile()
	defer file.Close()
	if err != nil {
		return err
	}
	return t.Execute(io.Writer(file), map[string]interface{}{
		"title": p.DocumentTitle,
		"apis":  p.APIs,
	})
}

func (p *Project) appendAPI(newAPI API) {
	for i, api := range p.APIs {
		if newAPI.equal(api) {
			// replace
			p.APIs[i] = newAPI
			return
		}
	}
	p.APIs = append(p.APIs, newAPI)
}
