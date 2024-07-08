package scw

import (
	"fmt"
	"math"
)

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

// reference: https://golangbyexample.com/comparing-floating-point-numbers-go/
const tolerance = 1e-9

func areFloatsEqual(a, b float32) bool {
	return math.Abs(float64(a-b)) <= tolerance
}

type Vector3 struct {
	X, Y, Z float32
}

func (v *Vector3) Equals(other *Vector3) bool {
	return areFloatsEqual(v.X, other.X) && areFloatsEqual(v.Y, other.Y) && areFloatsEqual(v.Z, other.Z)
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

func (q *Quaternion) Equals(other *Quaternion) bool {
	return q.Vector3.Equals(&other.Vector3) && areFloatsEqual(q.W, other.W)
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

func (f *KeyFrame) Encode(writer *Writer, u8 uint8, v58 uint16, Frames []KeyFrame) {
	writer.WriteU16(f.ID)

	_ = Frames

	v61 := int(u8)
	if v58 == 0 {
		v61 = -1
	}
	if v58 == 0 || (v61&1) != 0 {
		writer.WriteI16(int16(f.Rotation.X / 0.000030758))
		writer.WriteI16(int16(f.Rotation.Y / 0.000030758))
		writer.WriteI16(int16(f.Rotation.Z / 0.000030758))
		writer.WriteI16(int16(f.Rotation.W / 0.000030758))
	}

	if v58 == 0 || (v61&2) != 0 {
		writer.WriteFloat(f.Translation.X)
	}

	if v58 == 0 || (v61&4) != 0 {
		writer.WriteFloat(f.Translation.Y)
	}

	if v58 == 0 || (v61&8) != 0 {
		writer.WriteFloat(f.Translation.Z)
	}

	if v58 == 0 || (v61&0x10) != 0 {
		writer.WriteFloat(f.Scale.X)
	}

	if v58 == 0 || (v61&0x20) != 0 {
		writer.WriteFloat(f.Scale.Y)
	}

	if v58 == 0 || (v61&0x40) != 0 {
		writer.WriteFloat(f.Scale.Z)
	}

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

func (i *InstanceMaterial) Encode(writer *Writer) {
	writer.WriteStringUTF(i.Name)
	writer.WriteStringUTF(i.Target)
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

func (n *NodeInstance) Encode(writer *Writer) {
	writer.WriteStringChars(n.Type)

	writer.WriteStringUTF(n.Target)

	switch n.Type {
	case "GEOM", "CONT":
		writer.WriteU16(uint16(len(n.Materials)))
		for _, mat := range n.Materials {
			mat.Encode(writer)
		}
	case "LIGH":
		panic("LIGH is not supported yet")
	case "CAME":
		writer.WriteStringUTF(n.CameraTarget)
	default:
		panic(fmt.Errorf("invalid or unsupported instance material type: %s", n.Type))
	}
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

func computeFrameFlags(Frames []KeyFrame) byte {
	flags := byte(0)

	// check if all frames matches the properties (rotation, ...)  of the first frame
	rotation := true
	translationX := true
	translationY := true
	translationZ := true
	scaleX := true
	scaleY := true
	scaleZ := true

	firstFrame := Frames[0]

	for _, frame := range Frames[1:] {
		if rotation && !frame.Rotation.Equals(&firstFrame.Rotation) {
			rotation = false
		}

		if translationX && !areFloatsEqual(frame.Translation.X, firstFrame.Translation.X) {
			translationX = false
		}

		if translationY && !areFloatsEqual(frame.Translation.Y, firstFrame.Translation.Y) {
			translationY = false
		}

		if translationZ && !areFloatsEqual(frame.Translation.Z, firstFrame.Translation.Z) {
			translationZ = false
		}

		if scaleX && !areFloatsEqual(frame.Scale.X, firstFrame.Scale.X) {
			scaleX = false
		}

		if scaleY && !areFloatsEqual(frame.Scale.Y, firstFrame.Scale.Y) {
			scaleY = false
		}

		if scaleZ && !areFloatsEqual(frame.Scale.Z, firstFrame.Scale.Z) {
			scaleZ = false
		}
	}

	if rotation {
		flags |= 1 << 0
	}

	if translationX {
		flags |= 1 << 1
	}

	if translationY {
		flags |= 1 << 2
	}

	if translationZ {
		flags |= 1 << 3
	}

	if scaleX {
		flags |= 1 << 4
	}

	if scaleY {
		flags |= 1 << 5
	}

	if scaleZ {
		flags |= 1 << 6
	}

	return flags
}

func (n *Node) Encode(writer *Writer) {
	writer.WriteStringUTF(n.Name)
	writer.WriteStringUTF(n.ParentName)

	writer.WriteU16(uint16(len(n.Instances)))

	for _, instance := range n.Instances {
		instance.Encode(writer)
	}

	writer.WriteU16(uint16(len(n.Frames)))

	if len(n.Frames) > 0 {
		n.FramesFlags = 0

		if len(n.Frames) > 1 {
			n.FramesFlags = computeFrameFlags(n.Frames)
		}

		writer.WriteU8(n.FramesFlags)

		for i, frame := range n.Frames {
			frame.Encode(writer, n.FramesFlags, uint16(i), n.Frames)
		}
	}
}
