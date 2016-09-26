package apidoc

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// API has request and response info
type API struct {
	// Request
	RequestMethod            string            `json:"request_method"`
	RequestPath              string            `json:"request_path"`
	RequestHeaders           map[string]string `json:"request_headers"`
	RequestSuppressedHeaders map[string]bool   `json:"request_suppressed_headers"`
	RequestURLParams         map[string]string `json:"request_url_params"`
	RequestPostForms         map[string]string `json:"request_post_forms"`
	RequestBody              string            `json:"request_body"`

	// Response
	ResponseHeaders           map[string]string `json:"response_headers"`
	ResponseSuppressedHeaders map[string]bool   `json:"response_suppressed_headers"`
	ResponseStatusCode        int               `json:"response_status_code"`
	ResponseBody              string            `json:"response_body"`
}

// NewAPI new api instance
func NewAPI() API {
	return API{
		RequestHeaders:   map[string]string{},
		RequestURLParams: map[string]string{},
		RequestPostForms: map[string]string{},

		ResponseHeaders: map[string]string{},
	}
}

func (a API) equal(a2 API) bool {
	return a.RequestMethod == a2.RequestMethod && a.RequestPath == a2.RequestPath && a.ResponseStatusCode == a2.ResponseStatusCode
}

// SuppressedRequestHeaders ignore request headers
func (a *API) SuppressedRequestHeaders(headers ...string) {
	a.RequestSuppressedHeaders = make(map[string]bool, len(headers))
	for _, header := range headers {
		a.RequestSuppressedHeaders[header] = true
	}
}

// ReadRequestHeader read request http.Header
func (a *API) ReadRequestHeader(httpHeader http.Header) error {
	b := bytes.NewBuffer([]byte(""))
	if err := httpHeader.WriteSubset(b, a.RequestSuppressedHeaders); err != nil {
		return err
	}
	for _, header := range strings.Split(b.String(), "\n") {
		values := strings.Split(header, ":")
		if len(values) < 2 {
			continue
		}
		key := values[0]
		if key == "" {
			continue
		}
		a.RequestHeaders[key] = values[1]
	}
	return nil
}

// ReadRequestURLParams read request uri
func (a *API) ReadRequestURLParams(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	for _, param := range strings.Split(u.Query().Encode(), "&") {
		values := strings.Split(param, "=")
		if len(values) < 2 {
			continue
		}
		key := values[0]
		if key == "" {
			continue
		}
		a.RequestURLParams[key] = values[1]
	}
	return nil
}

// Reference https://golang.org/src/net/http/httputil/dump.go
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

// ReadRequestBody read request body
// Reference https://golang.org/src/net/http/httputil/dump.go
func (a *API) ReadRequestBody(req *http.Request) error {
	var err error
	save := req.Body
	if req.Body == nil {
		req.Body = nil
	} else {
		save, req.Body, err = drainBody(req.Body)
		if err != nil {
			return err
		}
	}

	var b bytes.Buffer

	if req.Body != nil {
		chunked := len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"
		var dest io.Writer = &b
		if chunked {
			dest = httputil.NewChunkedWriter(dest)
		}
		_, err = io.Copy(dest, req.Body)
		if chunked {
			err = dest.(io.Closer).Close()
		}
	}

	req.Body = save
	if err != nil {
		return err
	}

	contentType, ok := a.RequestHeaders["Content-Type"]
	if !ok {
		return nil
	}
	ct := strings.TrimSpace(contentType)
	switch {
	case strings.Contains(ct, "application/x-www-form-urlencoded"):
		for _, param := range strings.Split(b.String(), "&") {
			values := strings.Split(param, "=")
			if len(values) < 2 {
				continue
			}
			key := values[0]
			if key == "" {
				continue
			}
			a.RequestPostForms[key] = values[1]
		}
	case strings.Contains(ct, "application/json"):
		out, err := PrettyPrint(b.Bytes())
		if err != nil {
			return err
		}
		a.RequestBody = string(out)
	case strings.Contains(ct, "multipart/form-data"):
		// TODO handling multipart/form-data
	}
	return nil
}

func (a *API) getRequestURI(req *http.Request) string {
	reqURI := req.RequestURI
	if reqURI == "" {
		reqURI = req.URL.RequestURI()
	}
	return reqURI
}

// ReadRequest read values from http.Request
func (a *API) ReadRequest(req *http.Request, throwErr bool) error {
	a.RequestMethod = req.Method
	a.RequestPath = strings.Split(a.getRequestURI(req), "?")[0]
	if err := a.ReadRequestHeader(req.Header); err != nil {
		if throwErr {
			return err
		}
		log.Println(err)
	}
	if err := a.ReadRequestURLParams(a.getRequestURI(req)); err != nil {
		if throwErr {
			return err
		}
		log.Println(err)
	}
	if err := a.ReadRequestBody(req); err != nil {
		if throwErr {
			return err
		}
		log.Println(err)
	}
	return nil
}

// SuppressedResponseHeaders ignore response headers
func (a *API) SuppressedResponseHeaders(headers ...string) {
	a.ResponseSuppressedHeaders = make(map[string]bool, len(headers))
	for _, header := range headers {
		a.ResponseSuppressedHeaders[header] = true
	}
}

// ReadResponseHeader read http.Header
func (a *API) ReadResponseHeader(httpHeader http.Header) error {
	b := bytes.NewBuffer([]byte(""))
	if err := httpHeader.WriteSubset(b, a.ResponseSuppressedHeaders); err != nil {
		return err
	}
	for _, header := range strings.Split(b.String(), "\n") {
		values := strings.Split(header, ":")
		if len(values) < 2 {
			continue
		}
		key := values[0]
		if key == "" {
			continue
		}
		a.ResponseHeaders[key] = values[1]
	}
	return nil
}

// WrapResponseBody wrap body prettyprint if json
func (a *API) WrapResponseBody(body []byte) error {
	contentType, ok := a.ResponseHeaders["Content-Type"]
	if ok && strings.Contains(strings.TrimSpace(contentType), "application/json") {
		prettyBody, err := PrettyPrint(body)
		if err != nil {
			return err
		}
		a.ResponseBody = string(prettyBody)
		return nil
	}

	a.ResponseBody = string(body)

	return nil
}
