package scw

import "hash/crc32"

type Sc3dProperty interface {
	Tag() string
	Encode(writer *Writer)
}

func EncodeSc3dProperty(prop Sc3dProperty, writer *Writer) {
	w := NewWriter()
	w.WriteStringChars(prop.Tag())
	prop.Encode(w)
	bytes := w.Bytes()

	writer.WriteU32(uint32(len(bytes) - len(prop.Tag())))
	writer.WriteBytes(bytes)

	checksum := crc32.ChecksumIEEE(bytes)
	writer.WriteU32(checksum)
}
