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
	if exampleResult1 != 21 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 525152 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	lines := strings.Split(input, "\n")
	sum := 0
	for _, line := range lines {
		record := parseRecord(line)
		sum += record.countArrangements()
	}
	return sum
}

func part2(input string) int {
	lines := strings.Split(input, "\n")
	sum := len(lines)
	return sum
}

func parseRecord(s string) ConditionRecord {
	springs := strings.Split(s, " ")[0]
	groupSizesString := strings.Split(strings.Split(s, " ")[1], ",")
	groupSizes := make([]int, len(groupSizesString))
	for i := range groupSizes {
		groupSizes[i], _ = strconv.Atoi(groupSizesString[i])
	}
	return ConditionRecord{
		springs:    springs,
		groupSizes: groupSizes,
	}
}

type ConditionRecord struct {
	springs    string
	groupSizes []int
}

func (c ConditionRecord) countArrangements() int {
	return c.countArrangementsRec(0, 0)
}

func (c ConditionRecord) countArrangementsRec(springIndex, groupIndex int) int {
	if springIndex >= len(c.springs) {
		if groupIndex != len(c.groupSizes) {
			return 0
		}
		return 1
	}

	count := 0
	switch c.springs[springIndex] {
	case '?':
		count += c.countArrangementsRec(springIndex+1, groupIndex)
		fallthrough
	case '#':
		if groupIndex < len(c.groupSizes) {
			neededGroupSize := c.groupSizes[groupIndex]
			springEndIndex := springIndex + neededGroupSize
			if springEndIndex <= len(c.springs) &&
				(springEndIndex == len(c.springs) || c.springs[springEndIndex] == '?' || c.springs[springEndIndex] == '.') &&
				!strings.Contains(c.springs[springIndex:springEndIndex], ".") {
				return c.countArrangementsRec(springEndIndex+1, groupIndex+1) + count
			}
			return count
		}
		return count
	default:
		return c.countArrangementsRec(springIndex+1, groupIndex) + count
	}
}
