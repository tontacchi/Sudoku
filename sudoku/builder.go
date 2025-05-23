package sudoku

import (
	"sudoku-csp/solver"
	"math/rand/v2"
	"time"
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

			// boxIndex := (row % board.BoxRows) * board.BoxRows + (col / board.BoxCols)
			// fmt.Printf("(%d / %d) * %d + (%d / %d) -> %d\n", row, board.BoxRows, board.BoxCols, col, board.BoxCols, boxIndex)

			boxColsCount := board.BoardLen() / board.BoxCols
			boxIndex := (row / board.BoxRows) * boxColsCount + (col / board.BoxCols)

			// boxIndex := (row / board.BoxRows) * board.BoxCols + (col / board.BoxCols)
			// fmt.Printf("(%d / %d) * %d + (%d / %d) -> %d\n", row, board.BoxRows, board.BoxRows, col, board.BoxCols, boxIndex)

			// rowOffset := row / board.BoxRows
			// colOffset := col / board.BoxCols
			//
			// [0][0] = 0
			// [0][1] = 1
			// [1][0] = 2  // (2, 1) 
			// [1][1] = 3

			// boxIndex := (col / board.BoxCols) * board.BoxRows + (row / board.BoxRows)
			// fmt.Printf("(%d, %d) assigned box %d\n\n", row, col, boxIndex)
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

	// fmt.Printf("vars in row: %d, vars in col: %d, vars in box: %d\n", len(rowGroups), len(colGroups), len(boxGroups))
	// for i := range boardLen {
	// 	fmt.Printf("Row %d has %d vars\n", i, len(rowGroups[i]))
	// 	fmt.Printf("Col %d has %d vars\n", i, len(colGroups[i]))
	// 	fmt.Printf("Box %d has %d vars\n", i, len(boxGroups[i]))
	// }

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

func NewForwardCheckGeneratedBoard(boxRows, boxCols, numHints int) *Board {
	board   := NewEmptyBoard(boxRows, boxCols)
	network := NewNetworkFromBoard(board)
	trail   := solver.NewTrail()

	varSelector := solver.MRV{}
	valSelector := solver.LeastConstrainingValue{}
	checker     := solver.ForwardChecking{}

	numAssigned := 0
	
	for numAssigned < numHints {
		variable := varSelector.Select(network)
		if variable == nil { break }

		values := valSelector.OrderValues(variable, network)
		rand.Shuffle(len(values), func(left, right int) {
			values[left], values[right] = values[right], values[left]
		})

		placedValue := false
		for _, value := range values {
			trail.PlaceMarker()
			trail.Push(variable)

			variable.AssignValue(value)

			if checker.Enforce(network, trail) {
				numAssigned++
				placedValue = true

				break
			}

			trail.Undo()
		}

		if !placedValue { break }
	}

	return NewBoardFromNetwork(network, boxRows, boxCols)
}

func NewBoardFromSolved(boxRows, boxCols, numHints int) *Board {
	totalCells := boxRows * boxCols * boxRows * boxCols
	
	// fix later
	if numHints > totalCells {
		numHints = totalCells
	}

	board   := NewEmptyBoard(boxRows, boxCols)
	network := NewNetworkFromBoard(board)
	trail   := solver.NewTrail()
	
	varSelector := solver.MRV{}
	valSelector := solver.LeastConstrainingValue{}
	checker     := solver.ForwardChecking{}

	solver := solver.NewBacktrackSolver(network, trail, varSelector, valSelector, checker)
	solver.Solve(time.Minute * 2)

	board = NewBoardFromNetwork(network, boxRows, boxCols)
	fmt.Printf("[ DEBUG ] the solved board:\n%v\n", board)
	
	positions := allCellCoords(board.BoardLen())
	rand.Shuffle(len(positions), func(left, right int) {
		positions[left], positions[right] = positions[right], positions[left]
	})

	for _, pair := range positions[:len(positions) - numHints] {
		row, col := pair[0], pair[1]
		board.Cells[row][col] = 0
	}

	return board
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

func allCellCoords(dimension int) [][]int {
	numPairs := dimension * dimension

	pairs := make([][]int, numPairs)

	for row := range dimension {
		for col := range dimension {
			pairs[row * dimension + col] = []int{row, col}
		}
	}

	return pairs
}

