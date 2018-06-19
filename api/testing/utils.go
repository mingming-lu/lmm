package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"lmm/api/domain/factory"
)

func GenerateID() uint64 {
	id, err := factory.Default().GenerateID()
	if err != nil {
		panic(err)
	}
	return id
}

func StructToRequestBody(o interface{}) io.ReadCloser {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return ioutil.NopCloser(bytes.NewReader(b))
}
