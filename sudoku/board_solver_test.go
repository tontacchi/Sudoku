package sudoku

import (
	"testing"
	"sudoku-csp/solver"
	"time"
)

func TestSolver(t *testing.T) {
	const NUM_ROWS  = 10
	const NUM_COLS  = 10
	const NUM_HINTS = 2000

	board := NewForwardCheckGeneratedBoard(NUM_ROWS, NUM_COLS, NUM_HINTS)
	// board := NewRandomBoard(NUM_ROWS, NUM_COLS, NUM_HINTS)
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
		solver.NorvigCheck{},
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


// package sudoku
//
// import (
// 	"testing"
// 	"sudoku-csp/solver"
// 	"time"
// )
//
// func TestSolver(t *testing.T) {
// 	const NUM_ROWS  int = 10
// 	const NUM_COLS  int = 10
// 	const NUM_HINTS int = 2000
//
// 	board := NewRandomBoard(NUM_ROWS, NUM_COLS, NUM_HINTS)
// 	// board := NewForwardCheckGeneratedBoard(NUM_ROWS, NUM_COLS, NUM_HINTS)
// 	if board == nil {
// 		t.Fatal("NewRandomBoard returned nil")
// 	}
//
// 	const EXPECTED_BOARDLEN int = NUM_ROWS * NUM_COLS
// 	if board.boardLen != EXPECTED_BOARDLEN {
// 		t.Errorf("Expected boardLen to be %v, got %v\n", EXPECTED_BOARDLEN, board.boardLen)
// 		return
// 	}
//
// 	t.Log(board)
// 	network := NewNetworkFromBoard(board)
//
// 	solver := solver.NewBacktrackSolver(
// 		network,
// 		solver.NewTrail(),
// 		solver.MRVWithDegree{},
// 		solver.LeastConstrainingValue{},
// 		solver.NorvigCheck{},
// 	)
//
// 	solveRes := solver.Solve(time.Duration(1) * time.Minute)
//
// 	t.Log("\n")
// 	t.Log(solveRes)
// }
//
