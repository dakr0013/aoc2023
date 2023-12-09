package main

import (
	_ "embed"
	"log"
	"strconv"
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
	if exampleResult1 != 114 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 2 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	lines := strings.Split(input, "\n")
	sum := 0
	for _, line := range lines {
		sequence := parseSequence(line)
		sum += extrapolate(sequence)
	}
	return sum
}

func part2(input string) int {
	lines := strings.Split(input, "\n")
	sum := 0
	for _, line := range lines {
		sequence := parseSequence(line)
		sum += extrapolateBack(sequence)
	}
	return sum
}

func parseSequence(s string) []int {
	seqString := strings.Split(s, " ")
	sequence := make([]int, len(seqString))
	for i := range seqString {
		sequence[i], _ = strconv.Atoi(seqString[i])
	}
	return sequence
}

func extrapolate(seq []int) int {
	diffs := make([]int, len(seq)-1)
	allZero := true
	for i := range diffs {
		diff := seq[i+1] - seq[i]
		diffs[i] = diff
		if diff != 0 {
			allZero = false
		}
	}
	stepSize := 0
	if !allZero {
		stepSize = extrapolate(diffs)
	}
	return seq[len(seq)-1] + stepSize
}

func extrapolateBack(seq []int) int {
	diffs := make([]int, len(seq)-1)
	allZero := true
	for i := range diffs {
		diff := seq[i+1] - seq[i]
		diffs[i] = diff
		if diff != 0 {
			allZero = false
		}
	}
	stepSize := 0
	if !allZero {
		stepSize = extrapolateBack(diffs)
	}
	return seq[0] - stepSize
}
