package solver

import (
	"./sudoku"
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
	adjList     map[*Variable][]Constraint
}

func NewNetworkFromBoard(board *sudoku.Board) *Network {
	return &Network{}
}


// Mutators

func (n *Network) AddVariable(variable *Variable) {
	n.variables = append(n.variables, variable)
	n.adjList[variable] = []Constraint{}
}

func (n *Network) AddConstraint(constraint Constraint) {
	n.constraints = append(n.constraints, constraint)

	for _, variable := range constraint.Variables() {
		n.adjList[variable] = append(n.adjList[variable], constraint)
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
	return n.adjList[variable]
}

func (n *Network) GetNeighbors(variable *Variable) []*Variable {
	set := map[*Variable]bool{}

	for _, constraint := range n.adjList[variable] {
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
			for _, constraint := range n.adjList[variable] {
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

