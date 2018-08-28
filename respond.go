/*
Package respond is a http data responder

Source code and other details for the project are available at GitHub:

	https://github.com/gookit/respond

usage please see README and examples.
*/
package respond

import (
	"io"
	"net/http"
)

// create a default instance
var As = New()

// Initialize the default instance
func Initialize(config func(opts *Options)) {
	// config options
	config(As.opts)

	// init
	As.Initialize()
}

// Auto data response
func Auto(w http.ResponseWriter, req *http.Request, v interface{}) error {
	return As.Auto(w, req, v)
}

// Binary data response
func Binary(w http.ResponseWriter, status int, in io.Reader, outName string, inline bool) error {
	return As.Binary(w, status, in, outName, inline)
}

// Content data response
func Content(w http.ResponseWriter, status int, v []byte, contentType string) error {
	return As.Content(w, status, v, contentType)
}

// Data data response
func Data(w http.ResponseWriter, status int, v interface{}, contentType string) error {
	return As.Data(w, status, v, contentType)
}

// Empty data response
func Empty(w http.ResponseWriter) error {
	return As.Empty(w)
}

// HTML data response
func HTML(w http.ResponseWriter, status int, template string, v interface{}, layout ...string) error {
	return As.HTML(w, status, template, v, layout...)
}

// HTMLString data response
func HTMLString(w http.ResponseWriter, status int, tplContent string, v interface{}) error {
	return As.HTMLString(w, status, tplContent, v)
}

// HTMLText data response
func HTMLText(w http.ResponseWriter, status int, html string) error {
	return As.HTMLText(w, status, html)
}

// JSON data response
func JSON(w http.ResponseWriter, status int, v interface{}) error {
	return As.JSON(w, status, v)
}

// JSONP data response
func JSONP(w http.ResponseWriter, status int, callback string, v interface{}) error {
	return As.JSONP(w, status, callback, v)
}

// LoadTemplateFiles data response
func LoadTemplateFiles(files ...string) {
	As.LoadTemplateFiles(files...)
}

// LoadTemplateGlob data response
func LoadTemplateGlob(pattern string) {
	As.LoadTemplateGlob(pattern)
}

// NoContent data response
func NoContent(w http.ResponseWriter) error {
	return As.NoContent(w)
}

// String data response
func String(w http.ResponseWriter, status int, v string) error {
	return As.String(w, status, v)
}

// Text data response
func Text(w http.ResponseWriter, status int, v string) error {
	return As.Text(w, status, v)
}

// XML data response
func XML(w http.ResponseWriter, status int, v interface{}) error {
	return As.XML(w, status, v)
}
