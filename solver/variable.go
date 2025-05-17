package solver

/*
Represents a CSP variable
*/
type Variable struct {
	domain Domain

	row   int
	col   int
	block int

	assigned   bool
	modified   bool
	changeable bool
}

// Accessors
func (v *Variable) Size() int {
	return v.domain.Size()
}

func (v *Variable) Assigned() bool {
	return v.assigned
}

func (v *Variable) Modified() bool {
	return v.modified
}

func (v *Variable) Changeable() bool {
	return v.changeable
}

func (v *Variable) Domain() Domain {
	return v.domain
}

func (v *Variable) Values() []int {
	return v.domain.values
}

func (v *Variable) Assignment() int {
	if v.assigned == false {

	}
}

