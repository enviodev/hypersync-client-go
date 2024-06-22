package types

type DataResponse struct {
	Blocks       []Block       `json:"blocks,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
}
