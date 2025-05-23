package main

import (
	"sudoku-csp/solver"
	"sudoku-csp/sudoku"
	"time"
	"fmt"
)

func main() {
	const NUM_ROWS  = 3
	const NUM_COLS  = 2
	const NUM_HINTS = 10
	const SOLVE_TIME_LIMIT = time.Hour * 2

	board   := sudoku.NewRandomBoard(NUM_ROWS, NUM_COLS, NUM_HINTS)
	network := sudoku.NewNetworkFromBoard(board)
	trail   := solver.NewTrail()

	varSelector := solver.MRVWithDegree{}
	valSelector := solver.LeastConstrainingValue{}
	checker     := solver.NorvigCheck{}

	solver  := solver.NewBacktrackSolver(network, trail, varSelector, valSelector, checker)

	fmt.Printf("starting board:\n%v\n", board.String())

	start  := time.Now()
	result := solver.Solve(SOLVE_TIME_LIMIT)
	after  := time.Now()

	board = sudoku.NewBoardFromNetwork(solver.Network, NUM_ROWS, NUM_COLS)
	fmt.Printf("final board:\n%v\n", board.String())

	fmt.Printf("solution found: %v\n", result)
	fmt.Printf("solving time elapsed: %v\n", after.Sub(start))
}
