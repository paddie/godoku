package main

import (
	"io/ioutil"
	// "bufio"
	// "bytes"
	"fmt"
	// "io"
	// "os"
	"strconv"
	"strings"
)

func readMatrix(path string) [][]int {
	content, err := ioutil.ReadFile(path)

	// fmt.Printf("%#v", string(content))

	if err != nil {
		//Do something
	}
	lines := strings.Split(string(content), "\r")

	if len(lines) != 9 {
		panic(fmt.Sprintf("strings.Split: too many rows: %v", len(lines)))
	}

	matrix := make([][]int, 9, 9)
	for i, line := range lines {
		//fmt.Printf("%v: %v\n", i, line)

		stringRows := strings.Split(line, " ")

		ints := make([]int, 9)
		for j, str := range stringRows {
			val, err := strconv.Atoi(str)
			if err != nil {
				panic(err)
			}
			ints[j] = val
		}
		matrix[i] = ints
	}

	return matrix
}

func isValidPlacement(row, col, val int, matrix [][]int) bool {
	if validInSquare(row, col, val, matrix) &&
		validInColumnAndRow(row, col, val, matrix) {
		// validInRow(row, val, matrix) {
		return true
	}

	return false
}

func validInSquare(row, col, val int, matrix [][]int) bool {
	// fmt.Printf("validSqaure: row, col = %v, %v", row, col)
	row, col = int(row/3)*3, int(col/3)*3
	// fmt.Printf(" => row, col = %v, %v\n", row, col)

	for i := row; i < row+3; i++ {
		for j := col; j < col+3; j++ {
			//fmt.Printf("row, col = %v, %v\n", i, j)
			if matrix[i][j] == val {
				return false
			}
		}
	}
	return true
}

func validInColumnAndRow(row, col, val int, matrix [][]int) bool {
	for i := 0; i < 9; i++ {
		if matrix[row][i] == val ||
			matrix[i][col] == val {
			return false
		}
	}
	return true
}

func attemptToPlace(row, col int, matrix [][]int) {
	// ignore board values
	if matrix[row][col] != 0 {
		attemptNextPosition(row, col, matrix)

	} else {
		// Brute force: try to place values [1...9] on board
		for i := 1; i < 10; i++ {
			if isValidPlacement(row, col, i, matrix) {
				// place the value and attempt to solve
				matrix[row][col] = i
				// attempt to solve the sudoku with placed value
				attemptNextPosition(row, col, matrix)
				// clean up after attempt
				matrix[row][col] = 0
			}
		}
	}

}

func attemptNextPosition(row, col int, matrix [][]int) {
	if col < 8 {
		// fmt.Printf("1:at (%v,%v)\n", row, col+1)
		attemptToPlace(row, col+1, matrix)
	} else {
		if row < 8 {
			// fmt.Printf("2: at (%v,%v)\n", row+1, 0)
			attemptToPlace(row+1, 0, matrix)
		} else {
			printMatrix(matrix)
		}
	}
}

func printMatrix(matrix [][]int) {
	for _, row := range matrix {
		fmt.Println(row)
	}
	fmt.Println("")
}

func solveSudoku(matrix [][]int) {
	attemptToPlace(0, 0, matrix)
}

func main() {
	matrix := readMatrix("solvable88.txt")

	printMatrix(matrix)

	fmt.Println("")
	// fmt.Println("(3,3) : ")
	// for i := 1; i < 10; i++ {
	// 	fmt.Printf("[(%v=%v) ", i, isValidPlacement(3, 3, i, matrix))
	// }
	// fmt.Println("]")

	solveSudoku(matrix)
}
