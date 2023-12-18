package main

import (
	_ "embed"
	"log"
	"math"
	"slices"
	"strconv"
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
	exampleResult1 := part1(example1)
	if exampleResult1 != 102 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 71 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))
}

func part1(input string) int {
	startTime := time.Now()
	heatLoss := parseHeatLoss(input)
	start := Node{Vector{0, 0}, right, 0}
	target := Vector{len(heatLoss) - 1, len(heatLoss[0]) - 1}
	result, found := aStar(start, func(n Node) bool {
		return n.pos == target
	}, func(n Node) []Node {
		return n.successors(len(heatLoss)-1, len(heatLoss[0])-1)
	}, func(n Node) int {
		return heatLoss[n.pos.row][n.pos.col]
	})
	if !found {
		log.Fatalln("Path not found")
	}
	elapsedTime := time.Since(startTime)
	println("part1:", elapsedTime.String())
	return result
}

func part2(input string) int {
	startTime := time.Now()
	heatLoss := parseHeatLoss(input)
	start := Node{Vector{0, 0}, right, 0}
	target := Vector{len(heatLoss) - 1, len(heatLoss[0]) - 1}
	result, found := aStar(start, func(n Node) bool {
		return n.pos == target && n.straightCount >= 4
	}, func(n Node) []Node {
		return n.ultraSuccessors(len(heatLoss)-1, len(heatLoss[0])-1)
	}, func(n Node) int {
		return heatLoss[n.pos.row][n.pos.col]
	})
	if !found {
		log.Fatalln("Path not found")
	}
	elapsedTime := time.Since(startTime)
	println("part2:", elapsedTime.String())
	return result
}

type successorsFunc func(n Node) []Node

type weigthFunc func(n Node) int

type isTargetFunc func(n Node) bool

func aStar(start Node, isTarget isTargetFunc, successors successorsFunc, weight weigthFunc) (result int, found bool) {
	g := make(map[Node]int)
	g[start] = 0

	f := func(n Node) int {
		if curG, ok := g[n]; ok {
			return curG
		}
		return math.MaxInt
	}

	openList := []Node{start}

	closedList := make(map[Node]bool, 0)

	for len(openList) != 0 {
		slices.SortFunc(openList, func(a Node, b Node) int {
			return f(a) - f(b)
		})
		currentNode := openList[0]
		openList = openList[1:]
		if isTarget(currentNode) {
			return g[currentNode], true
		}
		closedList[currentNode] = true

		for _, successor := range successors(currentNode) {
			if _, ok := closedList[successor]; ok {
				continue
			}
			tentativeG := g[currentNode] + weight(successor)
			if curG, ok := g[successor]; !ok || tentativeG < curG {
				g[successor] = tentativeG
			}
			if !slices.Contains(openList, successor) {
				openList = append(openList, successor)
			}
		}
	}
	return -1, false
}

type Node struct {
	pos           Vector
	direction     Vector
	straightCount int
}

func (n Node) ultraSuccessors(maxRow, maxCol int) []Node {
	successors := make([]Node, 0)
	if n.straightCount < 10 {
		straightNode := Node{
			pos:           n.pos.move(n.direction),
			direction:     n.direction,
			straightCount: n.straightCount + 1,
		}
		if straightNode.isValid(maxRow, maxCol) {
			successors = append(successors, straightNode)
		}
	}
	if n.straightCount >= 4 || n.straightCount == 0 {
		left := n.direction.turnLeft()
		leftNode := Node{
			pos:           n.pos.move(left),
			direction:     left,
			straightCount: 1,
		}
		if leftNode.isValid(maxRow, maxCol) {
			successors = append(successors, leftNode)
		}
		right := n.direction.turnRight()
		rightNode := Node{
			pos:           n.pos.move(right),
			direction:     right,
			straightCount: 1,
		}
		if rightNode.isValid(maxRow, maxCol) {
			successors = append(successors, rightNode)
		}
	}
	return successors
}

func (n Node) successors(maxRow, maxCol int) []Node {
	successors := make([]Node, 0)
	if n.straightCount < 3 {
		straightNode := Node{
			pos:           n.pos.move(n.direction),
			direction:     n.direction,
			straightCount: n.straightCount + 1,
		}
		if straightNode.isValid(maxRow, maxCol) {
			successors = append(successors, straightNode)
		}
	}
	left := n.direction.turnLeft()
	leftNode := Node{
		pos:           n.pos.move(left),
		direction:     left,
		straightCount: 1,
	}
	if leftNode.isValid(maxRow, maxCol) {
		successors = append(successors, leftNode)
	}
	right := n.direction.turnRight()
	rightNode := Node{
		pos:           n.pos.move(right),
		direction:     right,
		straightCount: 1,
	}
	if rightNode.isValid(maxRow, maxCol) {
		successors = append(successors, rightNode)
	}
	return successors
}

func (n Node) isValid(maxRow, maxCol int) bool {
	return n.pos.row >= 0 && n.pos.row <= maxRow &&
		n.pos.col >= 0 && n.pos.col <= maxCol
}

func parseHeatLoss(s string) [][]int {
	lines := strings.Split(s, "\n")
	result := make([][]int, len(lines))
	for i := range result {
		result[i] = make([]int, len(lines[0]))
		for j := range result[i] {
			heatLoss, _ := strconv.Atoi(string(lines[i][j]))
			result[i][j] = heatLoss
		}
	}
	return result
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
