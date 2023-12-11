package main

import (
	_ "embed"
	"log"
	"strings"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed input.txt
var input string

func main() {
	exampleResult1 := part1(example1)
	if exampleResult1 != 374 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2, 100)
	if exampleResult2 != 8410 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input, 1000000))

}

func part1(input string) int {
	universe := strings.Split(input, "\n")
	return sumOfDistances(universe, 1)
}

func part2(input string, expansionSize int) int {
	universe := strings.Split(input, "\n")
	return sumOfDistances(universe, expansionSize)
}

func sumOfDistances(universe []string, expansionSize int) int {
	rowShouldExpand := make([]bool, len(universe))
	colShouldExpand := make([]bool, len(universe[0]))
	for i := range rowShouldExpand {
		rowShouldExpand[i] = true
	}
	for i := range colShouldExpand {
		colShouldExpand[i] = true
	}

	galaxies := make([]Galaxy, 0)
	for i := range universe {
		for j := range universe {
			if universe[i][j] == '#' {
				rowShouldExpand[i] = false
				colShouldExpand[j] = false
				galaxies = append(galaxies, Galaxy{row: i, col: j})
			}
		}
	}

	sum := 0
	for i, galaxy1 := range galaxies[:len(galaxies)-1] {
		for _, galaxy2 := range galaxies[i+1:] {
			rowDiff := abs(galaxy2.row-galaxy1.row) + countExpansions(rowShouldExpand, galaxy1.row, galaxy2.row)*expansionSize
			colDiff := abs(galaxy2.col-galaxy1.col) + countExpansions(colShouldExpand, galaxy1.col, galaxy2.col)*expansionSize
			if expansionSize > 1 {
				rowDiff -= countExpansions(rowShouldExpand, galaxy1.row, galaxy2.row)
				colDiff -= countExpansions(colShouldExpand, galaxy1.col, galaxy2.col)
			}
			distance := rowDiff + colDiff
			sum += distance
		}
	}

	return sum
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func countExpansions(shouldExpand []bool, from, to int) int {
	if from > to {
		from, to = to, from
	}
	count := 0
	for i := from; i < to; i++ {
		if shouldExpand[i] {
			count++
		}
	}
	return count
}

type Galaxy struct {
	row int
	col int
}
