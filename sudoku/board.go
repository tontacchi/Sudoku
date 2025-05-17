package sudoku

import (
	// "os"
	// "strings"
	"strconv"
	"math/rand/v2"
	// "time"
	// "fmt"
)

// Represents a Sudoku board
type Board struct {
	boxRows  int
	boxCols  int
	boardLen int
	cells    [][]int
}

func NewRandomBoard(boxRows, boxCols, numHints int) *Board {
	boardLen := boxRows * boxCols

	board := &Board{
		boxRows:  boxRows,
		boxCols:  boxCols,
		boardLen: boardLen,
		cells:    make([][]int, boardLen),
	}
	
	// initalize board matrix
	for index := range boardLen {
		board.cells[index] = make([]int, boardLen)
	}

	// seed automatically set as of go 1.24


	// fill board w/ randomly placed hints
	for numHints > 0 {
		row, col := rand.IntN(boardLen), rand.IntN(boardLen)
		value := rand.IntN(boardLen + 1)

		if board.cells[row][col] == 0 && board.isValidPlacement(row, col, value) {
			board.cells[row][col] = value
			numHints--
		}
	}

	return board
}

func (b *Board) String() string {
	res := ""

	for _, row := range b.cells {
		for _, value := range row {
			res += strconv.Itoa(value)
			res += " "
		}
		res += "\n"
	}

	return res
}

//---[ Internal Helpers ]-------------------------------------------------------

func (b *Board) isValidPlacement(row, col, value int) bool {
	// row & col check: value doesn't elsewhere in row or col
	for check := range b.boardLen {
		if b.cells[row][check] == value || b.cells[check][col] == value {
			return false
		}
	}

	// box check: value doesn't appear elsewhere in its smaller box
	startRow := (row / b.boxRows) * b.boxRows
	startCol := (col / b.boxCols) * b.boxCols

	for r := startRow; r < startRow + b.boxRows; r++ {
		for c := startCol; c < startCol + b.boxCols; c++ {
			if b.cells[r][c] == value {
				return false
			}
		}
	}

	return true
}



