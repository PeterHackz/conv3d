package scw

// world3d end?

// Wend stop tag for scw parser
type Wend struct{}

func (w *Wend) Tag() string {
	return "WEND"
}

func (w *Wend) Encode(writer *Writer) {}
