package main

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/PeterHackz/conv3d/models"
	"github.com/PeterHackz/conv3d/models/scw"
)

func main() {

	inputFile := flag.String("in-file", "", "the input file path")
	outputFile := flag.String("out-file", "", "the output file")
	scwSubVersion := flag.Int("minor-version", 5, "scw minor version (5 for v8+, not set for others)")
	scwOutVersion := flag.Int("out-version", 2, "scw output version")

	var scw2scw bool
	flag.BoolVar(&scw2scw, "scw2scw", false, "converts an scw model to another scw version")

	flag.Parse()

	if len(*inputFile) == 0 {
		panic("expected an input file")
	}

	outputJson := false

	if !scw2scw {
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
	} else {
		if len(*outputFile) == 0 {
			*outputFile = *inputFile
		}
	}

	model, err := models.LoadFromFile(*inputFile)
	if err != nil {
		panic(err)
	}

	if scw2scw || outputJson {

		if *scwSubVersion != 0 {
			switch m := model.(type) {
			case *scw.File:
				m.MinorVersion = *scwSubVersion
			}
		}

		if err := model.Load(); err != nil {
			panic(err)
		}

	} else {
		if err := model.LoadJSON(); err != nil {
			panic(err)
		}
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

	switch m := model.(type) {
	case *scw.File:
		if *scwOutVersion == 2 {
			m.Version = 2
			m.MinorVersion = 0
			if m.Unknown == -1 {
				m.Unknown = 0
			}
			for _, geom := range m.Geometries {
				for i := range geom.Vertices {
					if geom.Vertices[i].Name == "VERTEX" {
						geom.Vertices[i].Name = "POSITION"
					}
				}
			}
		} else {
			m.Unknown = -1
			m.MinorVersion = 5
			for _, geom := range m.Geometries {
				for i := range geom.Vertices {
					if geom.Vertices[i].Name == "POSITION" {
						geom.Vertices[i].Name = "VERTEX"
					}
				}
			}
		}
	}

	if !scw2scw && outputJson {
		r, err := json.MarshalIndent(model, "", "  ")
		if err != nil {
			panic(err)
		}

		if _, err = file.WriteString(string(r)); err != nil {
			panic(err)
		}
	} else {
		if _, err = file.Write(model.Encode()); err != nil {
			panic(err)
		}
	}
}
