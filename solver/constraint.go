package solver

// generic CSP constraint interface
type Constraint interface {
	Variables()   []*Variable
	IsSatisfied() bool
	IsModified()  bool
	String()      string
}

// Constraint: all values in the set are unique
type AllDiffConstraint struct {
	variables []*Variable
}

func NewAllDiffConstraint(variables []*Variable) *AllDiffConstraint {
	return &AllDiffConstraint{
		variables: variables,
	}
}

// Mutators
func (c *AllDiffConstraint) AddVariable(variable *Variable) {
	c.variables = append(c.variables, variable)
}

// Accessors & Constraint Interface
func (c *AllDiffConstraint) Variables() []*Variable {
	return c.variables
}

func (c *AllDiffConstraint) IsModified() bool {
	for _, variable := range c.variables {
		if variable.modified { return true }
	}

	return false
}

func (c *AllDiffConstraint) IsSatisfied() bool {
	set := make(map[int]bool)

	for _, variable := range c.variables {
		if !variable.assigned { continue }

		assignment := variable.Assignment()
		if set[assignment] {
			return false
		}

		set[assignment] = true
	}

	return true
}

func (c *AllDiffConstraint) String() string {
	repr := "{"

	for index, variable := range c.variables {
		if index > 0 {
			repr += ", "
		}

		repr += variable.String()
	}

	return repr
}

