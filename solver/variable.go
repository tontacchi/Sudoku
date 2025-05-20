package solver

import "fmt"

/*
Represents a CSP variable
*/
type Variable struct {
	domain *Domain

	row   int
	col   int
	block int

	assigned   bool
	modified   bool
	changeable bool
}

func NewVariable(values []int, row, col, block int) *Variable {
	variable := &Variable{
		domain:     NewDomain(values...),
		row:        row,
		col:        col,
		block:      block,
		changeable: len(values) > 1,
		assigned:   len(values) == 1,
		modified:   len(values) == 1,
	}

	return variable
}

func (v *Variable) Copy() *Variable {
	variable := &Variable{
		domain:     v.domain.Copy(),
		row:        v.row,
		col:        v.col,
		block:      v.block,
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
	return fmt.Sprintf("row: %d, col: %d, block: %d, value: %d", v.row, v.col, v.block, v.Assignment())
}

