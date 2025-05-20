package solver

import (
	"sudoku/sudoku"
	"fmt"
)

/*
TODO:
- NewNetworkFromBoard()
- String()
*/

type Network struct {
	variables   []*Variable
	constraints []Constraint
	varToConst  map[*Variable][]Constraint
}

func NewNetworkFromBoard(board *sudoku.Board) *Network {
	network := &Network{
		variables:   []*Variable{},
		constraints: []Constraint{},
		varToConst:  make(map[*Variable][]Constraint),
	}

	boardLen := board.BoardLen()
	tempVars := make([]*Variable, 0, boardLen * boardLen)

	// 1) build variables from board cells
	for row := range boardLen {
		for col := range boardLen {
			val := board.Cells[row][col]

			var domain []int

			if val == 0 {
				for i := 1; i <= boardLen; i++ {
					domain = append(domain, i)
				}
			} else {
				domain = []int{val}
			}
			boxIndex := (row / board.BoxRows) * board.BoxCols + (col / board.BoxCols)

			variable := NewVariable(domain, row, col, boxIndex)
			tempVars = append(tempVars, variable)

			network.AddVariable(variable)
		}
	}

	// 2) group variables by rows, cols, & boxes
	rowGroups := make(map[int][]*Variable)
	colGroups := make(map[int][]*Variable)
	boxGroups := make(map[int][]*Variable)

	for _, variable := range tempVars {
		rowGroups[variable.row] = append(rowGroups[variable.row], variable)
		rowGroups[variable.col] = append(colGroups[variable.col], variable)
		rowGroups[variable.block] = append(boxGroups[variable.block], variable)
	}

	// 3) assign constraints for rows, cols, & boxes
	for _, group := range []map[int][]*Variable{rowGroups, colGroups, boxGroups} {
		for _, vars := range group {
			constraint := NewAllDiffConstraint(vars)
			network.AddConstraint(constraint)
		}
	}

	return network
}


// Mutators

func (n *Network) AddVariable(variable *Variable) {
	n.variables = append(n.variables, variable)
	n.varToConst[variable] = []Constraint{}
}

func (n *Network) AddConstraint(constraint Constraint) {
	n.constraints = append(n.constraints, constraint)

	for _, variable := range constraint.Variables() {
		n.varToConst[variable] = append(n.varToConst[variable], constraint)
	}
}


// Accessors

func (n *Network) Variables() []*Variable {
	return n.variables
}

func (n *Network) Constraints() []Constraint {
	return n.constraints
}

func (n *Network) GetConstraints(variable *Variable) []Constraint {
	return n.varToConst[variable]
}

func (n *Network) GetNeighbors(variable *Variable) []*Variable {
	set := map[*Variable]bool{}

	for _, constraint := range n.varToConst[variable] {
		for _, other := range constraint.Variables() {
			if !set[other] && other != variable {
				set[other] = true
			}
		}
	}

	neighbors := make([]*Variable, 0, len(set))
	for neighbor := range set {
		neighbors = append(neighbors, neighbor)
	}

	return neighbors
}

func (n *Network) IsConsistent() bool {
	for _, constraint := range n.constraints {
		if !constraint.IsSatisfied() { return false }
	}

	return true
}

func (n *Network) GetModifiedConstraints() []Constraint {
	set := map[Constraint]bool{}
	modified := []Constraint{}
	
	for _, variable := range n.variables {
		if variable.modified {
			for _, constraint := range n.varToConst[variable] {
				if !set[constraint] {
					modified = append(modified, constraint)
					set[constraint] = true
				}
			}
		}

		variable.modified = false
	}

	return modified
}

func (n *Network) String() string {
	res := fmt.Sprintf("")



	return res
}

