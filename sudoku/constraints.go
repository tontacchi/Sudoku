package sudoku

import "sudoku-csp/solver"

func MakeAllDiff(variables []*solver.Variable) solver.Constraint {
	return solver.NewAllDiffConstraint(variables)
}


