package main

import (
	"errors"
	"fmt"
	"io"
	// "strings"
)

// TODO: Replace with proper JSON serialization? Originally was written to be quick&dirty for maximum perf.

func genEthCall(w io.Writer, s State) error {  
	// We eth_call the block before the call actually happened to avoid collision reverts
	to, from, input, block := s.RandomCall()
	var err error
	if to != "" {
		_, err = fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_call","params":[{"to":%q,"from":%q,"data":%q},"0x%x"]}`+"\n", s.ID(), to, from, input, block-1)
	} else {
		_, err = fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_call","params":[{"from":%q,"data":%q},"0x%x"]}`+"\n", s.ID(), from, input, block-1)
	}
	return err
}

func genEthGetTransactionReceipt(w io.Writer, s State) error {
	txID := s.RandomTransaction()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getTransactionReceipt","params":["%s"]}`+"\n", s.ID(), txID)
	return err
}

func genEthGetBalance(w io.Writer, s State) error {
	r := s.RandInt64()
	// TODO: ~half of the block numbers are further from head
	blockNum := s.CurrentBlock() - uint64(r%5) // Within the last ~minute
	addr := s.RandomAddress()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getBalance","params":["%s","%#x"]}`+"\n", s.ID(), addr, blockNum)
	return err
}

func genEthGetBlockByNumber(w io.Writer, s State) error {
	r := s.RandInt64()
	// TODO: ~half of the block numbers are further from head
	blockNum := s.CurrentBlock() - uint64(r%5) // Within the last ~minute
	full := "true"
	if r%2 >= 0 {
		full = "false"
	}

	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getBlockByNumber","params":["0x%x",%s]}`+"\n", s.ID(), blockNum, full)
	return err
}

func genEthGetBlockByHash(w io.Writer, s State) error {
	blockHash := s.BlockHash()
	full := "true"

	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getBlockByHash","params":["%v",%s]}`+"\n", s.ID(), blockHash, full)
	return err
}

func genEthEstimateGas(w io.Writer, s State) error {  
	// We eth_call the block before the call actually happened to avoid collision reverts
	to, from, input, block := s.RandomCall()
	var err error
	if to != "" {
		_, err = fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_estimateGas","params":[{"to":%q,"from":%q,"data":%q},"0x%x"]}`+"\n", s.ID(), to, from, input, block-1)
	} else {
		_, err = fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_estimateGas","params":[{"from":%q,"data":%q},"0x%x"]}`+"\n", s.ID(), from, input, block-1)
	}
	return err
}



func genEthGetTransactionCount(w io.Writer, s State) error {
	r := s.RandInt64()
	// TODO: ~half of the block numbers are further from head
	blockNum := s.CurrentBlock() - uint64(r%5) // Within the last ~minute
	addr := s.RandomAddress()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getTransactionCount","params":["%s","%#x"]}`+"\n", s.ID(), addr, blockNum)
	return err
}

func genEthChainId(w io.Writer, s State) error {
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_chainId"}`+"\n", s.ID())
	return err
}

// Not useful for current tests
// func genEthSyncing(w io.Writer, s State) error {
// 	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_syncing"}`+"\n", s.ID())
// 	return err
// } 

// func genNetVersion(w io.Writer, s State) error {
// 	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"net_version"}`+"\n", s.ID())
// 	return err
// }

// func genWeb3ClientVersion(w io.Writer, s State) error {
// 	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"web3_clientVersion"}`+"\n", s.ID())
// 	return err
// }

func genEthGasPrice(w io.Writer, s State) error {
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_GasPrice"}`+"\n", s.ID())
	return err
}

// func genEthBlockNumber(w io.Writer, s State) error {
// 	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_blockNumber"}`+"\n", s.ID())
// 	return err
// }

func genEthGetTransactionByHash(w io.Writer, s State) error {
	txID := s.RandomTransaction()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getTransactionByHash","params":["%s"]}`+"\n", s.ID(), txID)
	return err
}

func genEthGetLogs(w io.Writer, s State) error {
	r := s.RandInt64()
	// TODO: Favour latest/recent block on a curve
	fromBlock := s.CurrentBlock() - uint64(r%5000) // Pick a block within the last ~day
	toBlock := s.CurrentBlock() - uint64(r%5)      // Within the last ~minute
	address, topics := s.RandomContract()

	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getLogs","params":[{"fromBlock":"0x%x","toBlock":"0x%x","address":"%s","topics":%v}]}`+"\n", s.ID(), fromBlock, toBlock, address, topics)
	return err
}

