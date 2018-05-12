package category

import (
	"encoding/json"
	"io"
)

type Category struct {
	ID   uint64 `json:"id"`
	User int64  `json:"user"`
	Name string `json:"name"`
}

func Read(reader io.Reader) (*Category, error) {
	model := Category{}
	if err := json.NewDecoder(reader).Decode(&model); err != nil {
		return nil, err
	}
	return &model, nil
}
