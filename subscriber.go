package hypersyncgo

import (
	"context"
	"math/big"
	"sync"
)

type Subscriber struct {
	ctx              context.Context
	client           *Client
	startBlockNumber *big.Int
	endBlockNumber   *big.Int
	mu               sync.Mutex
	ch               chan *Block
}

func NewBlockSubscriber(ctx context.Context, client *Client, startBlockNumber, endBlockNumber *big.Int) (*Subscriber, error) {
	toReturn := &Subscriber{
		ctx:              ctx,
		client:           client,
		startBlockNumber: startBlockNumber,
		endBlockNumber:   endBlockNumber,
		mu:               sync.Mutex{},
		ch:               make(chan *Block),
	}

	return toReturn, nil
}

func (s *Subscriber) Subscribe() error {

	return nil
}

func (s *Subscriber) UnSubscribe() error {
	return nil
}

func (s *Subscriber) Ch() <-chan *Block {
	return s.ch
}
