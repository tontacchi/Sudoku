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
	NewNetworkFromBoard(board)

}

