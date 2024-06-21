package types

type TransactionSelection struct {
	From            []Address `json:"from"`
	To              []Address `json:"to"`
	SigHash         []SigHash `json:"sighash"`
	Status          *uint8    `json:"status"`
	Kind            []uint8   `json:"type"`
	ContractAddress []Address `json:"contract_address"`
}
