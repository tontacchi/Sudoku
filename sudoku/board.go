package sudoku

import (
	// "os"
	"strings"
	// "strconv"
	"math/rand/v2"
	// "time"
	"fmt"
)

// Represents a Sudoku board
type Board struct {
	BoxRows  int
	BoxCols  int
	boardLen int
	Cells    [][]int
}

func NewEmptyBoard(boxRows, boxCols int) *Board {
	board := &Board{
		BoxRows:  boxRows,
		BoxCols:  boxCols,
		boardLen: boxRows * boxCols,
		Cells:    make([][]int, boxRows * boxCols),
	}

	// initalize rows
	for row := range board.boardLen {
		board.Cells[row] = make([]int, board.boardLen)
	}

	return board
}

func NewRandomBoard(boxRows, boxCols, numHints int) *Board {
	if numHints > (boxRows * boxRows * boxCols * boxCols) {
		numHints = (boxRows * boxRows * boxCols * boxCols)
	}

	boardLen := boxRows * boxCols

	board := &Board{
		BoxRows:  boxRows,
		BoxCols:  boxCols,
		boardLen: boardLen,
		Cells:    make([][]int, boardLen),
	}
	
	// initalize board matrix
	for row := range boardLen {
		board.Cells[row] = make([]int, boardLen)
	}

	// go 1.24 -> seed automatically set

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

// Accessors

func (b *Board) BoardLen() int {
	return b.boardLen
}

func (b *Board) String() string {
	maxVal   := b.boardLen
	digitLen := len(fmt.Sprintf("%d", maxVal))
	blankSymbol := strings.Repeat("~", digitLen)

	var builder strings.Builder
	for row := range b.boardLen {
		// padding on the side
		builder.WriteString(strings.Repeat(" ", digitLen))

		for col := range b.boardLen {
			value := b.Cells[row][col]

			// display cell content
			if value == 0 {
				builder.WriteString(fmt.Sprintf("%*s", digitLen, blankSymbol))
			} else {
				builder.WriteString(fmt.Sprintf("%*d", digitLen, value))
			}

			// display separator
			atBoxBorder := (col + 1) % b.BoxCols == 0 && col != b.boardLen - 1
			if atBoxBorder {
				builder.WriteString(" | ")
			} else {
				builder.WriteString(" ")
			}
		}
		builder.WriteString("\n")

		// horizontal box separator
		atBoxBottom := (row + 1) % b.BoxRows == 0 && row != b.boardLen - 1
		if atBoxBottom {
			totalWidth := b.boardLen * (digitLen + 1) + (b.BoxCols - 1) * 3
			builder.WriteString(strings.Repeat("-", totalWidth) + "\n")
		}
	}

	return builder.String()
}

// Internal Helpers

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

