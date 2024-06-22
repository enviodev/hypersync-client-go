package types

type TraceSelection struct {
	From       []Address `json:"from"`
	To         []Address `json:"to"`
	Address    []Address `json:"address"`
	CallType   []string  `json:"call_type"`
	RewardType []string  `json:"reward_type"`
	Kind       []string  `json:"type"`
	SigHash    []SigHash `json:"sighash"`
}
