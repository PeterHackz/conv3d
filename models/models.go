package models

import (
	"github.com/PeterHackz/conv3d/models/scw"
	"io"
	"os"
	"strings"
)

// Model wraps the models formats
//
// support for GLB/GLTF and SCW models is a priority
type Model interface {
	Load() error
}

var formats = map[string]func(data []byte) Model{
	"scw": func(data []byte) Model {
		return scw.New(data)
	},
}

// LoadFromFile Loads a model from a file
//
// for now, it checks the file type using the extension
//
// TODO: use header MAGIC to detect file type
func LoadFromFile(filename string) (Model, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	for k, v := range formats {
		if strings.HasSuffix(filename, k) {
			return v(data), nil
		}
	}
	return nil, nil
}
