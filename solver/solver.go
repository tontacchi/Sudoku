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

// check that this can be converted into a boolean function. we only return -1 and 0? so it can be true and false?
func (bt *BacktrackSolver) Solve(timeLeft time.Duration) int {
	if timeLeft <= time.Minute {
		return -1
	} else if bt.HasSolution {
		return 0
	}
	
	start := time.Now()
	
	variable := bt.Select(bt.Network)
	if variable == nil {
		bt.HasSolution = true
		return 0
	}

	for _, value := range bt.OrderValues(variable, bt.Network) {
		bt.Trail.PlaceMarker()
		bt.Trail.Push(variable)
		variable.AssignValue(value)

		if bt.Enforce(bt.Network, bt.Trail) {
			remainingTime := timeLeft - time.Since(start)

			if bt.Solve(remainingTime) == -1 {
				return -1
			}
		}

		if bt.HasSolution {
			return 0
		}

		bt.Trail.Undo()
	}

	return 0
}


// Default Strategy Implementations

type FirstUnassigned struct{}

func (FirstUnassigned) Select(network *Network) *Variable {
	for _, variable := range network.variables {
		if !variable.assigned { return variable }
	}

	return nil
}


type MostRecentValue struct{}

func (MostRecentValue) Select(network *Network) *Variable {
	var bestVar *Variable

	minSize := -1
	for _, variable := range network.variables {
		if !variable.assigned {
			size := variable.Size()

			if minSize == -1 || size < minSize {
				minSize = size
				bestVar = variable
			}
		}
	}

	return bestVar
}

