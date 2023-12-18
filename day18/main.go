package main

import (
	_ "embed"
	"log"
	"strconv"
	"strings"

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
	if exampleResult1 != 62 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 952408144115 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	digPlan := parseDigPlan(input)
	return digPlan.lagoonArea()
}

func part2(input string) int {
	digPlan := parseDigPlan2(input)
	return digPlan.lagoonArea()
}

func parseDigPlan(s string) DigPlan {
	lines := strings.Split(s, "\n")
	instructions := make([]Instruction, len(lines))
	for i, line := range lines {
		instructions[i] = parseInstruction(line)
	}
	return DigPlan{instructions}
}

func parseInstruction(s string) Instruction {
	direction := strings.Split(s, " ")[0]
	distance, _ := strconv.Atoi(strings.Split(s, " ")[1])
	directionVector := Vector{}
	switch direction {
	case "R":
		directionVector = right
	case "L":
		directionVector = left
	case "U":
		directionVector = up
	case "D":
		directionVector = down
	default:
		log.Fatalln("error")
	}
	return Instruction{
		direction: directionVector,
		distance:  distance,
	}
}

func parseDigPlan2(s string) DigPlan {
	lines := strings.Split(s, "\n")
	instructions := make([]Instruction, len(lines))
	for i, line := range lines {
		instructions[i] = parseInstruction2(line)
	}
	return DigPlan{instructions}
}

func parseInstruction2(s string) Instruction {
	rawInstruction := strings.Trim(strings.Split(s, " ")[2], "(#)")
	direction := rawInstruction[len(rawInstruction)-1]
	distance, _ := strconv.ParseInt(rawInstruction[:len(rawInstruction)-1], 16, 0)
	directionVector := Vector{}
	switch direction {
	case '0':
		directionVector = right
	case '1':
		directionVector = down
	case '2':
		directionVector = left
	case '3':
		directionVector = up
	default:
		log.Fatalln("error")
	}
	return Instruction{
		direction: directionVector,
		distance:  int(distance),
	}
}

type DigPlan struct {
	instructions []Instruction
}

func (d DigPlan) lagoonArea() int {
	area := 0
	curPos := Vector{0, 0}
	for i, instruction := range d.instructions {
		previousDirection := d.instructions[(i-1+len(d.instructions))%len(d.instructions)].direction
		currentDirection := instruction.direction
		nextDirection := d.instructions[(i+1)%len(d.instructions)].direction
		distance := instruction.distance
		// if edge (#) is like this (two right turns)
		// >########
		//  inside # outside
		// <########
		// then edge is actually 1 longer then distance (so each # is accounted for in area calculation)
		if previousDirection.isRightTurn(currentDirection) && currentDirection.isRightTurn(nextDirection) {
			distance++
		}

		// check if edge (#) is like this (two left turns)
		//         ########<
		//  inside # outside
		//         ########>
		// then edge is actually 1 shorter then distance because tiles outside are 1 less then distance)
		if previousDirection.isLeftTurn(currentDirection) && currentDirection.isLeftTurn(nextDirection) {
			distance--
		}
		nextPos := curPos.move(instruction.direction.scalarMul(distance))
		area += ((curPos.row + nextPos.row) * (curPos.col - nextPos.col))
		curPos = nextPos
	}
	if area < 0 {
		return (-area / 2)
	}
	return (area / 2)
}

type Instruction struct {
	direction Vector
	distance  int
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
