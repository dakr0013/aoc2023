package main

import (
	_ "embed"
	"log"
	"strings"
	"unicode"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed input.txt
var input string

func main() {
	exampleResult1 := part1(example2)
	if exampleResult1 != 4361 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 467835 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	lines := strings.Split(input, "\n")
	sum := 0
	for i, line := range lines {
		num := 0
		isPartNumber := false
		symbols := "!\"§$%&/()=?*+~#-@/"
		for j, char := range line {
			if unicode.IsDigit(char) {
				num *= 10
				num += int(char - '0')
				isValidIndex := func(index, length int) bool {
					return index >= 0 && index < length
				}
				isSymbolAt := func(lines []string, i, j int) bool {
					return isValidIndex(i, len(lines)) && isValidIndex(j, len(lines[i])) && strings.Contains(symbols, string(lines[i][j]))
				}
				if isSymbolAt(lines, i-1, j-1) ||
					isSymbolAt(lines, i-1, j) ||
					isSymbolAt(lines, i-1, j+1) ||
					isSymbolAt(lines, i, j-1) ||
					isSymbolAt(lines, i, j+1) ||
					isSymbolAt(lines, i+1, j-1) ||
					isSymbolAt(lines, i+1, j) ||
					isSymbolAt(lines, i+1, j+1) {
					isPartNumber = true
				}
			} else {
				if isPartNumber {
					sum += num
				}
				num = 0
				isPartNumber = false
			}
		}
		if isPartNumber {
			sum += num
		}
		num = 0
		isPartNumber = false
	}
	return sum
}

func part2(input string) int {
	lines := strings.Split(input, "\n")
	sum := 0

	gears := make([][]int, len(lines))
	for i := range gears {
		gears[i] = make([]int, len(lines[0]))
	}

	for i, line := range lines {
		num := 0
		isGear := false
		gearI := 0
		gearJ := 0
		for j, char := range line {
			if unicode.IsDigit(char) {
				num *= 10
				num += int(char - '0')
				isValidIndex := func(index, length int) bool {
					return index >= 0 && index < length
				}
				isGearAt := func(lines []string, i, j int) bool {
					if isValidIndex(i, len(lines)) && isValidIndex(j, len(lines[i])) && string(lines[i][j]) == "*" {
						gearI = i
						gearJ = j
						return true
					} else {
						return false
					}
				}
				if isGearAt(lines, i-1, j-1) ||
					isGearAt(lines, i-1, j) ||
					isGearAt(lines, i-1, j+1) ||
					isGearAt(lines, i, j-1) ||
					isGearAt(lines, i, j+1) ||
					isGearAt(lines, i+1, j-1) ||
					isGearAt(lines, i+1, j) ||
					isGearAt(lines, i+1, j+1) {
					isGear = true
				}
			} else {
				if isGear {
					if gears[gearI][gearJ] != 0 {
						sum += (num * gears[gearI][gearJ])
					} else {
						gears[gearI][gearJ] = num
					}
				}
				num = 0
				isGear = false
			}
		}
		if isGear {
			if gears[gearI][gearJ] != 0 {
				sum += (num * gears[gearI][gearJ])
			} else {
				gears[gearI][gearJ] = num
			}
		}
		num = 0
		isGear = false
	}
	return sum
}
