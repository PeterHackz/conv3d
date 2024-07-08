package scw

type Material struct {
	SCWFile    *File `json:"-"`
	Name       string
	ShaderFile string
	BlendMode  byte
	Variables  struct {
		Ambient                                    Variable
		Diffuse                                    Variable
		Specular                                   Variable
		StencilTex2D                               string
		NormalTex2D                                string
		Colorize                                   Variable
		Emission                                   Variable
		OpacityTex2D                               string
		Opacity, Unk                               float32
		LightmapTex2D, LightmapSpecularTex2D, Unk2 string
	}
	ShaderConfig       uint32
	StencilScaleOffset [4]float32
}

func (m *Material) Tag() string {
	return "MATE"
}

type Variable struct {
	UseText2D bool
	Texture2D string
	Color     RGBA
}

func (v *Variable) Decode(reader *Reader) error {
	var err error
	if v.UseText2D, err = reader.ReadBool(); err != nil {
		return err
	} else if v.UseText2D {
		if v.Texture2D, err = reader.ReadUTF(); err != nil {
			return err
		}
	} else if err = v.Color.Decode(reader); err != nil {
		return err
	}
	return nil
}

func (v *Variable) Encode(writer *Writer) {
	writer.WriteBool(v.UseText2D)
	if v.UseText2D {
		writer.WriteStringUTF(v.Texture2D)
	} else {
		v.Color.Encode(writer)
	}
}

type RGBA [4]byte

func (r *RGBA) Decode(reader *Reader) (err error) {
	for i := range 4 {
		if r[i], err = reader.ReadU8(); err != nil {
			return
		}
	}
	return
}

func (r *RGBA) Encode(writer *Writer) {
	writer.WriteBytes(r[:])
}

func (m *Material) Decode(reader *Reader) (err error) {
	m.Name, err = reader.ReadUTF()
	if err != nil {
		return err
	}

	m.ShaderFile, err = reader.ReadUTF()
	if err != nil {
		return err
	}

	m.BlendMode, err = reader.ReadU8() // Material::bind: Stage::bindBlendMode(Stage::sm_pInstance, this->blendMode << 7);
	if err != nil {
		return err
	}

	if err = m.Variables.Ambient.Decode(reader); err != nil {
		return
	}

	if err = m.Variables.Diffuse.Decode(reader); err != nil {
		return
	}

	if err = m.Variables.Specular.Decode(reader); err != nil {
		return
	}

	if m.Variables.StencilTex2D, err = reader.ReadUTF(); err != nil {
		return
	}

	if m.Variables.NormalTex2D, err = reader.ReadUTF(); err != nil {
		return
	}

	if err = m.Variables.Colorize.Decode(reader); err != nil {
		return
	}

	if err = m.Variables.Emission.Decode(reader); err != nil {
		return
	}

	if m.Variables.OpacityTex2D, err = reader.ReadUTF(); err != nil {
		return
	}

	if m.Variables.Opacity, err = reader.ReadFloat(); err != nil {
		return
	}

	if m.Variables.Unk, err = reader.ReadFloat(); err != nil {
		return
	}

	if m.Variables.LightmapTex2D, err = reader.ReadUTF(); err != nil {
		return
	}

	if m.Variables.LightmapSpecularTex2D, err = reader.ReadUTF(); err != nil {
		return
	}

	if m.SCWFile.Version >= 2 {
		if m.Variables.Unk2, err = reader.ReadUTF(); err != nil {
			return
		}
	}

	if m.ShaderConfig, err = reader.ReadU32(); err != nil {
		return
	}

	if m.ShaderConfig&0x8000 != 0 {
		for i := range 4 {
			if m.StencilScaleOffset[i], err = reader.ReadFloat(); err != nil {
				return err
			}
		}
	}

	return
}

func (m *Material) Encode(writer *Writer) {
	writer.WriteStringUTF(m.Name)
	writer.WriteStringUTF(m.ShaderFile)
	writer.WriteU8(m.BlendMode)

	m.Variables.Ambient.Encode(writer)
	m.Variables.Diffuse.Encode(writer)
	m.Variables.Specular.Encode(writer)

	writer.WriteStringUTF(m.Variables.StencilTex2D)
	writer.WriteStringUTF(m.Variables.NormalTex2D)

	m.Variables.Colorize.Encode(writer)
	m.Variables.Emission.Encode(writer)

	writer.WriteStringUTF(m.Variables.OpacityTex2D)
	writer.WriteFloat(m.Variables.Opacity)

	writer.WriteFloat(m.Variables.Unk)

	writer.WriteStringUTF(m.Variables.LightmapTex2D)
	writer.WriteStringUTF(m.Variables.LightmapSpecularTex2D)

	if m.SCWFile.Version >= 2 {
		writer.WriteStringUTF(m.Variables.Unk2)
	}

	writer.WriteU32(m.ShaderConfig)

	if m.ShaderConfig&0x8000 != 0 {
		for i := range 4 {
			writer.WriteFloat(m.StencilScaleOffset[i])
		}
	}

}
