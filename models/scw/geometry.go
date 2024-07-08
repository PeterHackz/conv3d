package scw

import (
	"fmt"
)

type Geometry struct {
	SCWFile *File `json:"-"`
	Name    string
	Group   string
	// this was used in older versions of scw format
	IgnoredMatrix Matrix4x4
	Vertices      []SourceArray
	HasBindMatrix bool
	BindMatrix    Matrix4x4
	Skins         struct {
		Joints              []string
		InverseBindMatrices []Matrix4x4
	}
	SkinWeights []Weight
	Materials   []IndexArray
}

func (g *Geometry) Tag() string {
	return "GEOM"
}

type Weight struct {
	Joints  [4]byte
	Weights [4]uint16
}

func (w *Weight) Decode(reader *Reader, scwVersion uint16, scwMinorVersion int) (err error) {
	for i := range 4 {
		if w.Joints[i], err = reader.ReadU8(); err != nil {
			return
		}
	}

	var readFn func() (uint16, error)

	// scw v0 needs more work
	// (it's support is planned)
	if scwVersion == 0 && scwMinorVersion != 5 {
		readFn = func() (uint16, error) {
			res, err := reader.ReadU8()
			return uint16(res), err
		}
	} else {
		readFn = reader.ReadU16
	}

	for i := range 4 {
		if w.Weights[i], err = readFn(); err != nil {
			return
		}
	}

	return
}

func (w *Weight) Encode(writer *Writer, scwVersion uint16, scwMinorVersion int) {
	writer.WriteBytes(w.Joints[:])

	var writeFn func(uint16)
	_ = writeFn

	if scwVersion == 0 && scwMinorVersion != 5 {
		writeFn = func(val uint16) {
			writer.WriteU8(byte(val))
		}
	} else {
		writeFn = writer.WriteU16
	}

	for i := range 4 {
		writeFn(w.Weights[i])
	}
}

type SourceArray struct {
	Name        string
	Index       byte
	SourceIndex byte      // uh, it is used for TEXTCOORD I think?
	Stride      byte      // stride(Color) / element size
	Scale       float32   // this can always be 0?? (not saying it is)
	Data        []float64 // vertex/coordinate data?
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

	s.Data = make([]float64, count)

	var val int16
	for i := range count {
		if val, err = reader.ReadI16(); err != nil {
			return
		}
		s.Data[i] = float64(val) * float64(s.Scale)
	}

	return
}

func (s *SourceArray) Encode(writer *Writer) {
	writer.WriteStringUTF(s.Name)

	writer.WriteU8(s.Index)
	writer.WriteU8(s.SourceIndex)
	writer.WriteU8(s.Stride)
	writer.WriteFloat(s.Scale)

	count := uint32(len(s.Data))
	count /= uint32(s.Stride)

	writer.WriteU32(count)

	for _, data := range s.Data {
		val := int16(data / float64(s.Scale))
		writer.WriteI16(val)
	}

}

type IndexArray struct {
	Name            string
	IndexBufferSize byte
	IndexBuffer     []uint32
	TrianglesCount  uint32
	InputsCount     byte
}

func (i *IndexArray) Decode(reader *Reader) (err error) {
	if i.Name, err = reader.ReadUTF(); err != nil {
		return
	}

	if i.TrianglesCount, err = reader.ReadU32(); err != nil {
		return
	}

	if i.InputsCount, err = reader.ReadU8(); err != nil {
		return
	}

	if i.IndexBufferSize, err = reader.ReadU8(); err != nil {
		return
	}

	totalIndices := 3 * i.TrianglesCount * uint32(i.InputsCount)

	i.IndexBuffer = make([]uint32, totalIndices)

	for v := range i.IndexBuffer {
		switch i.IndexBufferSize {
		case 1:
			var idx uint8
			if idx, err = reader.ReadU8(); err != nil {
				return
			}
			i.IndexBuffer[v] = uint32(idx)
		case 2:
			var idx uint16
			if idx, err = reader.ReadU16(); err != nil {
				return
			}
			i.IndexBuffer[v] = uint32(idx)
		case 4:
			if i.IndexBuffer[v], err = reader.ReadU32(); err != nil {
				return
			}
		default:
			return fmt.Errorf("unsupported index buffer size: %d", i.IndexBufferSize)
		}
	}

	return
}

func (i *IndexArray) Encode(writer *Writer) {
	writer.WriteStringUTF(i.Name)

	writer.WriteU32(i.TrianglesCount)
	writer.WriteU8(i.InputsCount)
	writer.WriteU8(i.IndexBufferSize)

	for v := range i.IndexBuffer {
		switch i.IndexBufferSize {
		case 1:
			writer.WriteU8(byte(i.IndexBuffer[v]))
		case 2:
			writer.WriteU16(uint16(i.IndexBuffer[v]))
		case 4:
			writer.WriteU32(i.IndexBuffer[v])
		default:
			// this should not be reached unless manual bad modification for the model was done
			panic(fmt.Errorf("unsupported index buffer size: %d", i.IndexBufferSize))
		}
	}

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

func (m *Matrix4x4) Encode(writer *Writer) {
	// we do not want to transpose the original one in case the user wanted to use it somewhere else
	// *after* encoding ,so we make a copy
	matrix := *m
	matrix.Transpose() // return it back to the original state before Decode
	for i := 0; i < 16; i++ {
		writer.WriteFloat(matrix[i/4][i%4])
	}
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

	if g.HasBindMatrix, err = reader.ReadBool(); err != nil {
		return err
	} else if g.HasBindMatrix {
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
		if err = g.SkinWeights[i].Decode(reader, g.SCWFile.Version, g.SCWFile.MinorVersion); err != nil {
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

func (g *Geometry) Encode(writer *Writer) {

	writer.WriteStringUTF(g.Name)
	writer.WriteStringUTF(g.Group)

	if g.SCWFile.Version <= 1 {
		g.IgnoredMatrix.Encode(writer)
	}

	writer.WriteU8(byte(len(g.Vertices)))

	for _, sourceArray := range g.Vertices {
		sourceArray.Encode(writer)
	}

	writer.WriteBool(g.HasBindMatrix)
	if g.HasBindMatrix {
		g.BindMatrix.Encode(writer)
	}

	writer.WriteU8(byte(len(g.Skins.Joints)))

	for i, joint := range g.Skins.Joints {
		writer.WriteStringUTF(joint)
		g.Skins.InverseBindMatrices[i].Encode(writer)
	}

	writer.WriteU32(uint32(len(g.SkinWeights)))
	for _, skinWeight := range g.SkinWeights {
		skinWeight.Encode(writer, g.SCWFile.Version, g.SCWFile.MinorVersion)
	}

	writer.WriteU8(byte(len(g.Materials)))

	for _, mat := range g.Materials {
		mat.Encode(writer)
	}

}
