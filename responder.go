package respond

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gookit/view"
)

// Options for the Responder
type Options struct {
	Debug bool

	JSONIndent bool
	JSONPrefix string

	XMLIndent bool
	XMLPrefix string

	// template render
	TplLayout   string
	TplDelims   view.TplDelims
	TplFuncMap  template.FuncMap
	TplViewsDir string
	TplSuffixes []string

	// supported content types
	ContentBinary, ContentHTML, ContentXML, ContentText, ContentJSON, ContentJSONP string

	// Charset default content data charset
	Charset string
	// AddCharset on response content
	AddCharset bool
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
			AddCharset:    true,

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

// NewInitialized create new instance and initialization it.
func NewInitialized(config func(*Options)) *Responder {
	return New(config).Initialize()
}

// Initialize the default instance
func Initialize(config func(opts *Options)) {
	// config options
	config(_def.opts)

	// init
	_def.Initialize()
}

// Initialize the responder
func (r *Responder) Initialize() *Responder {
	if r.initialized {
		return r
	}

	if r.opts.AddCharset {
		r.appendCharset()
	}

	opts := r.opts
	// init view renderer
	v := view.NewRenderer()
	v.Debug = opts.Debug
	v.Delims = opts.TplDelims
	v.Layout = opts.TplLayout
	v.FuncMap = opts.TplFuncMap
	v.ViewsDir = opts.TplViewsDir
	v.MustInitialize()

	r.renderer = v
	r.initialized = true
	return r
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

// LoadTemplateGlob data response
func LoadTemplateGlob(pattern string) {
	_def.LoadTemplateGlob(pattern)
}

// LoadTemplateGlob load templates by glob
// Usage:
// 		LoadTemplateGlob("views/*")
// 		LoadTemplateGlob("views/**/*")
func (r *Responder) LoadTemplateGlob(pattern string) {
	r.renderer.LoadByGlob(pattern)
}

// LoadTemplateFiles data response
func LoadTemplateFiles(files ...string) {
	_def.LoadTemplateFiles(files...)
}

// LoadTemplateFiles load template files.
// Usage:
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

// HTML data response
func HTML(w http.ResponseWriter, status int, template string, v interface{}, layout ...string) error {
	return _def.HTML(w, status, template, v, layout...)
}

// HTML render HTML template file to http.ResponseWriter
func (r *Responder) HTML(w http.ResponseWriter, status int, template string, v interface{}, layout ...string) error {
	w.Header().Set(ContentType, r.opts.ContentHTML)
	w.WriteHeader(status)

	return r.renderer.Render(w, template, v, layout...)
}

// HTMLString data response
func HTMLString(w http.ResponseWriter, status int, tplContent string, v interface{}) error {
	return _def.HTMLString(w, status, tplContent, v)
}

// HTMLString render HTML template string to http.ResponseWriter
func (r *Responder) HTMLString(w http.ResponseWriter, status int, tplContent string, v interface{}) error {
	w.Header().Set(ContentType, r.opts.ContentHTML)
	w.WriteHeader(status)

	t := template.Must(template.New("").Parse(tplContent))
	if err := t.Execute(w, v); err != nil {
		// return err
		panic(err)
	}

	return nil
}

// HTMLText data response
func HTMLText(w http.ResponseWriter, status int, html string) error {
	return _def.HTMLText(w, status, html)
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

// Auto data response
func Auto(w http.ResponseWriter, req *http.Request, v interface{}) error {
	return _def.Auto(w, req, v)
}

// Auto response data by request accepted header
func (r *Responder) Auto(w http.ResponseWriter, req *http.Request, data interface{}) error {
	resContentType := w.Header().Get(ContentType)
	if resContentType != "" {

	}

	accepted := req.Header.Get("Accepted")
	if accepted == "" {
		return nil
	}

	// default response text
	return nil
}

// Empty data response
func Empty(w http.ResponseWriter) error {
	return _def.Empty(w)
}

// Empty alias method of the NoContent()
func (r *Responder) Empty(w http.ResponseWriter) error {
	return r.NoContent(w)
}

// NoContent data response
func NoContent(w http.ResponseWriter) error {
	return _def.NoContent(w)
}

// NoContent serve success but no content response
func (r *Responder) NoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// Content data response
func Content(w http.ResponseWriter, status int, v []byte, contentType string) error {
	return _def.Content(w, status, v, contentType)
}

// Content serve success but no content response
func (r *Responder) Content(w http.ResponseWriter, status int, v []byte, contentType string) error {
	w.WriteHeader(status)
	w.Header().Set(ContentType, contentType)
	_, err := w.Write(v)
	return err
}

// Data data response
func Data(w http.ResponseWriter, status int, v interface{}, contentType string) error {
	return _def.Data(w, status, v, contentType)
}

// Data is the generic function called by XML, JSON, Data, HTML, and can be called by custom implementations.
func (r *Responder) Data(w http.ResponseWriter, status int, v interface{}, contentType string) error {
	w.WriteHeader(status)
	_, err := w.Write(v.([]byte))

	return err
}

// String data response
func String(w http.ResponseWriter, status int, v string) error {
	return _def.String(w, status, v)
}

// String alias method of the Text()
func (r *Responder) String(w http.ResponseWriter, status int, v string) error {
	return r.Text(w, status, v)
}

// Text data response
func Text(w http.ResponseWriter, status int, v string) error {
	return _def.Text(w, status, v)
}

// Text serve string content as text/plain response
func (r *Responder) Text(w http.ResponseWriter, status int, v string) error {
	w.WriteHeader(status)
	w.Header().Set(ContentType, r.opts.ContentText)
	_, err := w.Write([]byte(v))

	return err
}

// JSON data response
func JSON(w http.ResponseWriter, status int, v interface{}) error {
	return _def.JSON(w, status, v)
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
		_, err = w.Write([]byte(r.opts.JSONPrefix))
		if err != nil {
			return err
		}
	}

	_, err = w.Write(bs)
	return err
}

// JSONP data response
func JSONP(w http.ResponseWriter, status int, callback string, v interface{}) error {
	return _def.JSONP(w, status, callback, v)
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

	_,_ = w.Write([]byte(callback + "("))
	_, err = w.Write(bs)
	_,_ = w.Write([]byte(");"))

	return err
}

// XML data response
func XML(w http.ResponseWriter, status int, v interface{}) error {
	return _def.XML(w, status, v)
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
		_, err = w.Write([]byte(r.opts.XMLPrefix))
		if err != nil {
			return err
		}
	}

	_, err = w.Write(bs)
	return err
}

// Binary data response
func Binary(w http.ResponseWriter, status int, in io.Reader, outName string, inline bool) error {
	return _def.Binary(w, status, in, outName, inline)
}

// Binary serve data as Binary response.
// Usage:
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

// json converts the data as bytes using json encoder
func jsonMarshal(v interface{}, indent, unEscapeHTML bool) ([]byte, error) {
	var bs []byte
	var err error
	if indent {
		bs, err = MarshalIndent(v, "", "  ")
	} else {
		bs, err = Marshal(v)
	}

	if err != nil {
		return bs, err
	}

	if unEscapeHTML {
		bs = bytes.Replace(bs, []byte("\\u003c"), []byte("<"), -1)
		bs = bytes.Replace(bs, []byte("\\u003e"), []byte(">"), -1)
		bs = bytes.Replace(bs, []byte("\\u0026"), []byte("&"), -1)
	}

	return bs, nil
}
