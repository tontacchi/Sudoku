package main

import (
	"sudoku-csp/solver"
	"sudoku-csp/sudoku"
	"time"
	"fmt"
)

func main() {
	const NUM_ROWS  = 4
	const NUM_COLS  = 4
	const NUM_HINTS = 120
	const SOLVE_TIME_LIMIT = time.Minute * 2

	fmt.Println("Creating sudoku board...")
	board   := sudoku.NewBoardFromSolved(NUM_ROWS, NUM_COLS, NUM_HINTS)
	network := sudoku.NewNetworkFromBoard(board)
	trail   := solver.NewTrail()
	
	varSelector := solver.MRV{}
	valSelector := solver.LeastConstrainingValue{}
	checker     := solver.ForwardChecking{}

	solver := solver.NewBacktrackSolver(network, trail, varSelector, valSelector, checker)

	fmt.Println("Done!")
	fmt.Printf("starting board:\n%v\n", board.String())

	start  := time.Now()
	result := solver.Solve(SOLVE_TIME_LIMIT)
	after  := time.Now()

	board = sudoku.NewBoardFromNetwork(network, NUM_ROWS, NUM_COLS)
	if solver.HasSolution {
		fmt.Printf("final board:\n%v\n", board.String())
	}

	fmt.Printf("solution found: %v\n", result)
	fmt.Printf("solving time elapsed: %v\n", after.Sub(start))
}
