# apidoc

Automatic API Document Generator.

## Supporting web framework

- gin

## Usage

### gin
  
```go
type ginBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w ginBodyWriter) Body() []byte {
	return w.body.Bytes()
}

func (w ginBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ginBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func newGinBodyWriter(w gin.ResponseWriter) *ginBodyWriter {
	return &ginBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: w}
}

func apidocMiddleware(c *gin.Context) {
	if apidoc.IsDisabled() {
		return
	}

	api := apidoc.NewAPI()
	// Ignore header names
	api.SuppressedRequestHeaders("Cache-Control", "Content-Length", "X-Request-Id", "ETag", "Set-Cookie")
	api.ReadRequest(c.Request, false)

	// Need to implement own Writer intercepting Write(), WriteString() calls to get response body.
	gbw := newGinBodyWriter(c.Writer)
	c.Writer = gbw

	// before processing request

	c.Next()

	// after processing request

	// Ignore header names
	api.SuppressedResponseHeaders("Cache-Control", "Content-Length", "X-Request-Id", "X-Runtime", "X-XSS-Protection", "ETag")
	api.ReadResponseHeader(c.Writer.Header())
	api.WrapResponseBody(gbw.Body())
	api.ResponseStatusCode = c.Writer.Status()

	apidoc.Gen(api)
}

func getEngine() *gin.Engine {
	r := gin.Default()
	r.Use(apidocMiddleware)

	...

	return r
}

func init() {
	d := flag.Bool("d", false, "disable api doc")
	flag.Parse()
	if *d {
		apidoc.Disable()
	}
	if apidoc.IsDisabled() {
		return
	}
	apidoc.Init(apidoc.Project{
		DocumentTitle: "readme",
		DocumentPath:  "readme-apidoc.html",
		TemplatePath:  "readme.tpl.html",
	})
}

func main() {
	// Listen and Server in 0.0.0.0:8080
	getEngine().Run(":8080")
}
```

## View

![view.png](https://github.com/gotokatsuya/apidoc/blob/master/example/gin/view.v1.png)

https://gotokatsuya.github.io/apidoc/example/gin/custom-apidoc.html
