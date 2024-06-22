package types

type LogSelection struct {
	Address []Address        `json:"address"`
	Topics  [4][]LogArgument `json:"topics"`
}
