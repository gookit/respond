package respond

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/gookit/easytpl"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/netutil/httpreq"
)

// Options for the Responder
type Options struct {
	Debug bool

	JSONIndent bool
	JSONPrefix string

	XMLIndent bool
	XMLPrefix string

	// template render options, please see easytpl.Options
	TplLayout   string
	TplDelims   easytpl.TplDelims
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

// OptionFn option function
type OptionFn func(opts *Options)

// Responder definition
type Responder struct {
	opts *Options
	// view renderer
	renderer *easytpl.Renderer
	// mark init is completed
	initialized bool
}

/*************************************************************
 * create and initialize
 *************************************************************/

// New responder instance
func New(fns ...OptionFn) *Responder {
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

			TplDelims:   easytpl.TplDelims{Left: "{{", Right: "}}"},
			TplSuffixes: []string{"tpl", "html"},
		},
	}

	// apply config options
	for _, fn := range fns {
		fn(r.opts)
	}
	return r
}

func NewInited(opFns ...OptionFn) *Responder {
	return New(opFns...).Initialize()
}

// NewInitialized create new instance and initialization it.
func NewInitialized(opFns ...OptionFn) *Responder {
	return New(opFns...).Initialize()
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
	v := easytpl.NewRenderer()
	v.Debug = opts.Debug
	v.Delims = opts.TplDelims
	v.Layout = opts.TplLayout
	v.FuncMap = opts.TplFuncMap
	v.ViewsDir = opts.TplViewsDir
	v.MustInit()

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

// LoadTemplateGlob load templates by glob
//
// Usage:
//
//	LoadTemplateGlob("views/*")
//	LoadTemplateGlob("views/**/*")
func (r *Responder) LoadTemplateGlob(pattern string) {
	r.renderer.LoadByGlob(pattern)
}

// LoadTemplateFiles load template files.
// Usage:
//
//	LoadTemplateFiles("path/file1.tpl", "path/file2.tpl")
func (r *Responder) LoadTemplateFiles(files ...string) {
	r.renderer.LoadFiles(files...)
}

/*************************************************************
 * render and response HTML
 *************************************************************/

// HTML render HTML template file to http.ResponseWriter
//
// NOTE: layout only available on enable the view layout
func (r *Responder) HTML(w http.ResponseWriter, status int, template string, v any, layout ...string) error {
	w.Header().Set(ContentType, r.opts.ContentHTML)
	w.WriteHeader(status)

	return r.renderer.Render(w, template, v, layout...)
}

// HTMLString render HTML template string to http.ResponseWriter
func (r *Responder) HTMLString(w http.ResponseWriter, status int, tplContent string, v any) error {
	w.Header().Set(ContentType, r.opts.ContentHTML)
	w.WriteHeader(status)

	t := template.Must(template.New("temp").Parse(tplContent))
	if err := t.Execute(w, v); err != nil {
		// return err
		panic(err)
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
func (r *Responder) Auto(w http.ResponseWriter, req *http.Request, data any) error {
	resContentType := w.Header().Get(ContentType)
	if resContentType != "" {
		// TODO: check the accepted header
	}

	accepted := req.Header.Get("Accepted")
	if accepted == "" {
		return nil
	}

	ls := httpreq.ParseAccept(accepted)
	for _, val := range ls {
		if len(val) == 0 {
			continue
		}

		switch val {
		case "application/json":
			return r.JSON(w, http.StatusOK, data)
		case "application/xml":
			return r.XML(w, http.StatusOK, data)
		case "text/html":
			return r.HTML(w, http.StatusOK, "index", data)
		case "text/plain":
			return r.Text(w, http.StatusOK, goutil.String(data))
		}
		break
	}

	// TODO default response text
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
func (r *Responder) Data(w http.ResponseWriter, status int, v any, contentType string) error {
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
func (r *Responder) JSON(w http.ResponseWriter, status int, v any) error {
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

// JSONP serve data as JSONP response
func (r *Responder) JSONP(w http.ResponseWriter, status int, callback string, v any) error {
	w.Header().Set(ContentType, r.opts.ContentJSONP)
	w.WriteHeader(status)

	bs, err := jsonMarshal(v, false, false)
	if err != nil {
		return err
	}

	if callback == "" {
		return errors.New("renderer: callback can not bet empty")
	}

	_, _ = w.Write([]byte(callback + "("))
	_, err = w.Write(bs)
	_, _ = w.Write([]byte(");"))
	return err
}

// XML serve data as XML response
func (r *Responder) XML(w http.ResponseWriter, status int, v any) (err error) {
	w.Header().Set(ContentType, r.opts.ContentXML)
	w.WriteHeader(status)

	var bs []byte
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

// Binary serve data as Binary response.
// Usage:
//
//	var reader io.Reader
//	reader, _ = os.Open("./README.md")
//	r.Binary(w, http.StatusOK, reader, "readme.md", true)
func (r *Responder) Binary(w http.ResponseWriter, status int, in io.Reader, outName string, inline bool) error {
	bs, err := io.ReadAll(in)
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

// Options get readonly options
func (r *Responder) Options() Options {
	return *r.opts
}

// Renderer get template renderer
func (r *Responder) Renderer() *easytpl.Renderer {
	return r.renderer
}
