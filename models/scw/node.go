package scw

import "fmt"

type Node struct {
	SCWFile          *File `json:"-"`
	Name, ParentName string
	Instances        []NodeInstance
	Frames           []KeyFrame
	FramesFlags      byte
}

type KeyFrame struct {
	ID                 uint16
	Rotation           Quaternion
	Translation, Scale Vector3
}

type Vector3 struct {
	X, Y, Z float32
}

// Quaternion q= w + xi + yj + zk
//
// w is the scalar component
//
// x, y and z are the vector components
type Quaternion struct {
	Vector3
	W float32
}

func (f *KeyFrame) Decode(reader *Reader, u8 uint8, v58 uint16, Frames []KeyFrame) (err error) {
	var val int16
	if f.ID, err = reader.ReadU16(); err != nil {
		return
	}

	v61 := int(u8)
	if v58 == 0 {
		v61 = -1
	}
	if v58 == 0 || (v61&1) != 0 {
		if val, err = reader.ReadI16(); err != nil {
			return
		}
		f.Rotation.X = float32(val) * 0.000030758

		if val, err = reader.ReadI16(); err != nil {
			return
		}
		f.Rotation.Y = float32(val) * 0.000030758

		if val, err = reader.ReadI16(); err != nil {
			return
		}
		f.Rotation.Z = float32(val) * 0.000030758

		if val, err = reader.ReadI16(); err != nil {
			return
		}
		f.Rotation.W = float32(val) * 0.000030758
	} else {
		f.Rotation = Frames[0].Rotation
	}

	if v58 == 0 || (v61&2) != 0 {
		if f.Translation.X, err = reader.ReadFloat(); err != nil {
			return
		}
	} else {
		f.Translation.X = Frames[0].Translation.X
	}

	if v58 == 0 || (v61&4) != 0 {
		if f.Translation.Y, err = reader.ReadFloat(); err != nil {
			return
		}
	} else {
		f.Translation.Y = Frames[0].Translation.Y
	}

	if v58 == 0 || (v61&8) != 0 {
		if f.Translation.Z, err = reader.ReadFloat(); err != nil {
			return
		}
	} else {
		f.Translation.Z = Frames[0].Translation.Z
	}

	if v58 == 0 || (v61&0x10) != 0 {
		if f.Scale.X, err = reader.ReadFloat(); err != nil {
			return
		}
	} else {
		f.Scale.X = Frames[0].Scale.X
	}

	if v58 == 0 || (v61&0x20) != 0 {
		if f.Scale.Y, err = reader.ReadFloat(); err != nil {
			return
		}
	} else {
		f.Scale.Y = Frames[0].Scale.Y
	}

	if v58 == 0 || (v61&0x40) != 0 {
		if f.Scale.Z, err = reader.ReadFloat(); err != nil {
			return
		}
	} else {
		f.Scale.Z = Frames[0].Scale.Z
	}

	return
}

type NodeInstance struct {
	Type, Target string
	CameraTarget string
	Materials    []InstanceMaterial
}

type InstanceMaterial struct {
	Name, Target string
}

func (i *InstanceMaterial) Decode(reader *Reader) (err error) {
	if i.Name, err = reader.ReadUTF(); err != nil {
		return
	}

	if i.Target, err = reader.ReadUTF(); err != nil {
		return
	}
	return
}

func (n *NodeInstance) Decode(reader *Reader) (err error) {
	if n.Type, err = reader.ReadUTFWithLength(4); err != nil {
		return
	}

	if n.Target, err = reader.ReadUTF(); err != nil {
		return
	}

	switch n.Type {
	case "GEOM", "CONT":
		var count uint16
		if count, err = reader.ReadU16(); err != nil {
			return
		}
		n.Materials = make([]InstanceMaterial, count)
		for i := range count {
			if err = n.Materials[i].Decode(reader); err != nil {
				return
			}
		}
	case "LIGH":
		panic("LIGH is not supported yet")
	case "CAME":
		if n.CameraTarget, err = reader.ReadUTF(); err != nil {
			return
		}
	default:
		return fmt.Errorf("invalid or unsupported instance material type: %s", n.Type)
	}
	return
}

func (n *Node) Decode(reader *Reader) (err error) {

	if n.Name, err = reader.ReadUTF(); err != nil {
		return
	}

	if n.ParentName, err = reader.ReadUTF(); err != nil {
		return
	}

	// v26 = SCWFileInput::readU16(a2);
	var count uint16
	if count, err = reader.ReadU16(); err != nil {
		return
	}

	n.Instances = make([]NodeInstance, count)

	for i := range count {
		if err = n.Instances[i].Decode(reader); err != nil {
			return
		}
	}

	if count, err = reader.ReadU16(); err != nil {
		return
	}

	n.Frames = make([]KeyFrame, count)

	if count > 0 {
		var u8 uint8
		if u8, err = reader.ReadU8(); err != nil {
			return
		}

		n.FramesFlags = u8

		for i := range count {
			if err = n.Frames[i].Decode(reader, u8, i, n.Frames); err != nil {
				return
			}
		}
	}

	return
}
