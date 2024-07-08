package main

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/PeterHackz/conv3d/models"
)

func main() {

	inputFile := flag.String("in-file", "", "the input file path")
	outputFile := flag.String("out-file", "", "the output file")

	flag.Parse()

	if len(*inputFile) == 0 {
		panic("expected an input file")
	}

	outputJson := false

	if strings.HasSuffix(*inputFile, "scw") {
		outputJson = true
		if len(*outputFile) == 0 {
			*outputFile = "output.scw.json"
		}
	} else if strings.HasSuffix(*inputFile, "json") {
		if len(*outputFile) == 0 {
			*outputFile = "output.scw"
		}
	} else {
		if len(*outputFile) == 0 {
			*outputFile = "output.scw"
		}
	}

	model, err := models.LoadFromFile(*inputFile)
	if err != nil {
		panic(err)
	}

	if outputJson {
		if err := model.Load(); err != nil {
			panic(err)
		}
	} else {
		if err := model.LoadJSON(); err != nil {
			panic(err)
		}
	}

	r, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		panic(err)
	}

	file, err := os.Create(*outputFile)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	if outputJson {
		if _, err = file.WriteString(string(r)); err != nil {
			panic(err)
		}
	} else {
		if _, err = file.Write(model.Encode()); err != nil {
			panic(err)
		}
	}
}
