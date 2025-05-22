package sudoku

import (
	"sudoku-csp/solver"
	"math/rand/v2"
)

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

func NewForwardCheckGeneratedBoard(boxRows, boxCols, clueCount int) *Board {
	board   := NewEmptyBoard(boxRows, boxCols)
	network := NewNetworkFromBoard(board)
	trail   := solver.NewTrail()

	varSelector := solver.MRV{}
	valSelector := solver.LeastConstrainingValue{}
	checker     := solver.ForwardChecking{}

	numAssigned := 0

	// go 1.24 -> seed automatically set

	for numAssigned < clueCount {
		variable := varSelector.Select(network)

		// no more assignable values
		if variable == nil { break }

		values := valSelector.OrderValues(variable, network)
		rand.Shuffle(len(values), func(left, right int) {
			values[left], values[right] = values[right], values[left]
		})

		isAssigned := false
		for _, value := range values {
			trail.PlaceMarker()
			trail.Push(variable)

			variable.AssignValue(value)

			if checker.Enforce(network, trail) {
				numAssigned++
				isAssigned = true
				break
			}

			trail.Undo()
		}

		// cannot assign to variable safely -> backtrack
		if !isAssigned {
			break
		}
	}

	finalBoard := NewBoardFromNetwork(network, boxRows, boxCols)
	return finalBoard
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

