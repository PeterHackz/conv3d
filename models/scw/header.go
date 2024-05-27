package scw

type Header struct {
	Version       uint16
	FrameRate     uint16
	FirstFrame    uint16
	LastFrame     uint16
	MaterialsFile string
	Unknown       int
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
