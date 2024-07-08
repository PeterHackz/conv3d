package scw

type Scene struct {
	Nodes []Node
}

func (s *Scene) Tag() string {
	return "NODE"
}

func (s *Scene) Encode(writer *Writer) {
	writer.WriteU16(uint16(len(s.Nodes)))
	for _, node := range s.Nodes {
		node.Encode(writer)
	}
}
