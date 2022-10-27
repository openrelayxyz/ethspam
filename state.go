package main

import (
	"context"
	"errors"
	"math/rand"
	"sync/atomic"

	"github.com/INFURA/go-ethlibs/eth"
	"github.com/INFURA/go-ethlibs/node"
)

var errEmptyBlock = errors.New("the sampled block is empty")

type State interface {
	RandInt64() int64
	ID() int64
	CurrentBlock() uint64
	RandomContract() (addr string, topics []string)
	RandomAddress() string
	RandomTransaction() string
	RandomCall() (to, from, input string, block uint64)
	BlockHash() string
}

type idGenerator struct {
	id int64
}

func (gen *idGenerator) Next() int64 {
	return atomic.AddInt64(&gen.id, 1)
}

// liveState implements State but it seeds the state dataset from live sources
// (Etherscan, etc)
type liveState struct {
	idGen   *idGenerator
	randSrc rand.Source

	currentBlock uint64
	transactions []eth.Transaction
	blockHash *eth.Hash
}

func (s *liveState) ID() int64 {
	return s.idGen.Next()
}

func (s *liveState) CurrentBlock() uint64 {
	return s.currentBlock
}

func (s *liveState) RandInt64() int64 {
	return s.randSrc.Int63()
}

func (s *liveState) BlockHash() string {
	return string(*s.blockHash)
}

func (s *liveState) RandomTransaction() string {
	if len(s.transactions) == 0 {
		return ""
	}
	idx := int(s.randSrc.Int63()) % len(s.transactions)
	return s.transactions[idx].Hash.String()
}

func (s *liveState) RandomAddress() string {
	if len(s.transactions) == 0 {
		return ""
	}
	idx := int(s.randSrc.Int63()) % len(s.transactions)
	return s.transactions[idx].From.String()
}

func (s *liveState) RandomCall() (to, from, input string, block uint64) {
	if len(s.transactions) == 0 {
		return
	}
	tx := s.transactions[int(s.randSrc.Int63())%len(s.transactions)]
	if tx.To != nil {
		to = tx.To.String()
	}
	from = tx.From.String()
	input = tx.Input.String()
	block = tx.BlockNumber.UInt64()
	return
}

func (s *liveState) RandomContract() (addr string, topics []string) {

	addresses, _ := fuzzAddress()
	topics, _ = fuzzTopics()

	return addresses, topics
}

type stateProducer struct {
	client node.Client
}

func (p *stateProducer) Refresh(oldState *liveState) (*liveState, error) {
	if oldState == nil {
		return nil, errors.New("must provide old state to refresh")
	}

	b, err := p.client.BlockByNumberOrTag(context.Background(), *(eth.MustBlockNumberOrTag("latest")), true)
	if err != nil {
		return nil, err
	}
	// Short circuit if the sampled block is empty
	if len(b.Transactions) == 0 {
		return nil, errEmptyBlock
	}
	// txs will grow to the maximum contract transaction list size we'll see in a block, and the higher-indexed ones will stick around longer
	txs := oldState.transactions
	for i, tx := range b.Transactions {
		if tx.Transaction.Value.Int64() > 0 {
			// Only take 0-value transactions, hopefully these are all contract calls.
			continue
		}
		if len(oldState.transactions) < 50 || i > len(txs) {
			txs = append(txs, tx.Transaction)
			continue
		}
		// Keep some old transactions randomly
		if oldState.RandInt64()%6 > 2 {
			txs[i] = tx.Transaction
		}
	}

	// var block eth.Block
	// *block = b

	state := liveState{
		idGen:   oldState.idGen,
		randSrc: oldState.randSrc,

		currentBlock: b.Number.UInt64(),
		transactions: txs,
		blockHash:  b.Hash,
	}

	return &state, nil
}
