/*
Package respond is a http data responder

Source code and other details for the project are available at GitHub:

	https://github.com/gookit/respond

usage please see README and examples.
*/
package respond

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gookit/goutil/netutil/httpctype"
)

const (
	defaultCharset   = "UTF-8"
	defaultXMLPrefix = `<?xml version="1.0" encoding="ISO-8859-1" ?>\n`
)

const (
	// ContentType header key
	ContentType = httpctype.Key

	// ContentText represents content type text/plain
	ContentText = httpctype.Text
	// ContentJSON represents content type application/json
	ContentJSON = httpctype.JSON
	// ContentJSONP represents content type application/javascript
	ContentJSONP = httpctype.JS
	// ContentXML represents content type application/xml
	ContentXML = httpctype.XML
	// ContentYAML represents content type application/x-yaml
	// ContentYAML = "application/x-yaml"

	// ContentHTML represents content type text/html
	ContentHTML = httpctype.HTML
	// ContentBinary represents content type application/octet-stream
	ContentBinary = httpctype.Binary
)

const (
	// ContentDisposition describes contentDisposition
	ContentDisposition = "Content-Disposition"

	// describes content disposition type
	dispositionInline = "inline"
	// describes content disposition type
	dispositionAttachment = "attachment"
)

/*************************************************************
 * JSON driver config
 *************************************************************/

var (
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	JSONDecoder   = json.NewDecoder
)

/*************************************************************
 * default instance
 *************************************************************/

// create a default instance
var std = New()

// Default instance get
func Default() *Responder {
	return std
}

// Init the default instance
func Init(fns ...OptionFn) {
	Initialize(fns...)
}

// Initialize the default instance
func Initialize(fns ...OptionFn) {
	// apply config options
	for _, fn := range fns {
		fn(std.opts)
	}

	// initialize instance
	std.Initialize()
}

// LoadTemplateGlob data response
func LoadTemplateGlob(pattern string) {
	std.LoadTemplateGlob(pattern)
}

// LoadTemplateFiles data response
func LoadTemplateFiles(files ...string) {
	std.LoadTemplateFiles(files...)
}

// HTML data response
func HTML(w http.ResponseWriter, status int, template string, v any, layout ...string) error {
	return std.HTML(w, status, template, v, layout...)
}

// HTMLString data response
func HTMLString(w http.ResponseWriter, status int, tplContent string, v any) error {
	return std.HTMLString(w, status, tplContent, v)
}

// HTMLText output raw HTML contents response
func HTMLText(w http.ResponseWriter, status int, html string) error {
	return std.HTMLText(w, status, html)
}

// Auto data response
func Auto(w http.ResponseWriter, req *http.Request, v any) error {
	return std.Auto(w, req, v)
}

// Empty data response
func Empty(w http.ResponseWriter) error {
	return std.Empty(w)
}

// NoContent data response
func NoContent(w http.ResponseWriter) error {
	return std.NoContent(w)
}

// Content data response
func Content(w http.ResponseWriter, status int, v []byte, contentType string) error {
	return std.Content(w, status, v, contentType)
}

// Data data response
func Data(w http.ResponseWriter, status int, v any, contentType string) error {
	return std.Data(w, status, v, contentType)
}

// String data response
func String(w http.ResponseWriter, status int, v string) error {
	return std.String(w, status, v)
}

// Text data response
func Text(w http.ResponseWriter, status int, v string) error {
	return std.Text(w, status, v)
}

// JSON data response
func JSON(w http.ResponseWriter, status int, v any) error {
	return std.JSON(w, status, v)
}

// JSONP data response
func JSONP(w http.ResponseWriter, status int, callback string, v any) error {
	return std.JSONP(w, status, callback, v)
}

// XML data response
func XML(w http.ResponseWriter, status int, v any) error {
	return std.XML(w, status, v)
}

// Binary data response
func Binary(w http.ResponseWriter, status int, in io.Reader, outName string, inline bool) error {
	return std.Binary(w, status, in, outName, inline)
}
