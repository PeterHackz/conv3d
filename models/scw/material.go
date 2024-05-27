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

type Variable struct {
	Texture2D string
	Color     RGBA
}

func (v *Variable) Decode(reader *Reader) error {
	if hasVariable, err := reader.ReadBool(); err != nil {
		return err
	} else if hasVariable {
		if v.Texture2D, err = reader.ReadUTF(); err != nil {
			return err
		}
	} else if err = v.Color.Decode(reader); err != nil {
		return err
	}
	return nil
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
