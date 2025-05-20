package solver

import "fmt"

/*
Represents a CSP variable
*/
type Variable struct {
	domain *Domain

	Row   int
	Col   int
	Block int

	assigned   bool
	modified   bool
	changeable bool
}

func NewVariable(values []int, row, col, block int) *Variable {
	variable := &Variable{
		domain:     NewDomain(values...),
		Row:        row,
		Col:        col,
		Block:      block,
		changeable: len(values) > 1,
		assigned:   len(values) == 1,
		modified:   len(values) == 1,
	}

	return variable
}

func (v *Variable) Copy() *Variable {
	variable := &Variable{
		domain:     v.domain.Copy(),
		Row:        v.Row,
		Col:        v.Col,
		Block:      v.Block,
		assigned:   v.assigned,
		modified:   v.modified,
		changeable: v.changeable,
	}

	return variable
}


// Accessors

// Number of values left in domain
func (v *Variable) Size() int {
	return v.domain.Size()
}

// all values in domain
func (v *Variable) Values() []int {
	return v.domain.values
}

// assignment is remaining value in domain + variable set to assigned. otherwise, 0 representing unassigned variable
func (v *Variable) Assignment() int {
	if v.assigned && v.Size() == 1 {
		return v.Values()[0]
	}

	return 0
}


// Mutators

func (v *Variable) AssignValue(value int) {
	if !v.changeable { return }

	v.domain = NewDomain()
	v.assigned, v.modified = true, true
}

func (v *Variable) Unassign() {
	if v.changeable {
		v.assigned = false
	}
}

func (v *Variable) RemoveValueFromDomain(value int) {
	if !v.changeable { return }

	if v.domain.Remove(value) {
		v.modified = true
	}
}


func (v *Variable) String() string {
	return fmt.Sprintf("Row: %d, Col: %d, Block: %d, value: %d", v.Row, v.Col, v.Block, v.Assignment())
}

