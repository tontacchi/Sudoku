package solver

import (
	"time"
	"slices"
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


type MinimumRemainingValue struct{}

func (MinimumRemainingValue) Select(network *Network) *Variable {
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


type DefaultValOrder struct{}

func (DefaultValOrder) OrderVals(variable *Variable, network *Network) []int {
	values := variable.Values()
	return append([]int{}, values...)
}


type BasicCheck struct{}

func (BasicCheck) Enforce(network *Network, trail *Trail) bool {
	for _, constraint := range network.constraints {
		if !constraint.IsSatisfied() { return false }
	}

	return true
}


type ForwardChecking struct{}

func (ForwardChecking) Enforce(network *Network, trail *Trail) bool {
	for _, modifiedConstraint := range network.GetModifiedConstraints() {
		for _, variable := range modifiedConstraint.Variables() {
			if !variable.assigned { continue }

			value := variable.Assignment()
			for _, neighbor := range network.GetNeighbors(variable) {
				if !neighbor.domain.Contains(value) { continue }

				trail.Push(neighbor)
				neighbor.domain.Remove(value)

				if neighbor.domain.Empty() { return false }

				neighbor.modified = true
			}
		}
	}

	for _, constraint := range network.constraints {
		if !constraint.IsSatisfied() { return false }
	}

	return true
}

