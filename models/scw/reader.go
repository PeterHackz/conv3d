package scw

import (
	"encoding/binary"
	"io"
	"math"
)

type Reader struct {
	data      []byte
	offset    int
	SkipBytes int
}

func NewReader(data []byte) *Reader {
	return &Reader{
		data:   data,
		offset: 0,
	}
}

func (r *Reader) hasData(n int) bool {
	return r.offset+n <= len(r.data)
}

func (r *Reader) Read(n int) ([]byte, error) {
	r.SkipBytes -= n
	if !r.hasData(n) {
		return nil, io.ErrShortBuffer
	}
	r.offset += n
	return r.data[r.offset-n : r.offset], nil
}

func (r *Reader) ReadUTF() (string, error) {
	length, err := r.ReadU16()
	if err != nil {
		return "", err
	}
	return r.ReadUTFWithLength(length)
}

func (r *Reader) ReadUTFWithLength(length uint16) (string, error) {
	if length == 0 {
		return "", nil
	}
	data, err := r.Read(int(length))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *Reader) ReadU32() (uint32, error) {
	data, err := r.Read(4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(data), nil
}

func (r *Reader) ReadU16() (uint16, error) {
	data, err := r.Read(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(data), nil
}

func (r *Reader) ReadU8() (uint8, error) {
	data, err := r.Read(1)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

func (r *Reader) Seek(n int) {
	r.offset += n
	r.SkipBytes -= n
}

func (r *Reader) ReadBool() (bool, error) {
	data, err := r.ReadU8()
	return data != 0, err
}

func (r *Reader) ReadFloat() (float32, error) {
	data, err := r.ReadU32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(data), nil
}

func (r *Reader) ReadI16() (int16, error) {
	data, err := r.Read(2)
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(data)), nil
}
