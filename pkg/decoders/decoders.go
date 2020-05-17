package decoders

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// StringDecoder decodes the HTTP response into a string.
type StringDecoder struct {
}

// Decode implements the decoders.ResponseDecoder interface.
func (d StringDecoder) Decode(resp *http.Response, v interface{}) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	s := v.(*string)
	*s = string(b)

	return nil
}

// XMLDecoder decodes the HTTP response into a XML-tagged struct value.
type XMLDecoder struct {
}

// Decode implements the decoders.ResponseDecoder interface.
func (d XMLDecoder) Decode(resp *http.Response, v interface{}) error {
	return xml.NewDecoder(resp.Body).Decode(v)
}
