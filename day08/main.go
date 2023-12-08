package main

import (
	_ "embed"
	"log"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed input.txt
var input string

func main() {
	exampleResult1 := part1(example1)
	if exampleResult1 != 2 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 6 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))
}

func part1(input string) int {
	desertMap := parseMap(input)
	return desertMap.stepCount("AAA", "ZZZ")
}

func part2(input string) int {
	desertMap := parseMap(input)
	startingNodes := []*Node{}
	for label, node := range desertMap.nodes {
		if strings.HasSuffix(label, "A") {
			startingNodes = append(startingNodes, node)
		}
	}
	individualStepCounts := make(map[int]bool)
	for _, node := range startingNodes {
		stepCount := desertMap.stepCount(node.label, "Z")
		individualStepCounts[stepCount] = true
	}
	ggtOfSteps := ggtAll(maps.Keys(individualStepCounts))
	result := 1
	for stepCount := range individualStepCounts {
		result *= stepCount
		result /= ggtOfSteps
	}
	return result * ggtOfSteps
}

func ggtAll(a []int) int {
	if len(a) == 1 {
		return a[0]
	}
	overallGgt := ggt(a[0], a[1])
	for _, num := range a[2:] {
		overallGgt = ggt(overallGgt, num)
	}
	return overallGgt
}

func ggt(a, b int) int {
	if a == 0 {
		return b
	}
	for b != 0 {
		if a > b {
			a = a - b
		} else {
			b = b - a
		}
	}
	return a
}

func parseMap(s string) Map {
	lines := strings.Split(s, "\n")
	instructions := lines[0]
	nodeStrings := lines[2:]
	nodes := make(map[string]*Node)
	re := regexp.MustCompile("(...) = \\((...), (...)\\)")
	for _, nodeString := range nodeStrings {
		matches := re.FindStringSubmatch(nodeString)
		nodeLabel := matches[1]
		leftChildLabel := matches[2]
		rightChildLabel := matches[3]

		leftChild, ok := nodes[leftChildLabel]
		if !ok {
			nodes[leftChildLabel] = &Node{
				label: leftChildLabel,
			}
			leftChild = nodes[leftChildLabel]
		}
		rightChild, ok := nodes[rightChildLabel]
		if !ok {
			nodes[rightChildLabel] = &Node{
				label: rightChildLabel,
			}
			rightChild = nodes[rightChildLabel]
		}

		_, ok = nodes[nodeLabel]
		if !ok {
			nodes[nodeLabel] = &Node{
				label:      nodeLabel,
				leftChild:  leftChild,
				rightChild: rightChild,
			}
		} else {
			nodes[nodeLabel].leftChild = leftChild
			nodes[nodeLabel].rightChild = rightChild
		}
	}

	return Map{
		instructions: instructions,
		nodes:        nodes,
	}
}

type Map struct {
	instructions string
	nodes        map[string]*Node
}

func (m Map) stepCount(startLabel, targetSuffix string) int {
	currentNode := m.nodes[startLabel]
	stepCount := 0
	for {
		if strings.HasSuffix(currentNode.label, targetSuffix) {
			return stepCount
		}
		nextMove := m.instructions[stepCount%len(m.instructions)]
		if nextMove == 'L' {
			currentNode = currentNode.leftChild
		} else {
			currentNode = currentNode.rightChild
		}
		stepCount++
	}
}

type Node struct {
	label      string
	leftChild  *Node
	rightChild *Node
}
