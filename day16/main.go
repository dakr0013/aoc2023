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
	if exampleResult1 != 46 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 51 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	contraption := Contraption{strings.Split(input, "\n")}
	return contraption.energizedCount(Vector{0, 0}, right)
}

func part2(input string) int {
	contraption := Contraption{strings.Split(input, "\n")}
	maxEnergized := 0
	lastRow := len(contraption.tiles) - 1
	lastCol := len(contraption.tiles[0]) - 1
	for col := range contraption.tiles[0] {
		result1 := contraption.energizedCount(Vector{0, col}, down)
		result2 := contraption.energizedCount(Vector{lastRow, col}, up)
		if result1 > maxEnergized {
			maxEnergized = result1
		}
		if result2 > maxEnergized {
			maxEnergized = result2
		}
	}
	for row := range contraption.tiles {
		result1 := contraption.energizedCount(Vector{row, 0}, right)
		result2 := contraption.energizedCount(Vector{row, lastCol}, left)
		if result1 > maxEnergized {
			maxEnergized = result1
		}
		if result2 > maxEnergized {
			maxEnergized = result2
		}
	}
	return maxEnergized
}

type Contraption struct {
	tiles []string
}

func (c Contraption) energizedCount(startingPos, direction Vector) int {
	energizedTiles := make([][]bool, len(c.tiles))
	for i := range energizedTiles {
		energizedTiles[i] = make([]bool, len(c.tiles[0]))
	}

	c.determineEnergizedTiles(energizedTiles, startingPos, direction)

	count := 0
	for i := range energizedTiles {
		for j := range energizedTiles[i] {
			if energizedTiles[i][j] {
				count++
			}
		}
	}
	return count
}

func (c Contraption) determineEnergizedTiles(energizedTiles [][]bool, curPos, direction Vector) {
	if curPos.row < 0 ||
		curPos.row >= len(c.tiles) ||
		curPos.col < 0 ||
		curPos.col >= len(c.tiles[0]) {
		return
	}

	switch c.tiles[curPos.row][curPos.col] {
	case '.':
		energizedTiles[curPos.row][curPos.col] = true
		c.determineEnergizedTiles(energizedTiles, curPos.move(direction), direction)
	case '-':
		if energizedTiles[curPos.row][curPos.col] {
			return
		}
		energizedTiles[curPos.row][curPos.col] = true
		if direction.isHorizontal() {
			c.determineEnergizedTiles(energizedTiles, curPos.move(direction), direction)
		} else {
			c.determineEnergizedTiles(energizedTiles, curPos.move(left), left)
			c.determineEnergizedTiles(energizedTiles, curPos.move(right), right)
		}
	case '|':
		if energizedTiles[curPos.row][curPos.col] {
			return
		}
		energizedTiles[curPos.row][curPos.col] = true
		if !direction.isHorizontal() {
			c.determineEnergizedTiles(energizedTiles, curPos.move(direction), direction)
		} else {
			c.determineEnergizedTiles(energizedTiles, curPos.move(up), up)
			c.determineEnergizedTiles(energizedTiles, curPos.move(down), down)
		}
	case '\\':
		energizedTiles[curPos.row][curPos.col] = true
		switch {
		case direction.equals(left):
			c.determineEnergizedTiles(energizedTiles, curPos.move(up), up)
		case direction.equals(right):
			c.determineEnergizedTiles(energizedTiles, curPos.move(down), down)
		case direction.equals(up):
			c.determineEnergizedTiles(energizedTiles, curPos.move(left), left)
		case direction.equals(down):
			c.determineEnergizedTiles(energizedTiles, curPos.move(right), right)
		}
	case '/':
		energizedTiles[curPos.row][curPos.col] = true
		switch {
		case direction.equals(left):
			c.determineEnergizedTiles(energizedTiles, curPos.move(down), down)
		case direction.equals(right):
			c.determineEnergizedTiles(energizedTiles, curPos.move(up), up)
		case direction.equals(up):
			c.determineEnergizedTiles(energizedTiles, curPos.move(right), right)
		case direction.equals(down):
			c.determineEnergizedTiles(energizedTiles, curPos.move(left), left)
		}
	}
}

type Vector struct {
	row, col int
}

func (v Vector) equals(other Vector) bool {
	return v.row == other.row && v.col == other.col
}

func (v Vector) move(direction Vector) Vector {
	return Vector{v.row + direction.row, v.col + direction.col}
}

func (v Vector) isHorizontal() bool {
	return v.row == 0
}

var left Vector = Vector{0, -1}
var right Vector = Vector{0, 1}
var up Vector = Vector{-1, 0}
var down Vector = Vector{1, 0}
