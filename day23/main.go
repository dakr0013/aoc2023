package main

import (
	_ "embed"
	"log"
	"strings"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed input.txt
var input string

func main() {
	exampleResult1 := part1(example1)
	if exampleResult1 != 94 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 154 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	startTime := time.Now()
	lines := strings.Split(input, "\n")
	result, ok := longestPath(Position{Vector{0, 1}, down}, Vector{len(lines) - 1, len(lines) - 2}, lines, map[Vector]bool{}, true)
	if !ok {
		log.Fatalln("no path found")
	}
	println("part1:", time.Since(startTime).String())
	return result
}

func part2(input string) int {
	startTime := time.Now()
	lines := strings.Split(input, "\n")
	result, ok := longestPath(Position{Vector{0, 1}, down}, Vector{len(lines) - 1, len(lines) - 2}, lines, map[Vector]bool{}, false)
	if !ok {
		log.Fatalln("no path found")
	}
	println("part2:", time.Since(startTime).String())
	return result
}

func longestPath(curPos Position, end Vector, tiles []string, visitedTiles map[Vector]bool, withSlopes bool) (int, bool) {
	if curPos.position == end {
		return 0, true
	}

	newVisitedTiles := maps.Clone(visitedTiles)
	newVisitedTiles[curPos.position] = true

	nextPositions := []Position{}
	if next, ok := nextPosition(curPos, curPos.direction.turnLeft(), tiles, visitedTiles, withSlopes); ok {
		nextPositions = append(nextPositions, next)
	}
	if next, ok := nextPosition(curPos, curPos.direction.turnRight(), tiles, visitedTiles, withSlopes); ok {
		nextPositions = append(nextPositions, next)
	}
	if next, ok := nextPosition(curPos, curPos.direction, tiles, visitedTiles, withSlopes); ok {
		nextPositions = append(nextPositions, next)
	}

	if len(nextPositions) == 0 {
		return 0, false
	}
	maxDistance, valid := 0, false
	for _, nextPos := range nextPositions {
		if distance, ok := longestPath(nextPos, end, tiles, newVisitedTiles, withSlopes); ok && distance >= maxDistance {
			maxDistance = distance
			valid = true
		}
	}

	return maxDistance + 1, valid

}

func nextPosition(curPos Position, direction Vector, tiles []string, visitedTiles map[Vector]bool, withSlopes bool) (Position, bool) {
	nextPos := curPos.position.move(direction)
	nextTile := tiles[nextPos.row][nextPos.col]
	if nextTile == '#' ||
		visitedTiles[nextPos] ||
		nextTile == '>' && direction != right && withSlopes ||
		nextTile == 'v' && direction != down && withSlopes {
		return curPos, false
	}

	return Position{nextPos, direction}, true
}

type Position struct {
	position  Vector
	direction Vector
}

type Vector struct {
	row, col int
}

func (v Vector) move(direction Vector) Vector {
	return Vector{v.row + direction.row, v.col + direction.col}
}

func (v Vector) turnRight() Vector {
	i := slices.Index(directions, v)
	return directions[(i+1)%len(directions)]
}

func (v Vector) turnLeft() Vector {
	i := slices.Index(directions, v)
	return directions[(i-1+len(directions))%len(directions)]
}

var directions []Vector = []Vector{up, right, down, left}
var left Vector = Vector{0, -1}
var right Vector = Vector{0, 1}
var up Vector = Vector{-1, 0}
var down Vector = Vector{1, 0}
