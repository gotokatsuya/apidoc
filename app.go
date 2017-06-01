package apidoc

import (
	"encoding/json"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

// App is
type App struct {
	Name string

	TemplateHTMLPath string
	OutputHTMLPath   string

	APIList []API
}

// NewApp is
func NewApp(name, templateHTMLPath, outputHTMLPath string) *App {
	return &App{
		Name:             name,
		TemplateHTMLPath: templateHTMLPath,
		OutputHTMLPath:   outputHTMLPath,
		APIList:          []API{},
	}
}

// OutputHTMLPathWithJSONExtension is
func (a *App) OutputHTMLPathWithJSONExtension() string {
	return a.OutputHTMLPath + ".json"
}

// OpenJSONHTMLPath is
func (a *App) OpenJSONHTMLPath() (*os.File, error) {
	filePath, err := filepath.Abs(a.OutputHTMLPathWithJSONExtension())
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// DecodeAPIList is
func (a *App) DecodeAPIList(jsonFile *os.File) error {
	if err := json.NewDecoder(io.Reader(jsonFile)).Decode(&a.APIList); err != nil {
		return err
	}
	return nil
}

// CreateJSONFile is
func (a *App) CreateJSONFile() (*os.File, error) {
	filePath, err := filepath.Abs(a.OutputHTMLPathWithJSONExtension())
	if err != nil {
		return nil, err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// RemoveJSONFile is
func (a *App) RemoveJSONFile() error {
	filePath, err := filepath.Abs(a.OutputHTMLPathWithJSONExtension())
	if err != nil {
		return err
	}
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

// WriteAPIList is
func (a *App) WriteAPIList(jsonFile *os.File) error {
	byteAry, err := json.Marshal(a.APIList)
	if err != nil {
		return err
	}
	out, err := JSONPrettyPrint(byteAry)
	if err != nil {
		return err
	}
	if _, err := jsonFile.Write(out); err != nil {
		return err
	}
	return nil
}

// CreateHTMLFile is
func (a *App) CreateHTMLFile() (*os.File, error) {
	filePath, err := filepath.Abs(a.OutputHTMLPath)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// RemoveHTMLFile is
func (a *App) RemoveHTMLFile() error {
	filePath, err := filepath.Abs(a.OutputHTMLPath)
	if err != nil {
		return err
	}
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

// ExecuteTemplate is
func (a *App) ExecuteTemplate(file *os.File) error {
	t := template.Must(template.ParseFiles(a.TemplateHTMLPath))
	return t.Execute(io.Writer(file), map[string]interface{}{
		"title": a.Name,
		"apis":  a.APIList,
	})
}

// Init initialize files
func (a *App) Init() error {
	jsonFile, err := a.CreateJSONFile()
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	if err := a.DecodeAPIList(jsonFile); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	htmlFile, err := a.CreateHTMLFile()
	if err != nil {
		return err
	}
	defer htmlFile.Close()
	if err := a.ExecuteTemplate(htmlFile); err != nil {
		return err
	}
	return nil
}

// Clear delete all files
func (a *App) Clear() error {
	if err := a.RemoveJSONFile(); err != nil {
		return err
	}
	if err := a.RemoveHTMLFile(); err != nil {
		return err
	}
	a.APIList = []API{}
	return nil
}

// AppendAPI is
func (a *App) AppendAPI(newAPI API) {
	for i, api := range a.APIList {
		if newAPI.Equal(api) {
			// replace new
			a.APIList[i] = newAPI
			return
		}
	}
	a.APIList = append(a.APIList, newAPI)
}

// GenerateDocument generate api document
func (a *App) GenerateDocument(api API) error {
	a.AppendAPI(api)

	file, err := a.CreateJSONFile()
	if err != nil {
		return err
	}
	defer file.Close()
	if err := a.WriteAPIList(file); err != nil {
		return err
	}

	file, err = a.CreateHTMLFile()
	if err != nil {
		return err
	}
	defer file.Close()
	if err := a.ExecuteTemplate(file); err != nil {
		return err
	}
	return nil
}
