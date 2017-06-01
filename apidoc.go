package apidoc

var app *App

func Setup(newApp *App) {
	app = newApp
}

func GenerateDocument(api API) error {
	return app.GenerateDocument(api)
}
