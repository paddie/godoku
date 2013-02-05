package godoku

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Sudoku struct {
	matrix        Matrix
	solved        bool
	solutionCount int
	print         bool
	dim           int
}
type Matrix [][]int

func (s *Sudoku) PrintMatrix() {
	for _, row := range s.matrix {
		fmt.Println(row)
	}
	fmt.Println("")
}

func (s *Sudoku) IsValidBoard() (bool, err) {
	if s.matrix == nil {
		return false, fmt.Errorf("Matrix has not been initialized")
	}

	for i, row := range s.matrix {
		for j, val := range row {
			if !s.ValidValueAtPosition(i, j, val) {
				return false, nil
			}
		}
	}

	return true, nil
}

func (s *Sudoku) String() string {
	var buffer bytes.Buffer
	for _, row := range s.matrix {
		buffer.WriteString(fmt.Sprintf("%v\n", row))
	}
	buffer.WriteString("\n")
	return buffer.String()
}

func NewSudokuFromFile(path string, dim int, doPrint bool) (*Sudoku, error) {
	s := new(Sudoku)
	var err error
	s.matrix, err = readMatrixFromFile(path, dim)

	if err != nil {
		return nil, err
	}

	s.dim = dim
	s.print = doPrint

	return s, nil
}

// Assumes a 9x9 Sudoku board
func NewSudokuFromString(path string, dim int, doPrint bool) (*Sudoku, error) {
	s := new(Sudoku)
	var err error
	s.matrix, err = readMatrixFromString(path, dim)

	if err != nil {
		return nil, err
	}

	s.dim = dim
	s.print = doPrint

	return s, nil
}

func (s *Sudoku) GetSolutionsCount() int {
	return s.solutionCount
}

// Could potentially make a copy of the matrix at this point
// to preserve the solution for further processing
func (s *Sudoku) registerSolution() {
	s.solutionCount++
	if s.print {
		s.PrintMatrix()
	}
	if !s.solved {
		s.solved = true
	}
}

func (s *Sudoku) IsSolved() bool {
	return s.solved
}

func (s *Sudoku) Dimension() int {
	return s.dim
}

func (s *Sudoku) Solve() error {

	s.solved = false

	if s.matrix == nil {
		return fmt.Errorf("No matrix has been loaded..")
	}
	s.bruteforcePosition(0, 0)

	return nil
}

func (s *Sudoku) bruteforcePosition(row, col int) {
	// we use '0' to indicate a non-filled block
	if s.matrix[row][col] == 0 {
		for i := 1; i < 10; i++ {
			if s.ValidValueAtPosition(row, col, i) {
				// place the value and attempt to solve
				s.matrix[row][col] = i
				// attempt to solve the sudoku with placed value
				s.nextPosition(row, col)
				// clean up after attempt
				s.matrix[row][col] = 0
			}
		}
	} else {
		s.nextPosition(row, col)
	}
}

// Does two things:
//	1) if the board is in a finished state, calls 
//		registerSolution() and returns - enables
//		bruteforcePostion to exhaust every remaining permutation
//	2) checks wether to move to next column or next row
func (s *Sudoku) nextPosition(row, col int) {
	// we run through the matrix row by row
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

// Verify that 'val' can be legally placed at (row,col)
// given restrictions in column, row and 3x3 square
func (s *Sudoku) ValidValueAtPosition(row, col, val int) bool {
	if s.ValidInSquare(row, col, val) &&
		s.ValidInColumnAndRow(row, col, val) {
		// validInRow(row, val, matrix) {
		return true
	}

	return false
}

// Checks that the 'val' does not already occur in the
// active 3x3 square.
// TODO: make square validate sudoku boards of random size
func (s *Sudoku) ValidInSquare(row, col, val int) bool {
	row, col = int(row/3)*3, int(col/3)*3

	for i := row; i < row+3; i++ {
		for j := col; j < col+3; j++ {
			//fmt.Printf("row, col = %v, %v\n", i, j)
			if s.matrix[i][j] == val {
				return false
			}
		}
	}
	return true
}

// Checks if 'val' already occurs in either the row or the column.
func (s *Sudoku) ValidInColumnAndRow(row, col, val int) bool {
	for i := 0; i < 9; i++ {
		if s.matrix[row][i] == val ||
			s.matrix[i][col] == val {
			return false
		}
	}
	return true
}

func readMatrixFromFile(path string, dim int) (Matrix, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return readMatrixFromString(string(content), dim)
}

func readMatrixFromString(m string, dim int) (Matrix, error) {
	lines := strings.Split(m, "\n")
	matrix := make(Matrix, dim, dim)

	for i := 0; i < dim; i++ {
		stringRows := strings.Split(lines[i], " ")

		integerRow := make([]int, dim, dim)
		for j := 0; j < dim; j++ {
			str := stringRows[j]
			val, err := strconv.Atoi(str)
			if err != nil {
				return nil, err
			}
			integerRow[j] = val
		}
		matrix[i] = integerRow
	}
	return matrix, nil
}
