package solver

import (
	"fmt"
)

type Network struct {
	variables   []*Variable
	constraints []Constraint
	varToConst  map[*Variable][]Constraint
}

func NewNetwork() *Network {
	return &Network {
		variables:   []*Variable{},
		constraints: []Constraint{},
		varToConst:  map[*Variable][]Constraint{},
	}
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
	// variables
	res := fmt.Sprintf("%d Variables {", len(n.variables))	
	for index, variable := range n.variables {
		if index > 0 {
			res += ", "
		}

		res += variable.String() + "\n"
	}
	res += "\n"

	// constriants
	res += fmt.Sprintf("%d Constraints {", len(n.constraints))
	for _, constraint := range n.constraints {
		res += constraint.String() + "\n"
	}
	res += "}"

	return res
}

