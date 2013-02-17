godoku
======

Sudoku solver in Golang.

See [Documentation](http://godoc.org/github.com/paddie/godoku) on [GoDoc.org](http://godoc.org/).

Example:
=======

```go
package main

import (
	"fmt"
	"github.com/paddie/godoku"
)

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
	s, err := godoku.NewSudokuFromString(solvable88, 9)
	if err != nil {
		fmt.Println(err)
		return
	}

	// check that board is valid
	if !s.IsValidBoard() {
		fmt.Println("Invalid board")
		return
	}

	fmt.Println(s)

	// solve the board
	s.Solve()
	
	// print the solution
	fmt.Println(s)

	// checks the number of solutions; 1 in this case
	if s.GetSolutionsCount() != 1 {
		fmt.Printf("Expected 1 != Actual %v\n", s.GetSolutionsCount())
		return
	}

	fmt.Println("This sudoku has one solution!\n")
}
```
Output:
=======

```
[0 7 0 0 0 0 0 8 0]
[0 3 0 7 6 2 0 0 1]
[0 0 1 9 8 0 0 0 0]
[1 0 0 0 0 0 0 0 0]
[8 0 3 0 0 0 0 0 2]
[0 0 6 0 0 0 0 0 8]
[0 0 0 0 3 1 6 0 0]
[5 0 0 2 4 9 0 1 0]
[0 1 0 0 0 0 0 9 0]

[2 7 9 1 5 3 4 8 6]
[4 3 8 7 6 2 9 5 1]
[6 5 1 9 8 4 3 2 7]
[1 4 5 3 2 8 7 6 9]
[8 9 3 6 1 7 5 4 2]
[7 2 6 4 9 5 1 3 8]
[9 8 2 5 3 1 6 7 4]
[5 6 7 2 4 9 8 1 3]
[3 1 4 8 7 6 2 9 5]
	
This sudoku has one solution!
```
