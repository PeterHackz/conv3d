package scw

type Camera3D struct {
	Name        string
	Yfov, Xfov  float32 // not sure...
	AspectRatio float32
	ZNear, ZFar float32
}

func (c *Camera3D) Tag() string {
	return "CAME"
}

func (c *Camera3D) Decode(reader *Reader) (err error) {
	if c.Name, err = reader.ReadUTF(); err != nil {
		return
	}

	if c.Yfov, err = reader.ReadFloat(); err != nil {
		return
	}

	if c.Xfov, err = reader.ReadFloat(); err != nil {
		return
	}

	if c.AspectRatio, err = reader.ReadFloat(); err != nil {
		return
	}

	if c.ZNear, err = reader.ReadFloat(); err != nil {
		return
	}

	if c.ZFar, err = reader.ReadFloat(); err != nil {
		return
	}
	return
}

func (c *Camera3D) Encode(writer *Writer) {
	writer.WriteStringUTF(c.Name)
	writer.WriteFloat(c.Yfov)
	writer.WriteFloat(c.Xfov)
	writer.WriteFloat(c.AspectRatio)
	writer.WriteFloat(c.ZNear)
	writer.WriteFloat(c.ZFar)
}
