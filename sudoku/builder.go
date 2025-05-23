package sudoku

import (
	"sudoku-csp/solver"
	"fmt"
)

func NewNetworkFromBoard(board *Board) *solver.Network {
	network := solver.NewNetwork()

	boardLen := board.BoardLen()
	tempVars := make([]*solver.Variable, 0, boardLen * boardLen)

	// 1) build variables from board cells
	for row := range boardLen {
		for col := range boardLen {
			value    := board.Cells[row][col]
			domain   := domainFor(value, boardLen)
			boxIndex := (row / board.BoxRows) * board.BoxCols + (col / board.BoxCols)
			fmt.Printf("(%d, %d) assigned box %d\n", row, col, boxIndex)
			variable := solver.NewVariable(domain, row, col, boxIndex)

			tempVars = append(tempVars, variable)
			network.AddVariable(variable)
		}
	}

	// 2) group variables by rows, cols, & boxes
	rowGroups := make(map[int][]*solver.Variable)
	colGroups := make(map[int][]*solver.Variable)
	boxGroups := make(map[int][]*solver.Variable)

	for _, variable := range tempVars {
		rowGroups[variable.Row] = append(rowGroups[variable.Row], variable)
		colGroups[variable.Col] = append(colGroups[variable.Col], variable)
		boxGroups[variable.Block] = append(boxGroups[variable.Block], variable)
	}

	fmt.Printf("vars in row: %d, vars in col: %d, vars in box: %d\n", len(rowGroups), len(colGroups), len(boxGroups))
	for i := range boardLen {
		fmt.Printf("Row %d has %d vars\n", i, len(rowGroups[i]))
		fmt.Printf("Col %d has %d vars\n", i, len(colGroups[i]))
		fmt.Printf("Box %d has %d vars\n", i, len(boxGroups[i]))
	}

	// 3) assign constraints for rows, cols, & boxes
	for _, group := range []map[int][]*solver.Variable{rowGroups, colGroups, boxGroups} {
		for _, vars := range group {
			constraint := MakeAllDiff(vars)
			network.AddConstraint(constraint)
		}
	}

	return network
}

func NewBoardFromNetwork(network *solver.Network, boxRows, boxCols int) *Board {
	boardLen := boxRows * boxCols
	cells    := make([][]int, boardLen)

	// initalize row slices
	for row := range cells {
		cells[row] = make([]int, boardLen)
	}

	for _, variable := range network.Variables() {
		if !variable.Assigned() { continue }

		cells[variable.Row][variable.Col] = variable.Assignment()
	}

	return &Board{
		BoxRows:  boxRows,
		BoxCols:  boxCols,
		boardLen: boardLen,
		Cells:    cells,
	}
}

func domainFor(value, size int) []int {
	if value != 0 {
		return []int{value}
	}

	domain := make([]int, size)
	for i := range size {
		domain[i] = i + 1
	}

	return domain
}

