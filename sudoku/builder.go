package sudoku

import "sudoku-csp/solver"

func NewNetworkFromBoard(board *Board) *solver.Network {
	network := solver.NewNetwork()

	boardLen := board.BoardLen()
	tempVars := make([]*solver.Variable, 0, boardLen * boardLen)

	// var domain []int
	//
	// if value == 0 {
	// 	for i := 1; i <= boardLen; i++ {
	// 		domain = append(domain, i)
	// 	}
	// } else {
	// 	domain = []int{value}
	// }

	// 1) build variables from board cells
	for row := range boardLen {
		for col := range boardLen {
			value    := board.Cells[row][col]
			domain   := domainFor(value, boardLen)
			boxIndex := (row / board.BoxRows) * board.BoxCols + (col / board.BoxCols)
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
		rowGroups[variable.Col] = append(colGroups[variable.Col], variable)
		rowGroups[variable.Block] = append(boxGroups[variable.Block], variable)
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

