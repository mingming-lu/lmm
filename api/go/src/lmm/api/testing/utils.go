package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// StructToRequestBody converts a interface to a io.ReadClose
func StructToRequestBody(o interface{}) io.ReadCloser {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return ioutil.NopCloser(bytes.NewReader(b))
}
