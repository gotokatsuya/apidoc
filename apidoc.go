package apidoc

var (
	// disable not gen docs if true
	disable bool

	// global project
	p Project
)

// Init initialize project setting
func Init(newProject Project) error {
	p = newProject
	p.APIs = []API{}
	if err := p.loadDocumentJSONFile(); err != nil {
		return err
	}
	if err := p.writeDocumentFile(); err != nil {
		return err
	}
	return nil
}

// Enable enable generator
func Enable() {
	disable = false
}

// Disable disable generator
func Disable() {
	disable = true
}

// IsDisabled ref disable
func IsDisabled() bool {
	return disable
}

// Clear delete all files
func Clear() error {
	if err := p.deleteDocumentJSONFile(); err != nil {
		return err
	}
	if err := p.deleteDocumentFile(); err != nil {
		return err
	}
	p.APIs = []API{}
	return nil
}

// Gen generate api document
func Gen(api API) error {
	p.appendAPI(api)
	if err := p.writeDocumentJSONFile(); err != nil {
		return err
	}
	if err := p.writeDocumentFile(); err != nil {
		return err
	}
	return nil
}
