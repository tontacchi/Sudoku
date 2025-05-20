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
	BoxRows  int
	BoxCols  int
	boardLen int
	Cells    [][]int
}

func NewRandomBoard(BoxRows, BoxCols, numHints int) *Board {
	boardLen := BoxRows * BoxCols

	board := &Board{
		BoxRows:  BoxRows,
		BoxCols:  BoxCols,
		boardLen: boardLen,
		Cells:    make([][]int, boardLen),
	}
	
	// initalize board matrix
	for index := range boardLen {
		board.Cells[index] = make([]int, boardLen)
	}

	// seed automatically set as of go 1.24


	// fill board w/ randomly placed hints
	for numHints > 0 {
		row, col := rand.IntN(boardLen), rand.IntN(boardLen)
		value := rand.IntN(boardLen + 1)

		if board.Cells[row][col] == 0 && board.isValidPlacement(row, col, value) {
			board.Cells[row][col] = value
			numHints--
		}
	}

	return board
}

func (b *Board) BoardLen() int {
	return b.boardLen
}

func (b *Board) String() string {
	res := ""

	for _, row := range b.Cells {
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
		if b.Cells[row][check] == value || b.Cells[check][col] == value {
			return false
		}
	}

	// box check: value doesn't appear elsewhere in its smaller box
	startRow := (row / b.BoxRows) * b.BoxRows
	startCol := (col / b.BoxCols) * b.BoxCols

	for r := startRow; r < startRow + b.BoxRows; r++ {
		for c := startCol; c < startCol + b.BoxCols; c++ {
			if b.Cells[r][c] == value {
				return false
			}
		}
	}

	return true
}



