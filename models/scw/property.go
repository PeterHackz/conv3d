package scw

import "hash/crc32"

type SC3DProperty interface {
	Tag() string
	Encode(writer *Writer)
}

func EncodeSc3dProperty(prop SC3DProperty, writer *Writer) {
	w := NewWriter()
	w.WriteStringChars(prop.Tag())
	prop.Encode(w)
	bytes := w.Bytes()

	writer.WriteU32(uint32(len(bytes) - len(prop.Tag())))
	writer.WriteBytes(bytes)

	checksum := crc32.ChecksumIEEE(bytes)
	writer.WriteU32(checksum)
}
