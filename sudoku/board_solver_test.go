package sudoku

import (
	"testing"
	"sudoku-csp/solver"
	"time"
)

func TestSolverRandom(t *testing.T) {
	const NUM_ROWS  = 2
	const NUM_COLS  = 2
	const NUM_HINTS = 0

	board := NewRandomBoard(NUM_ROWS, NUM_COLS, NUM_HINTS)
	if board == nil {
		t.Fatal("Generated board was nil")
	}

	t.Log("Initial Board:")
	t.Logf("\n%s\n", board)

	network := NewNetworkFromBoard(board)
	trail := solver.NewTrail()

	s := solver.NewBacktrackSolver(
		network,
		trail,
		solver.MRVWithDegree{},
		solver.LeastConstrainingValue{},
		solver.BasicCheck{},
	)

	res := s.Solve(time.Minute)
	t.Log("Solver result:", res)

	if s.HasSolution {
		solved := NewBoardFromNetwork(network, NUM_ROWS, NUM_COLS)
		t.Logf("Solved Board:\n%s\n", solved)
	} else {
		t.Error("Solver could not find a solution")
		solved := NewBoardFromNetwork(network, NUM_ROWS, NUM_COLS)
		t.Logf("Solved Board:\n%s\n", solved)
	}
}

func estSolver(t *testing.T) {
	const NUM_ROWS  = 3
	const NUM_COLS  = 2
	const NUM_HINTS = 500

	board := NewEmptyBoard(NUM_ROWS, NUM_COLS)

	if board == nil {
		t.Fatal("Generated board was nil")
	}

	t.Log("Initial Board:")
	t.Logf("\n%s\n", board)

	network := NewNetworkFromBoard(board)
	trail := solver.NewTrail()

	s := solver.NewBacktrackSolver(
		network,
		trail,
		solver.MRVWithDegree{},
		solver.LeastConstrainingValue{},
		solver.BasicCheck{},
	)

	res := s.Solve(time.Minute)
	t.Log("Solver result:", res)

	if s.HasSolution {
		solved := NewBoardFromNetwork(network, NUM_ROWS, NUM_COLS)
		t.Logf("Solved Board:\n%s\n", solved)
	} else {
		t.Error("Solver could not find a solution")
	}
}

