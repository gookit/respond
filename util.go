package respond

import "bytes"

// json converts the data as bytes using json encoder
func jsonMarshal(v any, indent, unEscapeHTML bool) ([]byte, error) {
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
