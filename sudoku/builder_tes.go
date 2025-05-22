package sudoku

import (
	"testing"
	// "fmt"
)

func TestNewNetworkFromBoard(t *testing.T) {
	const NUM_ROWS  int = 10
	const NUM_COLS  int = 10
	const NUM_HINTS int = 2000

	board := NewRandomBoard(NUM_ROWS, NUM_COLS, NUM_HINTS)
	// board := NewForwardCheckGeneratedBoard(NUM_ROWS, NUM_COLS, NUM_HINTS)
	if board == nil {
		t.Fatal("NewRandomBoard returned nil")
	}

	const EXPECTED_BOARDLEN int = NUM_ROWS * NUM_COLS
	if board.boardLen != EXPECTED_BOARDLEN {
		t.Errorf("Expected boardLen to be %v, got %v\n", EXPECTED_BOARDLEN, board.boardLen)
		return
	}

	t.Log(board)
	network := NewNetworkFromBoard(board)

	t.Log(network.String())
}

