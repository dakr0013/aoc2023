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
	exampleResult1 := solve(example1, 1)
	if exampleResult1 != 94 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", solve(input, 1))

	exampleResult2 := solve(example2, 2)
	if exampleResult2 != 154 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", solve(input, 2))
}

func solve(input string, part int) int {
	ignoreSlopes := part == 2

	startTime := time.Now()
	lines := strings.Split(input, "\n")
	start := Vector{0, 1}
	end := Vector{len(lines) - 1, len(lines) - 2}
	graph := buildGraph(lines, start, end, ignoreSlopes)

	result, ok := longestDistance(graph.get(start), graph.get(end), map[*Node]bool{})
	println("part", part, ":", time.Since(startTime).String())
	if !ok {
		println("path not found")
	}
	return result
}

func longestDistance(current *Node, target *Node, visitedNodes map[*Node]bool) (int, bool) {
	if current == target {
		return 0, true
	}
	if visitedNodes[current] {
		return 0, false
	}

	newVisitedNodes := maps.Clone(visitedNodes)
	newVisitedNodes[current] = true

	maxDistance := 0
	isValid := false
	for next, dstCurrentToNext := range current.neighbors {
		dstNextToTarget, ok := longestDistance(next, target, newVisitedNodes)
		newDistance := dstCurrentToNext + dstNextToTarget
		if ok && newDistance >= maxDistance {
			maxDistance = newDistance
			isValid = true
		}
	}

	return maxDistance, isValid
}

type Node struct {
	pos        Vector
	discovered bool
	neighbors  map[*Node]int
}

type Nodes map[Vector]*Node

func (this Nodes) get(v Vector) *Node {
	if _, ok := this[v]; !ok {
		this[v] = &Node{
			pos:       v,
			neighbors: map[*Node]int{},
		}
	}
	return this[v]
}

func buildGraph(tiles []string, start Vector, end Vector, ignoreSlopes bool) Nodes {
	nodes := Nodes{}
	nodes.get(end).discovered = true

	toDiscover := []Position{{start, down}}
	for len(toDiscover) != 0 {
		curPos := toDiscover[0]
		toDiscover = toDiscover[1:]
		currentNode := nodes.get(curPos.position)
		currentNode.discovered = true

		directions := []Vector{
			curPos.direction,
			curPos.direction.turnLeft(),
			curPos.direction.turnRight(),
		}

		for _, direction := range directions {
			neighbor, distance, isValid := Position{curPos.position.move(direction), direction}.discoverNode(tiles, ignoreSlopes)
			if isValid {
				neighborNode := nodes.get(neighbor.position)
				currentNode.neighbors[neighborNode] = distance + 1
				if ignoreSlopes {
					neighborNode.neighbors[currentNode] = distance + 1
				}
				if !neighborNode.discovered {
					toDiscover = append(toDiscover, neighbor)
				}
			}
		}
	}
	return nodes
}

func (this Position) next(direction Vector, tiles []string) (Position, bool) {
	nextPos := this.position.move(direction)
	nextTile := tiles[nextPos.row][nextPos.col]
	if nextTile != '#' {
		return Position{nextPos, direction}, true
	}
	return Position{nextPos, direction}, false
}

func (this Position) discoverNode(tiles []string, ignoreSlopes bool) (Position, int, bool) {
	current := this
	distance := 0
	start := Vector{0, 1}
	end := Vector{len(tiles) - 1, len(tiles) - 2}

	for {
		if current.position == start || current.position == end {
			return current, distance, true
		}

		if current.position.row < 0 || current.position.row >= len(tiles) ||
			current.position.col < 0 || current.position.col >= len(tiles) {
			return current, distance, false
		}

		currentTile := tiles[current.position.row][current.position.col]
		if currentTile == '#' {
			return current, distance, false
		}
		if !ignoreSlopes && currentTile == '>' && current.direction != right {
			return current, distance, false
		}
		if !ignoreSlopes && currentTile == 'v' && current.direction != down {
			return current, distance, false
		}

		nextPositions := []Position{}
		if next, ok := current.next(current.direction.turnLeft(), tiles); ok {
			nextPositions = append(nextPositions, next)
		}
		if next, ok := current.next(current.direction.turnRight(), tiles); ok {
			nextPositions = append(nextPositions, next)
		}
		if next, ok := current.next(current.direction, tiles); ok {
			nextPositions = append(nextPositions, next)
		}

		validMovesCount := len(nextPositions)
		if validMovesCount >= 2 {
			return current, distance, true
		}
		if validMovesCount == 0 {
			return current, distance, false
		}

		distance++
		current = nextPositions[0]
	}
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
