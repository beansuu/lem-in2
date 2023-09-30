package main

import (
	"fmt"
)

func main() {
	// var TwoDArray [8][8]int
	// TwoDArray[3][6] = 18
	// TwoDArray[7][4] = 3
	// fmt.Println(TwoDArray)

	var rows int
	var cols int
	rows = 7
	cols = 9
	var twodslices = make([][]int, rows)
	var i int
	for i = range twodslices {
		twodslices[i] = make([]int, cols)
	}
	fmt.Println(twodslices)
}
