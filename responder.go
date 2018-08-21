package respond

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gookit/view"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	defaultCharset   = "UTF-8"
	defaultXMLPrefix = `<?xml version="1.0" encoding="ISO-8859-1" ?>\n`
)

const (
	// ContentType header key
	ContentType = "Content-Type"

	// ContentText represents content type text/plain
	ContentText = "text/plain"
	// ContentJSON represents content type application/json
	ContentJSON = "application/json"
	// ContentJSONP represents content type application/javascript
	ContentJSONP = "application/javascript"
	// ContentXML represents content type application/xml
	ContentXML = "application/xml"
	// ContentYAML represents content type application/x-yaml
	// ContentYAML = "application/x-yaml"

	// ContentHTML represents content type text/html
	ContentHTML = "text/html"
	// ContentBinary represents content type application/octet-stream
	ContentBinary = "application/octet-stream"

	// ContentDisposition describes contentDisposition
	ContentDisposition = "Content-Disposition"

	// describes content disposition type
	dispositionInline = "inline"
	// describes content disposition type
	dispositionAttachment = "attachment"
)

// Options for the Responder
type Options struct {
	Debug bool

	JSONIndent bool
	JSONPrefix string

	XMLIndent bool
	XMLPrefix string

	// template render
	TplDir      string
	TplLayout   string
	TplDelims   view.TplDelims
	TplSuffixes []string
	TplFuncMap  template.FuncMap

	// supported content types
	ContentBinary, ContentHTML, ContentXML, ContentText, ContentJSON, ContentJSONP string

	// Charset default content data charset
	Charset string
	// AppendCharset on response content
	AppendCharset bool
}

// Responder definition
type Responder struct {
	opts *Options
	// view renderer
	renderer *view.Renderer
	// mark init is completed
	initialized bool
}

/*************************************************************
 * create and initialize
 *************************************************************/

// New instance
func New(config ...func(*Options)) *Responder {
	r := &Responder{
		opts: &Options{
			ContentXML:  ContentXML,
			ContentText: ContentText,
			ContentHTML: ContentHTML,
			ContentJSON: ContentJSON,

			ContentJSONP:  ContentJSONP,
			ContentBinary: ContentBinary,
			AppendCharset: true,

			Charset:   defaultCharset,
			XMLPrefix: defaultXMLPrefix,

			TplDelims:   view.TplDelims{Left: "{{", Right: "}}"},
			TplSuffixes: []string{"tpl", "html"},
		},
	}

	// apply user config
	if len(config) > 0 {
		config[0](r.opts)
	}

	return r
}

// Initialize the responder
func (r *Responder) Initialize() {
	if r.initialized {
		return
	}

	if r.opts.AppendCharset {
		r.appendCharset()
	}

	opts := r.opts
	// init view renderer
	v := view.NewRenderer()
	v.Debug = opts.Debug
	v.Delims = opts.TplDelims
	v.Layout = opts.TplLayout
	v.FuncMap = opts.TplFuncMap
	v.ViewsDir = opts.TplDir
	v.MustInitialize()

	r.renderer = v
	r.initialized = true
}

// append charset for all content types
func (r *Responder) appendCharset() {
	r.opts.ContentBinary += "; " + r.opts.Charset
	r.opts.ContentHTML += "; " + r.opts.Charset
	r.opts.ContentXML += "; " + r.opts.Charset
	r.opts.ContentText += "; " + r.opts.Charset
	r.opts.ContentJSON += "; " + r.opts.Charset
	r.opts.ContentJSONP += "; " + r.opts.Charset
}

// LoadTemplateGlob load templates by glob
// usage:
// 		LoadTemplateGlob("views/*")
// 		LoadTemplateGlob("views/**/*")
func (r *Responder) LoadTemplateGlob(pattern string) {
	r.renderer.LoadByGlob(pattern)
}

// LoadTemplateFiles load template files.
// usage:
// 		LoadTemplateFiles("path/file1.tpl", "path/file2.tpl")
func (r *Responder) LoadTemplateFiles(files ...string) {
	r.renderer.LoadFiles(files...)
}

// Renderer get view renderer
func (r *Responder) Renderer() *view.Renderer {
	return r.renderer
}

/*************************************************************
 * render and response HTML
 *************************************************************/

// HTML render HTML template file to http.ResponseWriter
func (r *Responder) HTML(w http.ResponseWriter, status int, template string, v interface{}, layout ...string) error {
	w.Header().Set(ContentType, r.opts.ContentHTML)
	w.WriteHeader(status)

	return r.renderer.Render(w, template, v, layout...)
}

