package errorshs

import "errors"

var (
	ErrContractNotFound = errors.New("contract not found")
	ErrWorkerCompleted  = errors.New("worker completed")
)
