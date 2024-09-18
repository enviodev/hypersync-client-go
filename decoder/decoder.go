package decoder

type Decoder struct{}

func New() (*Decoder, error) {
	return &Decoder{}, nil
}
