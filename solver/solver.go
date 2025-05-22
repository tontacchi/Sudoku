package solver

import (
	"time"
	"sort"
	"fmt"
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
func (bt *BacktrackSolver) Solve(timeLeft time.Duration) bool {
	if timeLeft <= 0 {
		return false
	} else if bt.HasSolution {
		return true
	}
	
	start := time.Now()
	
	variable := bt.Select(bt.Network)
	if variable == nil {
		bt.HasSolution = true
		return true
	}

	for _, value := range bt.OrderValues(variable, bt.Network) {
		bt.Trail.PlaceMarker()
		bt.Trail.Push(variable)
		variable.AssignValue(value)

		if bt.Enforce(bt.Network, bt.Trail) {
			remainingTime := timeLeft - time.Since(start)

			if bt.Solve(remainingTime) == false {
				return false
			}
		}

		if bt.HasSolution {
			return true
		}

		bt.Trail.Undo()
	}

	fmt.Println("darn no solution")
	return false
}


// Default Strategy Implementations

// 1) Variable Selectors

type FirstUnassigned struct{}

func (FirstUnassigned) Select(network *Network) *Variable {
	for _, variable := range network.variables {
		if !variable.assigned { return variable }
	}

	return nil
}


type MRV struct{}

func (MRV) Select(network *Network) *Variable {
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


type MRVWithDegree struct{}

func (MRVWithDegree) Select(network *Network) *Variable {
	var candidates []*Variable
	minDomainSize := -1

	for _, variable := range network.variables {
		if variable.assigned { continue }

		size := variable.Size()
		if minDomainSize == -1 || size < minDomainSize {
			minDomainSize = size
			candidates = []*Variable{variable}
		} else {
			candidates = append(candidates, variable)
		}
	}

	if len(candidates) == 1 {
		return candidates[0]
	}

	// tie-breaker: highest unassigned neighbor count
	var best *Variable
	maxDegree := -1

	for _, variable := range candidates {
		degree := 0

		for _, neighbor := range network.GetNeighbors(variable) {
			if !neighbor.assigned { degree++ }
		}

		if degree > maxDegree {
			maxDegree, best = degree, variable
		}
	}

	return best
}


// 2) Value Selectors

type DefaultValOrder struct{}

func (DefaultValOrder) OrderValues(variable *Variable, network *Network) []int {
	values := variable.Values()
	return append([]int{}, values...)
}


type LeastConstrainingValue struct{}

func (LeastConstrainingValue) OrderValues(variable *Variable, network *Network) []int {
	neighbors   := network.GetNeighbors(variable)
	valueImpact := make(map[int]int)

	for _, value := range variable.Values() {
		impact := 0

		for _, neighbor := range neighbors {
			if !neighbor.assigned && neighbor.domain.Contains(value) {
				impact++
			}
		}

		valueImpact[value] = impact
	}

	type valueImpactPair struct {
		value  int
		impact int
	}
	var pairs []valueImpactPair

	for value, impact := range valueImpact {
		pairs = append(pairs, valueImpactPair{value, impact})
	}

	sort.Slice(pairs, func(left, right int) bool {
		return pairs[left].impact < pairs[right].impact
	})

	ordered := make([]int, len(pairs))
	for index, pair := range pairs {
		ordered[index] = pair.value
	}

	return ordered
}


// 3) Consistency Checkers

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


type NorvigCheck struct{}

func (NorvigCheck) Enforce(network *Network, trail *Trail) bool {
	// 1) forward checking
	for _, constraint := range network.GetModifiedConstraints() {
		for _, variable := range constraint.Variables() {
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

	// 2) only-choice rule
	for _, constraint := range network.Constraints() {
		valueVariablesMap := make(map[int][]*Variable)

		for _, variable := range constraint.Variables() {
			if variable.assigned { continue }

			for _, value := range variable.domain.values {
				valueVariablesMap[value] = append(valueVariablesMap[value], variable)
			}
		}

		for value, variables := range valueVariablesMap {
			if len(variables) != 1 { continue }

			leadVariable := variables[0]
			trail.Push(leadVariable)

			leadVariable.AssignValue(value)
			leadVariable.assigned = true
		}
	}

	for _, constraint := range network.Constraints() {
		if !constraint.IsSatisfied() { return false }
	}

	return true
}


type ArcConsistency struct{}

func (ArcConsistency) Enforce(network *Network, trail *Trail) bool {
	queue := []*Variable{}

	for _, variable := range network.variables {
		if !variable.assigned { continue }	

		queue = append(queue, variable)
	}

	for len(queue) > 0 {
		variable := queue[0]
		value    := variable.Assignment()
		queue     = queue[1:]
	
		for _, neighbor := range network.GetNeighbors(variable) {
			if neighbor.assigned || !neighbor.domain.Contains(value) { continue }

			trail.Push(neighbor)
			neighbor.domain.Remove(value)

			if neighbor.domain.Empty() { return false }
			neighbor.modified = true

			if neighbor.Size() != 1 { continue }

			valueToAssign := neighbor.Values()[0]
			neighbor.AssignValue(valueToAssign)
			queue = append(queue, neighbor)
		}
	}

	for _, constraint := range network.constraints {
		if !constraint.IsSatisfied() { return false }
	}

	return true
}

