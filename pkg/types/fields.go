package types

type FieldSelection struct {
	Block       []string `json:"block"`
	Transaction []string `json:"transaction"`
	Log         []string `json:"log"`
	Trace       []string `json:"trace"`
}
