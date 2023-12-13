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
	if exampleResult1 != 405 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 400 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))
}

func part1(input string) int {
	patterns := parsePatterns(input)
	sum := 0
	for _, pattern := range patterns {
		sum += pattern.calcReflection(0)
	}
	return sum
}

func part2(input string) int {
	patterns := parsePatterns(input)
	sum := 0
	for _, pattern := range patterns {
		sum += pattern.calcReflection(1)
	}
	return sum
}

func parsePatterns(s string) []Pattern {
	patternsString := strings.Split(s, "\n\n")
	patterns := make([]Pattern, len(patternsString))
	for i := range patterns {
		patterns[i] = Pattern{
			lines: strings.Split(patternsString[i], "\n"),
		}
	}
	return patterns
}

type Pattern struct {
	lines []string
}

func (p Pattern) calcReflection(smudges int) int {
	if colToLeft := p.verticalReflection(smudges); colToLeft > 0 {
		return colToLeft
	}
	if rowsAbove := p.horizontalReflection(smudges); rowsAbove > 0 {
		return rowsAbove * 100
	}
	log.Fatalln("should not happen")
	return 0
}

func (p Pattern) verticalReflection(smudges int) int {
	for colsToLeft := 1; colsToLeft < len(p.lines[0]); colsToLeft++ {
		maxColsToCheck := min(colsToLeft, len(p.lines[0])-colsToLeft)
		smudgeCount := 0
		for i := 0; i < maxColsToCheck; i++ {
			smudgeCount += countSmudgesCol(p.lines, colsToLeft-1-i, colsToLeft+i)
			if smudgeCount > smudges {
				break
			}
		}
		if smudgeCount == smudges {
			return colsToLeft
		}
	}
	return 0
}

func (p Pattern) horizontalReflection(smudges int) int {
	for rowsAbove := 1; rowsAbove < len(p.lines); rowsAbove++ {
		maxLinesToCheck := min(rowsAbove, len(p.lines)-rowsAbove)
		smudgeCount := 0
		for i := 0; i < maxLinesToCheck; i++ {
			smudgeCount += countSmudgesRow(p.lines, rowsAbove-1-i, rowsAbove+i)
			if smudgeCount > smudges {
				break
			}
		}
		if smudgeCount == smudges {
			return rowsAbove
		}
	}
	return 0
}

func countSmudgesCol(lines []string, i, j int) int {
	smudges := 0
	for _, line := range lines {
		if line[i] != line[j] {
			smudges++
		}
	}
	return smudges
}

func countSmudgesRow(lines []string, i, j int) int {
	smudges := 0
	for col := 0; col < len(lines[i]); col++ {
		if lines[i][col] != lines[j][col] {
			smudges++
		}
	}
	return smudges
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
