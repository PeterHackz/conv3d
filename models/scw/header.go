package scw

type Header struct {
	Version       uint16
	FrameRate     uint16
	FirstFrame    uint16
	LastFrame     uint16
	MaterialsFile string
	Unknown       int
}

func (h *Header) Tag() string {
	return "HEAD"
}

func (h *Header) Decode(reader *Reader) (err error) {
	h.Version, err = reader.ReadU16()
	if err != nil {
		return
	}

	h.FrameRate, err = reader.ReadU16()
	if err != nil {
		return
	}

	h.FirstFrame, err = reader.ReadU16()
	if err != nil {
		return
	}

	h.LastFrame, err = reader.ReadU16()
	if err != nil {
		return
	}

	h.MaterialsFile, err = reader.ReadUTF()
	if err != nil {
		return
	}

	if reader.SkipBytes >= 1 {
		var unk byte
		unk, err = reader.ReadU8()
		if err != nil {
			return
		}
		h.Unknown = int(unk)
	} else {
		h.Unknown = -1
	}
	return
}

func (h *Header) Encode(writer *Writer) {
	writer.WriteU16(h.Version)
	writer.WriteU16(h.FrameRate)
	writer.WriteU16(h.FirstFrame)
	writer.WriteU16(h.LastFrame)
	writer.WriteStringUTF(h.MaterialsFile)
	if h.Unknown != -1 {
		writer.WriteU8(byte(h.Unknown))
	}
}
