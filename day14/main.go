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
	if exampleResult1 != 136 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 64 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	lines := strings.Split(input, "\n")
	farthestNorthRow := make([]int, len(lines[0]))
	sum := 0
	for i, line := range lines {
		for j, tile := range line {
			switch tile {
			case 'O':
				sum += len(lines) - farthestNorthRow[j]
				farthestNorthRow[j]++
			case '#':
				farthestNorthRow[j] = i + 1
			}
		}
	}
	return sum
}

func part2(input string) int {
	dish := parseDish(input)
	loopSize, loopStart, totalLoads := dish.findLoop(1000)
	cycles := (1000000000 - loopStart - 1) % loopSize
	return totalLoads[loopStart+cycles]
}

func (d *Dish) findLoop(maxCycles int) (loopSize, loopStart int, totalLoads []int) {
	totalLoads = make([]int, maxCycles)
	for i := 0; i < maxCycles; i++ {
		d.spinCycle()
		totalLoads[i] = d.totalLoad()
		loopSize, loopStart := findLoop(totalLoads, i+1)
		if loopSize > 0 {
			return loopSize, loopStart, totalLoads
		}
	}
	return -1, -1, totalLoads
}

func findLoop(numbers []int, length int) (loopSize, loopStart int) {
	for loopStart := 0; loopStart < length; loopStart++ {
		remaining := length - loopStart
		if remaining%2 == 0 && remaining >= 4 {
			loopSize := remaining / 2
			isLoop := true
			for i := loopStart; i < loopStart+loopSize; i++ {
				if numbers[i] != numbers[i+loopSize] {
					isLoop = false
					break
				}
			}
			if isLoop {
				return loopSize, loopStart
			}
		}
	}
	return -1, -1
}

func parseDish(s string) Dish {
	rows := strings.Split(s, "\n")
	positions := make([][]rune, len(rows))
	for i, row := range rows {
		positions[i] = make([]rune, len(row))
		for j, position := range row {
			positions[i][j] = position
		}
	}
	return Dish{positions: positions}
}

type Dish struct {
	positions [][]rune
}

func (d *Dish) String() string {
	builder := strings.Builder{}
	for _, row := range d.positions {
		for _, position := range row {
			builder.WriteRune(position)
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func (d *Dish) tiltNorth() {
	farthestNorthRow := make([]int, len(d.positions[0]))
	for i, row := range d.positions {
		for j, position := range row {
			switch position {
			case 'O':
				d.positions[i][j] = '.'
				d.positions[farthestNorthRow[j]][j] = 'O'
				farthestNorthRow[j]++
			case '#':
				farthestNorthRow[j] = i + 1
			}
		}
	}
}

func (d *Dish) tiltWest() {
	farthestWestCol := make([]int, len(d.positions))
	for j := 0; j < len(d.positions[0]); j++ {
		for i := 0; i < len(d.positions); i++ {
			switch d.positions[i][j] {
			case 'O':
				d.positions[i][j] = '.'
				d.positions[i][farthestWestCol[i]] = 'O'
				farthestWestCol[i]++
			case '#':
				farthestWestCol[i] = j + 1
			}
		}
	}
}

func (d *Dish) tiltSouth() {
	farthestSouthRow := make([]int, len(d.positions[0]))
	for i := range farthestSouthRow {
		farthestSouthRow[i] = len(d.positions) - 1
	}
	for i := len(d.positions) - 1; i >= 0; i-- {
		for j, position := range d.positions[i] {
			switch position {
			case 'O':
				d.positions[i][j] = '.'
				d.positions[farthestSouthRow[j]][j] = 'O'
				farthestSouthRow[j]--
			case '#':
				farthestSouthRow[j] = i - 1
			}
		}
	}
}

func (d *Dish) tiltEast() {
	farthestEastCol := make([]int, len(d.positions))
	for i := range farthestEastCol {
		farthestEastCol[i] = len(d.positions[0]) - 1
	}
	for j := len(d.positions[0]) - 1; j >= 0; j-- {
		for i := 0; i < len(d.positions); i++ {
			switch d.positions[i][j] {
			case 'O':
				d.positions[i][j] = '.'
				d.positions[i][farthestEastCol[i]] = 'O'
				farthestEastCol[i]--
			case '#':
				farthestEastCol[i] = j - 1
			}
		}
	}
}

func (d *Dish) spinCycle() {
	d.tiltNorth()
	d.tiltWest()
	d.tiltSouth()
	d.tiltEast()
}

func (d *Dish) totalLoad() int {
	sum := 0
	for i, row := range d.positions {
		for _, position := range row {
			if position == 'O' {
				sum += len(d.positions) - i
			}
		}
	}
	return sum
}
