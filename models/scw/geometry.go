package scw

import (
	"fmt"
	"math"
)

type Geometry struct {
	SCWFile *File `json:"-"`
	Name    string
	Group   string
	// this was used in older versions of scw format
	IgnoredMatrix Matrix4x4
	Vertices      []SourceArray
	BindMatrix    Matrix4x4
	Skins         struct {
		Joints              []string
		InverseBindMatrices []Matrix4x4
	}
	SkinWeights []Weight
	Materials   []IndexArray
}

type Weight struct {
	Joints  [4]byte
	Weights [4]uint16
}

func (w *Weight) Decode(reader *Reader) (err error) {
	for i := range 4 {
		if w.Joints[i], err = reader.ReadU8(); err != nil {
			return
		}
	}

	for i := range 4 {
		if w.Weights[i], err = reader.ReadU16(); err != nil {
			return
		}
	}

	return
}

type SourceArray struct {
	Name        string
	Index       byte
	SourceIndex byte    // uh, it is used for TEXTCOORD I think?
	Stride      byte    // stride(Color) / element size
	Scale       float32 // this can always be 0?? (not saying it is)
	Data        []int16 // vertex/coordinate data?
}

func (s *SourceArray) Decode(reader *Reader) (err error) {
	if s.Name, err = reader.ReadUTF(); err != nil {
		return
	}
	if s.Index, err = reader.ReadU8(); err != nil {
		return
	}
	if s.SourceIndex, err = reader.ReadU8(); err != nil {
		return
	}
	if s.Stride, err = reader.ReadU8(); err != nil {
		return
	}
	if s.Scale, err = reader.ReadFloat(); err != nil {
		return
	}
	var count uint32
	if count, err = reader.ReadU32(); err != nil {
		return
	}

	count *= uint32(s.Stride)

	s.Data = make([]int16, count)

	for i := range count {
		if s.Data[i], err = reader.ReadI16(); err != nil {
			return
		}
		s.Data[i] = int16(math.Round(float64(float32(s.Data[i]) * (s.Scale))))
	}
	return
}

type IndexArray struct {
	Name            string
	IndexBufferSize byte
	IndexBuffer     []uint32
}

func (i *IndexArray) Decode(reader *Reader) (err error) {
	if i.Name, err = reader.ReadUTF(); err != nil {
		return
	}
	var trianglesCount uint32
	if trianglesCount, err = reader.ReadU32(); err != nil {
		return
	}

	var inputsCount byte
	if inputsCount, err = reader.ReadU8(); err != nil {
		return
	}

	if i.IndexBufferSize, err = reader.ReadU8(); err != nil {
		return
	}

	totalIndices := 3 * trianglesCount * uint32(inputsCount)

	i.IndexBuffer = make([]uint32, totalIndices)

	switch i.IndexBufferSize {
	case 1:
		for v := range i.IndexBuffer {
			var idx uint8
			if idx, err = reader.ReadU8(); err != nil {
				return
			}
			i.IndexBuffer[v] = uint32(idx)
		}
	case 2:
		for v := range i.IndexBuffer {
			var idx uint16
			if idx, err = reader.ReadU16(); err != nil {
				return
			}
			i.IndexBuffer[v] = uint32(idx)
		}
	case 4:
		for v := range i.IndexBuffer {
			if i.IndexBuffer[v], err = reader.ReadU32(); err != nil {
				return
			}
		}
	default:
		return fmt.Errorf("unsupported index buffer size: %d", i.IndexBufferSize)
	}

	return
}

type Matrix4x4 [4][4]float32

func (m *Matrix4x4) Decode(reader *Reader) (err error) {
	for i := 0; i < 16; i++ {
		if m[i/4][i%4], err = reader.ReadFloat(); err != nil {
			return
		}
	}
	m.Transpose()
	return nil
}

func (m *Matrix4x4) Transpose() {
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			m[i][j], m[j][i] = m[j][i], m[i][j]
		}
	}
}

func (g *Geometry) Decode(reader *Reader) (err error) {
	if g.Name, err = reader.ReadUTF(); err != nil {
		return
	}

	if g.Group, err = reader.ReadUTF(); err != nil {
		return
	}

	if g.SCWFile.Version <= 1 {
		if err = g.IgnoredMatrix.Decode(reader); err != nil {
			return
		}
	}

	var verticesCount byte
	if verticesCount, err = reader.ReadU8(); err != nil {
		return
	}

	g.Vertices = make([]SourceArray, verticesCount)
	for i := range g.Vertices {
		if err = g.Vertices[i].Decode(reader); err != nil {
			return
		}
	}

	var b bool
	if b, err = reader.ReadBool(); err != nil {
		return err
	} else if b {
		if err = g.BindMatrix.Decode(reader); err != nil {
			return
		}
	}

	var skinsCount byte
	if skinsCount, err = reader.ReadU8(); err != nil {
		return err
	}

	g.Skins.Joints = make([]string, skinsCount)
	g.Skins.InverseBindMatrices = make([]Matrix4x4, skinsCount)

	for i := range skinsCount {
		if g.Skins.Joints[i], err = reader.ReadUTF(); err != nil {
			return err
		}
		if err = g.Skins.InverseBindMatrices[i].Decode(reader); err != nil {
			return err
		}
	}

	var skinWeightsCount uint32
	if skinWeightsCount, err = reader.ReadU32(); err != nil {
		return err
	}

	g.SkinWeights = make([]Weight, skinWeightsCount)
	for i := range skinWeightsCount {
		if err = g.SkinWeights[i].Decode(reader); err != nil {
			return
		}
	}

	var indexesCount byte
	if indexesCount, err = reader.ReadU8(); err != nil {
		return
	}

	g.Materials = make([]IndexArray, indexesCount)
	for i := range indexesCount {
		if err = g.Materials[i].Decode(reader); err != nil {
			return
		}
	}

	return
}
