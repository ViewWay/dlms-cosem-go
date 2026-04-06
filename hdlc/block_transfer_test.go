package hdlc

import (
	"testing"
)

func TestNewBlockTransfer(t *testing.T) {
	bt := NewBlockTransfer()
	if bt.state != StateIdle {
		t.Errorf("Expected state IDLE, got %v", bt.state)
	}
	if bt.currentBlockNumber != 0 {
		t.Errorf("Expected currentBlockNumber 0, got %d", bt.currentBlockNumber)
	}
	if len(bt.receivedBlocks) != 0 {
		t.Errorf("Expected 0 received blocks, got %d", len(bt.receivedBlocks))
	}
	if bt.isComplete {
		t.Error("Expected isComplete to be false")
	}
}

func TestBlockTransfer_StartTransfer(t *testing.T) {
	bt := NewBlockTransfer()
	expectedLength := 1000

	err := bt.StartTransfer(&expectedLength)
	if err != nil {
		t.Fatalf("StartTransfer failed: %v", err)
	}

	if bt.state != StateTransmitting {
		t.Errorf("Expected state TRANSMITTING, got %v", bt.state)
	}
	if bt.currentBlockNumber != 0 {
		t.Errorf("Expected currentBlockNumber 0, got %d", bt.currentBlockNumber)
	}
	if bt.totalExpectedLength == nil {
		t.Error("Expected totalExpectedLength to be set")
	}
}

func TestBlockTransfer_StartTransferTwiceFails(t *testing.T) {
	bt := NewBlockTransfer()

	err := bt.StartTransfer(nil)
	if err != nil {
		t.Fatalf("First StartTransfer failed: %v", err)
	}

	err = bt.StartTransfer(nil)
	if err == nil {
		t.Error("Expected error when starting transfer twice")
	}
}

func TestBlockTransfer_AddSingleBlock(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	block := Block{BlockNumber: 0, Data: []byte("hello"), LastBlock: true}
	err := bt.AddBlock(block)
	if err != nil {
		t.Fatalf("AddBlock failed: %v", err)
	}

	if bt.state != StateComplete {
		t.Errorf("Expected state COMPLETE, got %v", bt.state)
	}
	if !bt.isComplete {
		t.Error("Expected isComplete to be true")
	}
	if len(bt.receivedBlocks) != 1 {
		t.Errorf("Expected 1 received block, got %d", len(bt.receivedBlocks))
	}
}

func TestBlockTransfer_AddMultipleBlocks(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	block1 := Block{BlockNumber: 0, Data: []byte("part1"), LastBlock: false}
	err := bt.AddBlock(block1)
	if err != nil {
		t.Fatalf("AddBlock failed: %v", err)
	}

	if bt.state != StateWaitingForNextBlock {
		t.Errorf("Expected state WAITING_FOR_NEXT_BLOCK, got %v", bt.state)
	}
	if bt.currentBlockNumber != 1 {
		t.Errorf("Expected currentBlockNumber 1, got %d", bt.currentBlockNumber)
	}

	block2 := Block{BlockNumber: 1, Data: []byte("part2"), LastBlock: true}
	err = bt.AddBlock(block2)
	if err != nil {
		t.Fatalf("AddBlock failed: %v", err)
	}

	if bt.state != StateComplete {
		t.Errorf("Expected state COMPLETE, got %v", bt.state)
	}
	if bt.isComplete != true {
		t.Error("Expected isComplete to be true")
	}
	if len(bt.receivedBlocks) != 2 {
		t.Errorf("Expected 2 received blocks, got %d", len(bt.receivedBlocks))
	}
}

func TestBlockTransfer_AddBlockWrongNumberFails(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	block := Block{BlockNumber: 5, Data: []byte("test"), LastBlock: false}
	err := bt.AddBlock(block)
	if err == nil {
		t.Error("Expected error when adding block with wrong number")
	}
}

