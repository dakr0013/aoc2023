package main

import (
	_ "embed"
	"log"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed input.txt
var input string

func main() {
	exampleResult1 := part1(example1, 6)
	if exampleResult1 != 16 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input, 64))

	exampleResult2 := part2(example2)
	if exampleResult2 != 123 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string, steps int) int {
	startTime := time.Now()
	lines := strings.Split(input, "\n")
	start := findStart(lines)
	result := countGardenPlots(State{start, steps}, lines, make(map[Vector]bool))
	println(time.Since(startTime).String())
	return result
}

func part2(input string) int {
	lines := strings.Split(input, "\n")
	sum := len(lines)
	return sum
}

type State struct {
	curPos         Vector
	remainingSteps int
}

var cache map[State]int = make(map[State]int)

func countGardenPlots(s State, tiles []string, counted map[Vector]bool) int {
	if _, ok := cache[s]; ok {
		return 0
	}
	if tiles[s.curPos.row][s.curPos.col] == '#' {
		cache[s] = 0
		return 0
	}
	if s.remainingSteps == 0 {
		if counted[s.curPos] {
			return 0
		}
		counted[s.curPos] = true
		cache[s] = 1
		return 1
	}
	result := countGardenPlots(State{
		curPos:         s.curPos.move(right),
		remainingSteps: s.remainingSteps - 1,
	}, tiles, counted) +
		countGardenPlots(State{
			curPos:         s.curPos.move(down),
			remainingSteps: s.remainingSteps - 1,
		}, tiles, counted) +
		countGardenPlots(State{
			curPos:         s.curPos.move(left),
			remainingSteps: s.remainingSteps - 1,
		}, tiles, counted) +
		countGardenPlots(State{
			curPos:         s.curPos.move(up),
			remainingSteps: s.remainingSteps - 1,
		}, tiles, counted)

	cache[s] = result
	return result
}

func findStart(lines []string) Vector {
	var start Vector
	for i := range lines {
		for j := range lines[i] {
			if lines[i][j] == 'S' {
				start = Vector{i, j}
			}
		}
	}
	return start
}

type Vector struct {
	row, col int
}

func (v Vector) scalarMul(a int) Vector {
	return Vector{
		row: v.row * a,
		col: v.col * a,
	}
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

func (v Vector) isRightTurn(nextDirection Vector) bool {
	return v.turnRight() == nextDirection
}

func (v Vector) isLeftTurn(nextDirection Vector) bool {
	return v.turnLeft() == nextDirection
}

var directions []Vector = []Vector{up, right, down, left}
var left Vector = Vector{0, -1}
var right Vector = Vector{0, 1}
var up Vector = Vector{-1, 0}
var down Vector = Vector{1, 0}