// HTMLString render HTML template string to http.ResponseWriter
func (r *Responder) HTMLString(w http.ResponseWriter, status int, tplContent string, v interface{}) error {
	w.Header().Set(ContentType, r.opts.ContentHTML)
	w.WriteHeader(status)

	t := template.Must(template.New("").Parse(tplContent))
	if err := t.Execute(w, v); err != nil {
		panic(err)
		return err
	}

	return nil
}

// HTMLText response string as html/text
func (r *Responder) HTMLText(w http.ResponseWriter, status int, html string) error {
	w.Header().Set(ContentType, r.opts.ContentHTML)
	w.WriteHeader(status)

	_, err := w.Write([]byte(html))
	return err
}

/*************************************************************
 * respond data
 *************************************************************/

// Auto response data by request accepted header
func (r *Responder) Auto(w io.Writer, accepted string, data interface{}) error {
	return nil
}

// Empty alias method of the NoContent()
func (r *Responder) Empty(w http.ResponseWriter) error {
	return r.NoContent(w)
}

// NoContent serve success but no content response
func (r *Responder) NoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Content serve success but no content response
func (r *Responder) Content(w http.ResponseWriter, status int, v []byte, contentType string) error {
	w.WriteHeader(status)
	w.Header().Set(ContentType, contentType)
	_, err := w.Write(v)
	return err
}

// Data is the generic function called by XML, JSON, Data, HTML, and can be called by custom implementations.
func (r *Responder) Data(w http.ResponseWriter, status int, v interface{}, contentType string) error {
	w.WriteHeader(status)
	_, err := w.Write(v.([]byte))

	return err
}

// String alias method of the Text()
func (r *Responder) String(w http.ResponseWriter, status int, v string) error {
	return r.Text(w, status, v)
}

// Text serve string content as text/plain response
func (r *Responder) Text(w http.ResponseWriter, status int, v string) error {
	w.WriteHeader(status)
	w.Header().Set(ContentType, r.opts.ContentText)
	_, err := w.Write([]byte(v))

	return err
}

// JSON serve string content as json response
func (r *Responder) JSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set(ContentType, r.opts.ContentJSON)
	w.WriteHeader(status)

	bs, err := jsonMarshal(v, r.opts.JSONIndent, false)
	if err != nil {
		return err
	}

	if r.opts.JSONPrefix != "" {
		w.Write([]byte(r.opts.JSONPrefix))
	}

	_, err = w.Write(bs)
	return err
}

// JSONP serve data as JSONP response
func (r *Responder) JSONP(w http.ResponseWriter, status int, callback string, v interface{}) error {
	w.Header().Set(ContentType, r.opts.ContentJSONP)
	w.WriteHeader(status)

	bs, err := jsonMarshal(v, false, false)
	if err != nil {
		return err
	}

	if callback == "" {
		return errors.New("renderer: callback can not bet empty")
	}

	w.Write([]byte(callback + "("))
	_, err = w.Write(bs)
	w.Write([]byte(");"))

	return err
}

// XML serve data as XML response
func (r *Responder) XML(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set(ContentType, r.opts.ContentXML)
	w.WriteHeader(status)
	var bs []byte
	var err error

	if r.opts.XMLIndent {
		bs, err = xml.MarshalIndent(v, "", " ")
	} else {
		bs, err = xml.Marshal(v)
	}
	if err != nil {
		return err
	}

	if r.opts.XMLPrefix != "" {
		w.Write([]byte(r.opts.XMLPrefix))
	}

	_, err = w.Write(bs)
	return err
}

// Binary serve data as Binary response.
// usage:
// 		var reader io.Reader
// 		reader, _ = os.Open("./README.md")
// 		r.Binary(w, http.StatusOK, reader, "readme.md", true)
func (r *Responder) Binary(w http.ResponseWriter, status int, in io.Reader, outName string, inline bool) error {
	bs, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	dispositionType := dispositionAttachment
	if inline {
		dispositionType = dispositionInline
	}

	w.Header().Set(ContentType, r.opts.ContentBinary)
	w.Header().Set(ContentDisposition, fmt.Sprintf("%s; filename=%s", dispositionType, outName))
	w.WriteHeader(status)

	_, err = w.Write(bs)
	return err
}

/*************************************************************
 * helper methods
 *************************************************************/

// Options get options
func (r *Responder) Options() Options {
	return *r.opts
}
