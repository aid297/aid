package reflection

type (
	Attributer interface {
		Register(ref *Reflection)
	}

	AttrSerializeFormat struct{ format string }
)

func SerializeFormat(format string) AttrSerializeFormat {
	return AttrSerializeFormat{format: format}
}

func (my AttrSerializeFormat) Register(ref *Reflection) { ref.serializeFormat = my.format }