func TestBlockTransfer_AddBlockInWrongStateFails(t *testing.T) {
	bt := NewBlockTransfer()

	block := Block{BlockNumber: 0, Data: []byte("test"), LastBlock: true}
	err := bt.AddBlock(block)
	if err == nil {
		t.Error("Expected error when adding block in IDLE state")
	}
}

func TestBlockTransfer_ReassembleComplete(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	bt.AddBlock(Block{BlockNumber: 0, Data: []byte("Hello "), LastBlock: false})
	bt.AddBlock(Block{BlockNumber: 1, Data: []byte("World"), LastBlock: true})

	data, err := bt.Reassemble()
	if err != nil {
		t.Fatalf("Reassemble failed: %v", err)
	}

	expected := "Hello World"
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestBlockTransfer_ReassembleIncompleteFails(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	bt.AddBlock(Block{BlockNumber: 0, Data: []byte("part"), LastBlock: false})

	_, err := bt.Reassemble()
	if err == nil {
		t.Error("Expected error when reassembling incomplete transfer")
	}
}

func TestBlockTransfer_GetDataComplete(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	bt.AddBlock(Block{BlockNumber: 0, Data: []byte("data1"), LastBlock: false})
	bt.AddBlock(Block{BlockNumber: 1, Data: []byte("data2"), LastBlock: true})

	data, err := bt.GetData()
	if err != nil {
		t.Fatalf("GetData failed: %v", err)
	}

	expected := "data1data2"
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestBlockTransfer_GetDataIncomplete(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	bt.AddBlock(Block{BlockNumber: 0, Data: []byte("part"), LastBlock: false})

	data, err := bt.GetData()
	if err == nil {
		t.Error("Expected error when getting data from incomplete transfer")
	}
	if data != nil {
		t.Error("Expected nil data from incomplete transfer")
	}
}

func TestBlockTransfer_Reset(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)
	bt.AddBlock(Block{BlockNumber: 0, Data: []byte("test"), LastBlock: false})

	bt.Reset()

	if bt.state != StateIdle {
		t.Errorf("Expected state IDLE after reset, got %v", bt.state)
	}
	if bt.currentBlockNumber != 0 {
		t.Errorf("Expected currentBlockNumber 0 after reset, got %d", bt.currentBlockNumber)
	}
	if len(bt.receivedBlocks) != 0 {
		t.Errorf("Expected 0 received blocks after reset, got %d", len(bt.receivedBlocks))
	}
	if bt.isComplete {
		t.Error("Expected isComplete to be false after reset")
	}
}

func TestBlockTransfer_ValidateSequenceNumbersCorrect(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	bt.AddBlock(Block{BlockNumber: 0, Data: []byte("a"), LastBlock: false})
	bt.AddBlock(Block{BlockNumber: 1, Data: []byte("b"), LastBlock: false})
	bt.AddBlock(Block{BlockNumber: 2, Data: []byte("c"), LastBlock: true})

	err := bt.ValidateSequenceNumbers()
	if err != nil {
		t.Errorf("ValidateSequenceNumbers failed: %v", err)
	}
}

func TestBlockTransfer_ValidateSequenceNumbersIncorrect(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	bt.receivedBlocks = append(bt.receivedBlocks, Block{BlockNumber: 0, Data: []byte("a"), LastBlock: false})
	// Skip block 1
	bt.receivedBlocks = append(bt.receivedBlocks, Block{BlockNumber: 2, Data: []byte("c"), LastBlock: true})

	err := bt.ValidateSequenceNumbers()
	if err == nil {
		t.Error("Expected error when validating incorrect sequence numbers")
	}
}

func TestBlockTransfer_GetBlock(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	block1 := Block{BlockNumber: 0, Data: []byte("first"), LastBlock: false}
	block2 := Block{BlockNumber: 1, Data: []byte("second"), LastBlock: true}

	bt.AddBlock(block1)
	bt.AddBlock(block2)

	retrieved := bt.GetBlock(1)
	if retrieved == nil {
		t.Error("Expected to retrieve block 1")
	}
	if string(retrieved.Data) != "second" {
		t.Errorf("Expected 'second', got %s", string(retrieved.Data))
	}
}

func TestBlockTransfer_GetNonexistentBlock(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	result := bt.GetBlock(5)
	if result != nil {
		t.Error("Expected nil for nonexistent block")
	}
}

func TestBlockTransfer_GetTotalTransmitted(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	if bt.GetTotalTransmitted() != 0 {
		t.Errorf("Expected 0 transmitted blocks, got %d", bt.GetTotalTransmitted())
	}

	bt.AddBlock(Block{BlockNumber: 0, Data: []byte("a"), LastBlock: false})
	if bt.GetTotalTransmitted() != 1 {
		t.Errorf("Expected 1 transmitted block, got %d", bt.GetTotalTransmitted())
	}

	bt.AddBlock(Block{BlockNumber: 1, Data: []byte("b"), LastBlock: true})
	if bt.GetTotalTransmitted() != 2 {
		t.Errorf("Expected 2 transmitted blocks, got %d", bt.GetTotalTransmitted())
	}
}

func TestBlockTransfer_GetProgress(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	progress := bt.GetProgress()
	if progress["state"] != "TRANSMITTING" {
		t.Errorf("Expected state TRANSMITTING, got %v", progress["state"])
	}
	if progress["currentBlock"] != 0 {
		t.Errorf("Expected currentBlock 0, got %v", progress["currentBlock"])
	}
	if progress["totalBlocks"] != 0 {
		t.Errorf("Expected totalBlocks 0, got %v", progress["totalBlocks"])
	}
	if progress["isComplete"] != false {
		t.Error("Expected isComplete to be false")
	}

	bt.AddBlock(Block{BlockNumber: 0, Data: []byte("test"), LastBlock: true})

	progress = bt.GetProgress()
	if progress["state"] != "COMPLETE" {
		t.Errorf("Expected state COMPLETE, got %v", progress["state"])
	}
	if progress["currentBlock"] != 1 {
		t.Errorf("Expected currentBlock 1, got %v", progress["currentBlock"])
	}
	if progress["totalBlocks"] != 1 {
		t.Errorf("Expected totalBlocks 1, got %v", progress["totalBlocks"])
	}
	if progress["isComplete"] != true {
		t.Error("Expected isComplete to be true")
	}
}

func TestBlockTransfer_LargeDataSplit(t *testing.T) {
	bt := NewBlockTransfer()
	expectedLength := 300
	bt.StartTransfer(&expectedLength)

	// Simulate splitting 300 bytes into 3 blocks
	block1 := Block{BlockNumber: 0, Data: bytesRepeat([]byte("x"), 100), LastBlock: false}
	block2 := Block{BlockNumber: 1, Data: bytesRepeat([]byte("y"), 100), LastBlock: false}
	block3 := Block{BlockNumber: 2, Data: bytesRepeat([]byte("z"), 100), LastBlock: true}

	bt.AddBlock(block1)
	bt.AddBlock(block2)
	bt.AddBlock(block3)

	if !bt.isComplete {
		t.Error("Expected transfer to be complete")
	}

	data, err := bt.Reassemble()
	if err != nil {
		t.Fatalf("Reassemble failed: %v", err)
	}

	if len(data) != 300 {
		t.Errorf("Expected 300 bytes, got %d", len(data))
	}
}

func TestBlockTransfer_SingleBlockTransfer(t *testing.T) {
	bt := NewBlockTransfer()
	bt.StartTransfer(nil)

	block := Block{BlockNumber: 0, Data: []byte("short data"), LastBlock: true}
	bt.AddBlock(block)

	if !bt.isComplete {
		t.Error("Expected transfer to be complete")
	}

	data, err := bt.Reassemble()
	if err != nil {
		t.Fatalf("Reassemble failed: %v", err)
	}

	expected := "short data"
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

// Helper function to repeat bytes
func bytesRepeat(b []byte, count int) []byte {
	result := make([]byte, len(b)*count)
	for i := 0; i < count; i++ {
		copy(result[i*len(b):], b)
	}
	return result
}
