// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ptn

import (
	//"bytes"
	"context"
	//"errors"
	//"fmt"
	//"io/ioutil"
	//"runtime"
	//"sync"
	"time"

	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/hexutil"
	"github.com/palletone/go-palletone/common/rpc"
	"github.com/palletone/go-palletone/dag/state"
)

const (
	// defaultTraceTimeout is the amount of time a single transaction can execute
	// by default before being forcefully aborted.
	defaultTraceTimeout = 5 * time.Second

	// defaultTraceReexec is the number of blocks the tracer is willing to go back
	// and reexecute to produce missing historical state necessary to run a specific
	// trace.
	defaultTraceReexec = uint64(128)
)

// TraceConfig holds extra parameters to trace functions.
type TraceConfig struct {
	//*vm.LogConfig
	Tracer  *string
	Timeout *string
	Reexec  *uint64
}

// txTraceResult is the result of a single transaction trace.
type txTraceResult struct {
	Result interface{} `json:"result,omitempty"` // Trace results produced by the tracer
	Error  string      `json:"error,omitempty"`  // Trace failure produced by the tracer
}

// blockTraceTask represents a single block trace task when an entire chain is
// being traced.
type blockTraceTask struct {
	statedb *state.StateDB   // Intermediate state prepped for tracing
	rootref common.Hash      // Trie root reference held for this task
	results []*txTraceResult // Trace results procudes by the task
}

// blockTraceResult represets the results of tracing a single block when an entire
// chain is being traced.
type blockTraceResult struct {
	Block  hexutil.Uint64   `json:"block"`  // Block number corresponding to this trace
	Hash   common.Hash      `json:"hash"`   // Block hash corresponding to this trace
	Traces []*txTraceResult `json:"traces"` // Trace results produced by the task
}

// txTraceTask represents a single transaction trace task when an entire block
// is being traced.
type txTraceTask struct {
	statedb *state.StateDB // Intermediate state prepped for tracing
	index   int            // Transaction offset in the block
}

// TraceChain returns the structured logs created during the execution of EVM
// between two blocks (excluding start) and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceChain(ctx context.Context, start, end rpc.BlockNumber, config *TraceConfig) (*rpc.Subscription, error) {
	return &rpc.Subscription{}, nil
}

// TraceBlockByNumber returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceBlockByNumber(ctx context.Context, number rpc.BlockNumber, config *TraceConfig) ([]*txTraceResult, error) {
	/*
		// Fetch the block that we want to trace
		var block *types.Block

		switch number {
		case rpc.PendingBlockNumber:
			block = api.eth.miner.PendingBlock()
		case rpc.LatestBlockNumber:
			block = api.eth.blockchain.CurrentBlock()
		default:
			block = api.eth.blockchain.GetBlockByNumber(uint64(number))
		}
		// Trace the block if it was found
		if block == nil {
			return nil, fmt.Errorf("block #%d not found", number)
		}
		return api.traceBlock(ctx, block, config)
	*/
	return []*txTraceResult{}, nil
}

// TraceBlockByHash returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceBlockByHash(ctx context.Context, hash common.Hash, config *TraceConfig) ([]*txTraceResult, error) {
	/*
		block := api.eth.blockchain.GetBlockByHash(hash)
		if block == nil {
			return nil, fmt.Errorf("block #%x not found", hash)
		}
		return api.traceBlock(ctx, block, config)
	*/
	return []*txTraceResult{}, nil
}

// TraceBlock returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceBlock(ctx context.Context, blob []byte, config *TraceConfig) ([]*txTraceResult, error) {
	/*
		block := new(types.Block)
		if err := rlp.Decode(bytes.NewReader(blob), block); err != nil {
			return nil, fmt.Errorf("could not decode block: %v", err)
		}
		return api.traceBlock(ctx, block, config)
	*/
	return []*txTraceResult{}, nil
}

// TraceBlockFromFile returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceBlockFromFile(ctx context.Context, file string, config *TraceConfig) ([]*txTraceResult, error) {
	/*
		blob, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("could not read file: %v", err)
		}
		return api.TraceBlock(ctx, blob, config)
	*/
	return []*txTraceResult{}, nil
}

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceTransaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	/*
		// Retrieve the transaction and assemble its EVM context
		tx, blockHash, _, index := core.GetTransaction(api.eth.ChainDb(), hash)
		if tx == nil {
			return nil, fmt.Errorf("transaction %x not found", hash)
		}
		reexec := defaultTraceReexec
		if config != nil && config.Reexec != nil {
			reexec = *config.Reexec
		}
		msg, vmctx, statedb, err := api.computeTxEnv(blockHash, int(index), reexec)
		if err != nil {
			return nil, err
		}
		// Trace the transaction and return
		return api.traceTx(ctx, msg, vmctx, statedb, config)
	*/
	return nil, nil
}

/*
// computeTxEnv returns the execution environment of a certain transaction.
func (api *PrivateDebugAPI) computeTxEnv(blockHash common.Hash, txIndex int, reexec uint64) (core.Message, vm.Context, *state.StateDB, error) {

		// Create the parent state database
		block := api.eth.blockchain.GetBlockByHash(blockHash)
		if block == nil {
			return nil, vm.Context{}, nil, fmt.Errorf("block %x not found", blockHash)
		}
		parent := api.eth.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
		if parent == nil {
			return nil, vm.Context{}, nil, fmt.Errorf("parent %x not found", block.ParentHash())
		}
		statedb, err := api.computeStateDB(parent, reexec)
		if err != nil {
			return nil, vm.Context{}, nil, err
		}
		// Recompute transactions up to the target index.
		signer := types.MakeSigner(api.config, block.Number())

		for idx, tx := range block.Transactions() {
			// Assemble the transaction call message and return if the requested offset
			msg, _ := tx.AsMessage(signer)
			context := core.NewEVMContext(msg, block.Header(), api.eth.blockchain, nil)
			if idx == txIndex {
				return msg, context, statedb, nil
			}
			// Not yet the searched for transaction, execute on top of the current state
			vmenv := vm.NewEVM(context, statedb, api.config, vm.Config{})
			if _, _, _, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(tx.Gas())); err != nil {
				return nil, vm.Context{}, nil, fmt.Errorf("tx %x failed: %v", tx.Hash(), err)
			}
			// Ensure any modifications are committed to the state
			statedb.Finalise(true)
		}
		return nil, vm.Context{}, nil, fmt.Errorf("tx index %d out of range for block %x", txIndex, blockHash)


}
*/
