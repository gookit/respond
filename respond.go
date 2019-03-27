/*
Package respond is a http data responder

Source code and other details for the project are available at GitHub:

	https://github.com/gookit/respond

usage please see README and examples.
*/
package respond

import "encoding/json"

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

/*************************************************************
 * JSON driver config
 *************************************************************/

var (
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
)

// create a default instance
var _def = New()

// Default instance get
func Default() *Responder {
	return _def
}