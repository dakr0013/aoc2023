package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed input.txt
var input string

func main() {
	actual := []int{
		part1(example2, 6),
		part1(example2, 10),
		part1(example2, 50),
		part1(example2, 100),
	}
	expected := []int{
		16,
		50,
		1594,
		6536,
	}
	testPassed := true
	for i := range actual {
		if actual[i] != expected[i] {
			fmt.Printf("Result for steps: %d wrong;\nexpected: %d\nactual  : %d\n", (i+1)*10, expected[i], actual[i])
			testPassed = false
		}
	}
	if !testPassed {
		log.Fatalln("Some tests not passed")
	}

	log.Printf("Part 1: %d\n", part1(input, 64))
	log.Printf("Part 2: %d\n", part2(input, 26501365))
}

func part1(input string, steps int) int {
	cache = map[State]int{}
	startTime := time.Now()
	tiles := strings.Split(input, "\n")
	start := findStart(tiles)
	result := countGardenPlots(State{GlobalPosition{start, Vector{0, 0}}, steps}, tiles, make(map[GlobalPosition]bool))
	println("part1:", time.Since(startTime).String())
	return result
}

// only works for real input, not for example
func part2(input string, steps int) int {
	startTime := time.Now()
	tiles := strings.Split(input, "\n")
	mapSize := int64(len(tiles))

	// f(steps)->count function is a 2nd degree polynomial (drawing graph makes it clear),
	// which has a general form of f(x)=ax^2+bx+c
	// there are 3 variables: a,b,c
	// to solve those we need the results of 3 values for x
	x1 := int64(steps) % mapSize
	x2 := x1 + mapSize
	x3 := x2 + mapSize
	f_x1 := part1(input, int(x1))
	f_x2 := part1(input, int(x2))
	f_x3 := part1(input, int(x3))

	results, _ := solveLinearSystem([][]float64{
		{float64(x1 * x1), float64(x1), 1, float64(f_x1)},
		{float64(x2 * x2), float64(x2), 1, float64(f_x2)},
		{float64(x3 * x3), float64(x3), 1, float64(f_x3)},
	})

	a := results[0]
	b := results[1]
	c := results[2]

	f := func(x float64) float64 {
		return a*x*x + b*x + c
	}
	result := int(f(float64(steps)))
	println("part2:", time.Since(startTime).String())
	return result
}

func solveLinearSystem(A [][]float64) ([]float64, error) {
	n := len(A)

	// Step 2: Apply Gauss Elimination on Matrix A
	for i := 0; i < n-1; i++ {
		if A[i][i] == 0 {
			return nil, errors.New("Cannot solve")
		}

		for j := i + 1; j < n; j++ {
			ratio := A[j][i] / A[i][i]
			for k := 0; k < n+1; k++ {
				A[j][k] = A[j][k] - ratio*A[i][k]
			}
		}
	}

	// Step 3: Obtaining Solution by Back Substitution
	var X = make([]float64, n)
	X[n-1] = A[n-1][n] / A[n-1][n-1]
	for i := n - 2; i >= 0; i-- {
		X[i] = A[i][n]
		for j := i + 1; j < n; j++ {
			X[i] = X[i] - A[i][j]*X[j]
		}
		X[i] = X[i] / A[i][i]
	}

	return X, nil
}

type State struct {
	curPos         GlobalPosition
	remainingSteps int
}

type GlobalPosition struct {
	local  Vector // represents position within grid
	global Vector // represents on which grid it is (0,0) is the start grid, (1,0) is repeated grid below
}

func (this GlobalPosition) move(direction Vector, maxSize int) GlobalPosition {
	newLocal := this.local.move(direction)
	newGlobal := this.global
	if newLocal.isOutside(maxSize) {
		newGlobal = newGlobal.move(direction)
		newLocal = newLocal.mod(maxSize)
	}
	return GlobalPosition{
		local:  newLocal,
		global: newGlobal,
	}
}

var cache map[State]int = make(map[State]int)

func countGardenPlots(s State, tiles []string, counted map[GlobalPosition]bool) int {
	if _, ok := cache[s]; ok {
		return 0
	}
	if tiles[s.curPos.local.row][s.curPos.local.col] == '#' {
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
	maxSize := len(tiles)
	result := countGardenPlots(State{
		curPos:         s.curPos.move(right, maxSize),
		remainingSteps: s.remainingSteps - 1,
	}, tiles, counted) +
		countGardenPlots(State{
			curPos:         s.curPos.move(down, maxSize),
			remainingSteps: s.remainingSteps - 1,
		}, tiles, counted) +
		countGardenPlots(State{
			curPos:         s.curPos.move(left, maxSize),
			remainingSteps: s.remainingSteps - 1,
		}, tiles, counted) +
		countGardenPlots(State{
			curPos:         s.curPos.move(up, maxSize),
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

func (v Vector) move(direction Vector) Vector {
	return Vector{v.row + direction.row, v.col + direction.col}
}

func (this Vector) mod(maxSize int) Vector {
	return Vector{
		row: (this.row + maxSize) % maxSize,
		col: (this.col + maxSize) % maxSize,
	}
}

func (this Vector) isOutside(maxSize int) bool {
	return this.row < 0 || this.col < 0 || this.col >= maxSize || this.row >= maxSize
}

var left Vector = Vector{0, -1}
var right Vector = Vector{0, 1}
var up Vector = Vector{-1, 0}
var down Vector = Vector{1, 0}
