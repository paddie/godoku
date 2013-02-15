// Package godoku is a simple brute-force
// in-place sudoku solver
package godoku

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Sudoku struct {
	board         Board
	solved        bool
	solutionCount int
	doPrint       bool
	dim           int
	solveAll      bool
	solution      Board
}

type Board [][]int

func (s *Sudoku) PrintBoard() {
	for _, row := range s.board {
		fmt.Println(row)
	}
	fmt.Println("")
}

// IsValidBoard iterates through all initial 
// values on the board and verifies that they indeed
// abide by the 3 laws of Sudoku
func (s *Sudoku) IsValidBoard() bool {
	if s.board == nil {
		return false
	}
	for i, row := range s.board {
		for j, val := range row {
			if val == 0 {
				continue
			}
			s.board[i][j] = 0
			if !s.ValidValueAtPosition(i, j, val) {
				s.board[i][j] = val
				return false
			}
			s.board[i][j] = val
		}
	}

	return true
}

// String returns either the unsolved board if the
// sudoku has not been solved, or the solution 
// if such a solution has been found 
// (by running one of the Solve* methods)
func (s *Sudoku) String() string {
	var buffer bytes.Buffer
	if s.solved {
		for _, row := range s.solution {
			buffer.WriteString(fmt.Sprintf("%v\n", row))
		}
		return buffer.String()
	}

	for _, row := range s.board {
		buffer.WriteString(fmt.Sprintf("%v\n", row))
	}
	return buffer.String()
}

// GetSolution returns the solution
// BUG(paddie): doesn't check if board is solved
func (s *Sudoku) GetSolution() Board {
	return s.solution
}

// Load a sudoku from a path and a dimension argument
func NewSudokuFromFile(path string, dim int) (*Sudoku, error) {
	s := new(Sudoku)
	var err error
	s.board, err = readBoardFromFile(path, dim)

	if err != nil {
		return nil, err
	}
	s.dim = dim

	return s, nil
}

// Loads a sudoku-board in a string-representation;
// The values are in a 9x9 matrix, using space " " as delimiters and '\n' as linebreaks
func NewSudokuFromString(path string, dim int) (*Sudoku, error) {
	s := new(Sudoku)
	var err error
	s.board, err = readBoardFromString(path, dim)

	if err != nil {
		return nil, err
	}
	s.dim = dim

	return s, nil
}

// Returns the number of solutions found. 
// Returns 0 if a Solve() call has not been made
// and if the Sudoku has no solutions. If at least one solution
// has been found, the number of solutions are returned
// (The number of solutions obviously vary depending if 
// Find() or FindAll() was used.
func (s *Sudoku) GetSolutionsCount() int {
	return s.solutionCount
}

// registers the first solutions in the s.solution
// board, and prints if doPrint is set.
func (s *Sudoku) registerSolution() {
	s.solutionCount++
	if s.doPrint {
		s.PrintBoard()
	}

	if s.solved {
		return
	}

	s.solved = true
	s.solution = make(Board, 9, 9)

	for i, row := range s.board {
		s.solution[i] = make([]int, 9, 9)
		copy(s.solution[i], row)
	}
}

// Check if the solver has found a solution
func (s *Sudoku) IsSolved() bool {
	return s.solved
}

// The dimensions of the sudoku board
func (s *Sudoku) Dimension() int {
	return s.dim
}

// Solve and save the solution. Returns an error if no Sudoku has been loaded
func (s *Sudoku) Solve() error {

	s.solved = false

	if s.board == nil {
		return fmt.Errorf("No Board has been loaded..")
	}

	s.solveAll = false

	s.bruteforcePosition(0, 0)

	return nil
}

// Same as Solve(), but this one also prints
// the solution to stdin
func (s *Sudoku) SolveAndPrint() error {
	s.doPrint = true

	err := s.Solve()

	s.doPrint = false

	return err
}

