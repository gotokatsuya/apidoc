package main

import (
	"bytes"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gotokatsuya/apidoc"
)

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
	api := apidoc.NewAPI()
	api.SuppressedRequestHeaders("Cache-Control", "Content-Length", "X-Request-Id", "ETag", "Set-Cookie")
	api.ReadRequest(c.Request, false)

	gbw := newGinBodyWriter(c.Writer)
	c.Writer = gbw

	// before processing request
	c.Next()
	// after processing request

	api.SuppressedResponseHeaders("Cache-Control", "Content-Length", "X-Request-Id", "X-Runtime", "X-XSS-Protection", "ETag")
	api.ReadResponseHeader(c.Writer.Header())
	api.WrapResponseBody(gbw.Body())
	api.ResponseStatusCode = c.Writer.Status()

	// Handling error if you want
	_ = apidoc.GenerateDocument(api)
}

func getEngine() *gin.Engine {
	r := gin.Default()
	r.Use(apidocMiddleware)

	type user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	r.GET("/users", func(c *gin.Context) {
		limit := c.Query("limit")
		c.JSON(200, gin.H{
			"limit": limit,
			"users": []user{
				user{ID: 1, Name: "test1"},
				user{ID: 2, Name: "test2"}},
		})
	})
	r.GET("/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(200, gin.H{"user": user{ID: id, Name: "test"}})
	})
	r.POST("/users", func(c *gin.Context) {
		name := c.PostForm("name")
		c.JSON(200, gin.H{"user": user{ID: 1, Name: name}})
	})
	r.PUT("/users", func(c *gin.Context) {
		type req struct {
			Name string `form:"name" json:"name" binding:"required"`
		}
		var r req
		if err := c.BindJSON(&r); err != nil {
			log.Println(err)
			return
		}
		c.JSON(200, gin.H{"user": user{ID: 1, Name: r.Name}})
	})
	return r
}

func init() {
	app := apidoc.NewApp("gin-example", "custom.tpl.html", "custom-apidoc.html")
	if err := app.Init(); err != nil {
		panic(err)
	}
	apidoc.Setup(app)
}

func main() {
	// Listen and Server in 0.0.0.0:8080
	getEngine().Run(":8080")
}
