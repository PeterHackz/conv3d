package main

import (
	"encoding/json"
	"os"

	"github.com/PeterHackz/conv3d/models"
)

func main() {
	model, err := models.LoadFromFile("test_geo.scw")
	if err != nil {
		panic(err)
	}
	if err := model.Load(); err != nil {
		panic(err)
	}
	r, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		panic(err)
	}

	file, err := os.Create("output.json")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err = file.WriteString(string(r)); err != nil {
		panic(err)
	}
}