// Same as Solve, but keeps running until it has all
// the solutions and keeps a count. It only saves the first solution
func (s *Sudoku) SolveAll() error {

	s.solved = false

	if s.board == nil {
		return fmt.Errorf("No Board has been loaded..")
	}

	s.solveAll = true

	s.bruteforcePosition(0, 0)

	return nil
}

// Same as SolveAll but prints all the solutions
// to stdin
func (s *Sudoku) SolveAllAndPrint() error {
	s.doPrint = true

	err := s.SolveAll()

	s.doPrint = false

	return err
}

func (s *Sudoku) bruteforcePosition(row, col int) {
	// we use '0' to indicate a non-filled block
	if s.board[row][col] == 0 {
		for i := 1; i < 10; i++ {
			if s.ValidValueAtPosition(row, col, i) {
				// place the value and attempt to solve
				s.board[row][col] = i
				// attempt to solve the sudoku with placed value
				s.nextPosition(row, col)

				if s.solved && !s.solveAll {
					// if Solve() was used, we break
					// after first solution
					s.board[row][col] = 0
					return
				}

				// clean up after attempt
				s.board[row][col] = 0
			}
		}
	} else {
		s.nextPosition(row, col)
	}
}

// Does two things:
// 1) if the board is in a finished state, calls 
// registerSolution() and returns - enables
// bruteforcePostion to exhaust every remaining permutation
// 2) checks wether to move to next column or next row
func (s *Sudoku) nextPosition(row, col int) {
	// we run through the Board row by row
	// meaning we only change rows when we're in
	// the final column
	if col < 8 {
		s.bruteforcePosition(row, col+1)
	} else {
		// if we're in the final collumn in the final 
		// row; we have a solution
		// - else we iterate to next row and reset the collumn
		if row < 8 {
			s.bruteforcePosition(row+1, 0)
		} else {
			s.registerSolution()
		}
	}
}

// Verify that *val* can be legally placed at (row,col)
// given restrictions in column, row and 3x3 square
func (s *Sudoku) ValidValueAtPosition(row, col, val int) bool {
	if s.ValidInSquare(row, col, val) &&
		s.ValidInColumnAndRow(row, col, val) {
		// validInRow(row, val, Board) {
		return true
	}

	return false
}

// Checks that the *val* does not already occur in the
// active 3x3 square
func (s *Sudoku) ValidInSquare(row, col, val int) bool {
	row, col = int(row/3)*3, int(col/3)*3

	for i := row; i < row+3; i++ {
		for j := col; j < col+3; j++ {
			//fmt.Printf("row, col = %v, %v\n", i, j)
			if s.board[i][j] == val {
				return false
			}
		}
	}
	return true
}

// Checks if *val* already occurs in either the row or the column.
func (s *Sudoku) ValidInColumnAndRow(row, col, val int) bool {
	for i := 0; i < 9; i++ {
		if s.board[row][i] == val ||
			s.board[i][col] == val {
			return false
		}
	}
	return true
}

func readBoardFromFile(path string, dim int) (Board, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return readBoardFromString(string(content), dim)
}

func readBoardFromString(m string, dim int) (Board, error) {
	lines := strings.Split(m, "\n")

	if len(lines) < dim {
		return nil, fmt.Errorf("row count of input: %v does not match dim: %v", len(lines), dim)
	}

	Board := make(Board, dim, dim)

	for i := 0; i < dim; i++ {
		stringRows := strings.Split(lines[i], " ")

		if len(stringRows) < dim {
			return nil, fmt.Errorf("column count of input: %v does not match dim: %v", len(lines[i]), dim)
		}

		integerRow := make([]int, dim, dim)
		for j := 0; j < dim; j++ {
			str := stringRows[j]
			val, err := strconv.Atoi(str)
			if err != nil {
				return nil, err
			}
			integerRow[j] = val
		}
		Board[i] = integerRow
	}
	return Board, nil
}
