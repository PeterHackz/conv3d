package scw

import (
	"bytes"
	"encoding/binary"
	"math"
)

type Writer struct {
	bytes.Buffer
}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) WriteU8(b byte) {
	_ = w.WriteByte(b)
}

func (w *Writer) WriteU16(value uint16) {
	var data [2]byte
	binary.BigEndian.PutUint16(data[:], value)
	_, _ = w.Write(data[:])
}

func (w *Writer) WriteU32(value uint32) {
	var data [4]byte
	binary.BigEndian.PutUint32(data[:], value)
	_, _ = w.Write(data[:])
}

func (w *Writer) WriteBool(value bool) {
	if value {
		w.WriteU8(1)
	} else {
		w.WriteU8(0)
	}
}

func (w *Writer) WriteI16(value int16) {
	w.WriteU16(uint16(value))
}

func (w *Writer) WriteFloat(value float32) {
	bits := math.Float32bits(value)
	w.WriteU32(bits)
}

func (w *Writer) WriteStringUTF(str string) {
	w.WriteU16(uint16(len(str)))
	w.WriteStringChars(str)
}

func (w *Writer) WriteStringChars(str string) {
	_, _ = w.WriteString(str)
}

func (w *Writer) WriteBytes(bytes []byte) {
	_, _ = w.Write(bytes)
}
