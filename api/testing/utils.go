package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

func StructToRequestBody(o interface{}) io.ReadCloser {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return ioutil.NopCloser(bytes.NewReader(b))
}
