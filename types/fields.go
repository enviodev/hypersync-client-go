package types

type FieldSelection struct {
	Block       []string `json:"block,omitempty"`
	Transaction []string `json:"transaction,omitempty"`
	Log         []string `json:"log,omitempty"`
	Trace       []string `json:"trace,omitempty"`
}