func genEthGetCode(w io.Writer, s State) error {
	r := s.RandInt64()
	// TODO: ~half of the block numbers are further from head
	blockNum := s.CurrentBlock() - uint64(r%5) // Within the last ~minute
	addr, _ := s.RandomContract()
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"eth_getCode","params":["%s","%#x"]}`+"\n", s.ID(), addr, blockNum)
	return err
}

func genBorGetAuthor(w io.Writer, s State) error {
	r := s.RandInt64()
	// TODO: ~half of the block numbers are further from head
	blockNum := s.CurrentBlock() - uint64(r%5) // Within the last ~minute

	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"bor_getAuthor","params":["0x%x"]}`+"\n", s.ID(), blockNum)
	return err
}

func genBorGetRootHash(w io.Writer, s State) error {
	r := s.RandInt64()

	fromBlock := s.CurrentBlock() - uint64(r%5000) // Pick a block within the last ~day
	toBlock := s.CurrentBlock() - uint64(r%5)      // Within the last ~minute

	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"bor_getRootHash","params":["%v","%v"]}`+"\n", s.ID(), fromBlock, toBlock)
	return err
}

func genBorGetSnapshot(w io.Writer, s State) error {
	r := s.RandInt64()

	block := s.CurrentBlock() - uint64(r%5000) // Pick a block within the last ~day

	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"bor_getSnapshot","params":["0x%x"]}`+"\n", s.ID(), block)
	return err
}


func genBorGetSignersAtHash(w io.Writer, s State) error {
	r := s.RandInt64()

	block := s.CurrentBlock() - uint64(r%5000) // Pick a block within the last ~day

	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"bor_getSignersAtHash","params":["0x%x"]}`+"\n", s.ID(), block)
	return err
}

func genBorGetCurrentValidators(w io.Writer, s State) error {
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"bor_getCurrentValidators"}`+"\n", s.ID())
	return err
}

func genBorGetCurrentProposer(w io.Writer, s State) error {
	_, err := fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"method":"bor_getCurrentProposer"}`+"\n", s.ID())
	return err
}

func installDefaults(gen *generator, methods map[string]int64) error {
	// Top queries by weight, pulled from a 5000 Infura query sample on Dec 2019.
	//     3 "eth_accounts"
	//     4 "eth_getStorageAt"
	//     4 "eth_syncing"
	//     7 "net_peerCount"
	//    12 "net_listening"
	//    14 "eth_gasPrice"
	//    16 "eth_sendRawTransaction"
	//    25 "net_version"
	//    30 "eth_getTransactionByBlockNumberAndIndex"
	//    38 "eth_getBlockByHash"
	//    45 "eth_estimateGas"
	//    88 "eth_getCode"
	//   252 "eth_getLogs"
	//   255 "eth_getTransactionByHash"
	//   333 "eth_blockNumber"
	//   390 "eth_getTransactionCount"
	//   399 "eth_getBlockByNumber"
	//   545 "eth_getBalance"
	//   607 "eth_getTransactionReceipt"
	//  1928 "eth_call"

	rpcMethod := map[string]func(io.Writer, State) error{
		"eth_call":                  genEthCall,
		"eth_getTransactionReceipt": genEthGetTransactionReceipt,
		"eth_getBalance":            genEthGetBalance,
		"eth_getBlockByNumber":      genEthGetBlockByNumber,
		"eth_getTransactionCount":   genEthGetTransactionCount,
		// "eth_blockNumber":           genEthBlockNumber,
		"eth_getTransactionByHash":  genEthGetTransactionByHash,
		"eth_estimateGas":           genEthEstimateGas,
		"eth_getLogs":               genEthGetLogs,
		"eth_getCode":               genEthGetCode,
		"eth_chainId":               genEthChainId,
		"eth_getBlockByHash":        genEthGetBlockByHash,
		"eth_gasPrice":              genEthGasPrice,
		// "eth_syncing":               genEthSyncing,
		// "net_version":               genNetVersion,
		// "web3_clientVersion":        genWeb3ClientVersion,
		"bor_getAuthor":             genBorGetAuthor,
		"bor_getRootHash":           genBorGetRootHash,
		"bor_getSnapshot":           genBorGetSnapshot,
		"bor_getSignersAtHash":      genBorGetSignersAtHash,
		"bor_getCurrentValidators":  genBorGetCurrentValidators,
		"bor_getCurrentProposer":    genBorGetCurrentProposer,

	}

	for method, weight := range methods {
		if weight == 0 {
			continue
		}
		if _, err := rpcMethod[method]; err == false {
			return errors.New(method + " is not supported")
		}
		gen.Add(RandomQuery{
			Method:   method,
			Weight:   weight,
			Generate: rpcMethod[method],
		})
	}

	return nil
}
