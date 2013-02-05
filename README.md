godoku
======

Sudoku solver in Golang

======
Example:

```go

package main

import "github.com/paddie/godoku"

const solvable88 string = `0 7 0 0 0 0 0 8 0
0 3 0 7 6 2 0 0 1
0 0 1 9 8 0 0 0 0
1 0 0 0 0 0 0 0 0
8 0 3 0 0 0 0 0 2
0 0 6 0 0 0 0 0 8
0 0 0 0 3 1 6 0 0
5 0 0 2 4 9 0 1 0
0 1 0 0 0 0 0 9 0`

func main() {
  s, err := NewSudokuFromString(solvable88, 9, false)
  
  if err != nil {
		t.Error(err)
	}
	
  s.SolveAndPrint()

	if s.GetSolutionsCount() != 1 {
		fmt.Printf("Expected 1 != Actual %v\n", s.GetSolutionsCount())
	}
}
