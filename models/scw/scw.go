package scw

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
)

var (
	// ErrInvalidSCWMagic expected magic is to be SC3D (first 4 bytes of the File)
	ErrInvalidSCWMagic = errors.New("invalid scw file magic")
)

type File struct {
	reader *Reader
	Header
	Materials  []*Material
	Geometries []*Geometry
	Nodes      []Node
	Cameras    []*Camera3D
}

func New(data []byte) *File {
	return &File{
		reader: NewReader(data),
	}
}

func (f *File) LoadJSON() (err error) {
	if err = json.Unmarshal(f.reader.data, f); err != nil {
		return err
	}

	for _, geom := range f.Geometries {
		geom.SCWFile = f
	}

	for _, mat := range f.Materials {
		mat.SCWFile = f
	}

	for _, node := range f.Nodes {
		node.SCWFile = f
	}

	return nil
}

// Load loads a File, for now support of scw version 2 is under development
//
// in the future support for version 1 is planned.
func (f *File) Load() (err error) {
	reader := f.reader

	var magic string

	if magic, err = reader.ReadUTFWithLength(4); err != nil {
		return
	} else if magic != "SC3D" {
		return ErrInvalidSCWMagic
	}

	var (
		length uint32
		prop   string
	)

	end := false

	for !end {
		length, err = reader.ReadU32()
		if err != nil {
			return err
		}

		prop, err = reader.ReadUTFWithLength(4)
		reader.SkipBytes = int(length)
		if err != nil {
			return
		}
		switch prop {
		case "HEAD":
			if err = f.Header.Decode(reader); err != nil {
				return
			}
		case "MATE":
			material := new(Material)
			material.SCWFile = f
			if err = material.Decode(reader); err != nil {
				return
			}
			f.Materials = append(f.Materials, material)
		case "GEOM":
			geometry := new(Geometry)
			geometry.SCWFile = f
			if err = geometry.Decode(reader); err != nil {
				return
			}
			f.Geometries = append(f.Geometries, geometry)
		case "CAME":
			camera := new(Camera3D)
			if err = camera.Decode(reader); err != nil {
				return
			}
			f.Cameras = append(f.Cameras, camera)
		case "NODE":
			var nodesCount uint16
			if nodesCount, err = reader.ReadU16(); err != nil {
				return
			}
			f.Nodes = make([]Node, nodesCount)
			for i := range nodesCount {
				f.Nodes[i].SCWFile = f
				if err = f.Nodes[i].Decode(reader); err != nil {
					return
				}
			}
		case "WEND":
			reader.SkipBytes = 0
			end = true
		default:
			return fmt.Errorf("unsupported scw property: %s", prop)
		}
		if reader.SkipBytes != 0 {
			panic(fmt.Sprintf("failed to parse SCW model property infully, chunk: %s, bytes left: %d", prop, reader.SkipBytes))
		}

		_, err = reader.ReadU32() // crc32(prop + buffer)

		if err != nil {
			return
		}
	}

	return
}

type SC3DProperty interface {
	Encode(writer *Writer)
}

func encodeSC3DProperty(prop SC3DProperty, fileWriter *Writer, name string) {

	w := NewWriter()
	w.WriteStringChars(name)
	prop.Encode(w)
	bytes := w.Bytes()

	fileWriter.WriteU32(uint32(len(bytes) - len(name)))
	fileWriter.WriteBytes(bytes)

	checksum := crc32.ChecksumIEEE(bytes)
	fileWriter.WriteU32(checksum)

}

func (f *File) Encode() []byte {
	writer := NewWriter()

	writer.WriteStringChars("SC3D") // file magic

	encodeSC3DProperty(&f.Header, writer, "HEAD")

	for _, mat := range f.Materials {
		encodeSC3DProperty(mat, writer, "MATE")
	}

	for _, cam := range f.Cameras {
		encodeSC3DProperty(cam, writer, "CAME")
	}

	for _, geom := range f.Geometries {
		encodeSC3DProperty(geom, writer, "GEOM")
	}

	// TODO: separate Node encoding to a different function
	w := NewWriter()
	w.WriteStringChars("NODE")
	w.WriteU16(uint16(len(f.Nodes)))
	for _, node := range f.Nodes {
		node.Encode(w)
	}
	bytes := w.Bytes()

	writer.WriteU32(uint32(len(bytes) - len("NODE")))
	writer.WriteBytes(bytes)

	checksum := crc32.ChecksumIEEE(bytes)
	writer.WriteU32(checksum)

	writer.WriteU32(0)
	writer.WriteStringChars("WEND")
	writer.WriteU32(crc32.ChecksumIEEE([]byte("WEND")))

	return writer.Bytes()
}
