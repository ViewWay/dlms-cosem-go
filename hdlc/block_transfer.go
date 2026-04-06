package hdlc

import (
	"errors"
	"fmt"
	"time"
)

// BlockTransferState represents the state of a block transfer
type BlockTransferState int

const (
	StateIdle BlockTransferState = iota
	StateTransmitting
	StateLastBlockReceived
	StateBlockReceived
	StateWaitingForNextBlock
	StateError
	StateTimeout
	StateComplete
)

// String returns the string representation of the state
func (s BlockTransferState) String() string {
	return [...]string{
		"IDLE",
		"TRANSMITTING",
		"LAST_BLOCK_RECEIVED",
		"BLOCK_RECEIVED",
		"WAITING_FOR_NEXT_BLOCK",
		"ERROR",
		"TIMEOUT",
		"COMPLETE",
	}[s]
}

// Block represents a single data block
type Block struct {
	BlockNumber int    // Block number (0-based)
	Data       []byte // Block data
	LastBlock  bool   // Is this the last block?
}

// BlockTransfer manages block transfer operations
type BlockTransfer struct {
	state               BlockTransferState
	currentBlockNumber   int
	receivedBlocks       []Block
	isComplete          bool
	totalExpectedLength  *int
	timeoutMs           int
}

// NewBlockTransfer creates a new block transfer manager
func NewBlockTransfer() *BlockTransfer {
	return &BlockTransfer{
		state:             StateIdle,
		currentBlockNumber: 0,
		receivedBlocks:     make([]Block, 0),
		isComplete:        false,
		timeoutMs:          5000, // 5 second default
	}
}

// Reset resets the transfer state
func (bt *BlockTransfer) Reset() {
	bt.state = StateIdle
	bt.currentBlockNumber = 0
	bt.receivedBlocks = bt.receivedBlocks[:0]
	bt.isComplete = false
	bt.totalExpectedLength = nil
}

// StartTransfer starts a new block transfer
func (bt *BlockTransfer) StartTransfer(expectedLength *int) error {
	if bt.state != StateIdle {
		return fmt.Errorf("cannot start transfer in state %s", bt.state)
	}

	bt.state = StateTransmitting
	bt.currentBlockNumber = 0
	bt.isComplete = false
	bt.totalExpectedLength = expectedLength

	return nil
}

// AddBlock adds a received block
func (bt *BlockTransfer) AddBlock(block Block) error {
	if bt.state != StateTransmitting && bt.state != StateWaitingForNextBlock {
		return fmt.Errorf("cannot add block in state %s", bt.state)
	}

	if block.BlockNumber != bt.currentBlockNumber {
		return fmt.Errorf("expected block %d, got %d", bt.currentBlockNumber, block.BlockNumber)
	}

	bt.receivedBlocks = append(bt.receivedBlocks, block)
	bt.currentBlockNumber++

	if block.LastBlock {
		bt.state = StateComplete
		bt.isComplete = true
	} else {
		bt.state = StateWaitingForNextBlock
	}

	return nil
}

// GetData returns reassembled data if transfer is complete
func (bt *BlockTransfer) GetData() ([]byte, error) {
	if !bt.isComplete {
		return nil, errors.New("transfer not complete")
	}

	data := make([]byte, 0, bt.getTotalSize())
	for _, block := range bt.receivedBlocks {
		data = append(data, block.Data...)
	}

	return data, nil
}

// GetBlock returns a specific block by block number
func (bt *BlockTransfer) GetBlock(blockNumber int) *Block {
	for i := range bt.receivedBlocks {
		if bt.receivedBlocks[i].BlockNumber == blockNumber {
			return &bt.receivedBlocks[i]
		}
	}
	return nil
}

// Reassemble reassembles all received blocks into complete data
func (bt *BlockTransfer) Reassemble() ([]byte, error) {
	if !bt.isComplete {
		return nil, errors.New("cannot reassemble incomplete block transfer")
	}

	data := make([]byte, 0, bt.getTotalSize())
	for _, block := range bt.receivedBlocks {
		data = append(data, block.Data...)
	}

	return data, nil
}

// getTotalSize calculates the total size of all received blocks
func (bt *BlockTransfer) getTotalSize() int {
	size := 0
	for _, block := range bt.receivedBlocks {
		size += len(block.Data)
	}
	return size
}

// GetTotalTransmitted returns the total number of transmitted blocks
func (bt *BlockTransfer) GetTotalTransmitted() int {
	return len(bt.receivedBlocks)
}

// GetProgress returns transfer progress information
func (bt *BlockTransfer) GetProgress() map[string]interface{} {
	return map[string]interface{}{
		"state":        bt.state.String(),
		"currentBlock": bt.currentBlockNumber,
		"totalBlocks":  len(bt.receivedBlocks),
		"isComplete":   bt.isComplete,
	}
}

// ValidateSequenceNumbers validates that block sequence numbers are correct
func (bt *BlockTransfer) ValidateSequenceNumbers() error {
	expectedNumber := 0
	for _, block := range bt.receivedBlocks {
		if block.BlockNumber != expectedNumber {
			return fmt.Errorf("block sequence error: expected %d, got %d",
				expectedNumber, block.BlockNumber)
		}
		expectedNumber++
	}

	return nil
}

// SetTimeout sets the timeout in milliseconds
func (bt *BlockTransfer) SetTimeout(timeoutMs int) {
	bt.timeoutMs = timeoutMs
}

// GetTimeout returns the timeout in milliseconds
func (bt *BlockTransfer) GetTimeout() int {
	return bt.timeoutMs
}

// WaitForCompletion waits for transfer to complete or timeout
func (bt *BlockTransfer) WaitForCompletion() error {
	if bt.isComplete {
		return nil
	}

	timeout := time.After(time.Duration(bt.timeoutMs) * time.Millisecond)
	<-timeout

	if !bt.isComplete {
		bt.state = StateTimeout
		return errors.New("transfer timed out")
	}

	return nil
}
