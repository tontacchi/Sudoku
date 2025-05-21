package solver

import (
	"time"
)

// Strategy Interfaces
type VarSelector interface {
	Select(*Network) *Variable
}

type ValSelector interface {
	OrderValues(*Variable, *Network) []int
}

type ConsistencyChecker interface {
	Enforce(*Network, *Trail) bool
}


// Backtrack Solver w/ injected strategy

type BacktrackSolver struct {
	Network     *Network
	Trail       *Trail
	HasSolution bool

	VarSelector         // Select()
	ValSelector         // OrderValues()
	ConsistencyChecker  // Enforce()
}

func NewBacktrackSolver(
	network     *Network,
	trail       *Trail,
	varSelector VarSelector,
	valSelector ValSelector,
	checker     ConsistencyChecker,
) *BacktrackSolver {
	return &BacktrackSolver{	
		Network:            network,
		Trail:              trail,
		VarSelector:        varSelector,
		ValSelector:        valSelector,
		ConsistencyChecker: checker,
		HasSolution:        false,
	}
}


// Solver Logic

