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
	println("part 2 ok")
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
	sum := 0
	for _, line := range lines {
		record := parseRecord(line)
		sum += record.unfold(5).countArrangements()
	}
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

func (c ConditionRecord) unfold(times int) ConditionRecord {
	newGroupSizes := make([]int, len(c.groupSizes)*times)
	for i := 0; i < times; i++ {
		copy(newGroupSizes[i*len(c.groupSizes):], c.groupSizes)
	}
	return ConditionRecord{
		springs:    strings.Repeat(c.springs+"?", times)[:len(c.springs)*times+times-1],
		groupSizes: newGroupSizes,
	}
}

type State struct {
	springIndex, groupIndex int
}

var cache map[State]int

func (c ConditionRecord) countArrangements() int {
	cache = make(map[State]int)
	return c.countArrangementsRec(State{0, 0})
}

func (c ConditionRecord) countArrangementsRec(state State) int {
	if result, ok := cache[state]; ok {
		return result
	}

	springIndex, groupIndex := state.springIndex, state.groupIndex
	if springIndex >= len(c.springs) {
		if groupIndex != len(c.groupSizes) {
			cache[state] = 0
			return 0
		}
		cache[state] = 1
		return 1
	}

	count := 0
	switch c.springs[springIndex] {
	case '?':
		count += c.countArrangementsRec(State{springIndex + 1, groupIndex})
		fallthrough
	case '#':
		if groupIndex < len(c.groupSizes) {
			neededGroupSize := c.groupSizes[groupIndex]
			springEndIndex := springIndex + neededGroupSize
			if springEndIndex <= len(c.springs) &&
				(springEndIndex == len(c.springs) || c.springs[springEndIndex] == '?' || c.springs[springEndIndex] == '.') &&
				!strings.Contains(c.springs[springIndex:springEndIndex], ".") {
				result := c.countArrangementsRec(State{springEndIndex + 1, groupIndex + 1}) + count
				cache[state] = result
				return result
			}
			cache[state] = count
			return count
		}
		cache[state] = count
		return count
	default:
		result := c.countArrangementsRec(State{springIndex + 1, groupIndex}) + count
		cache[state] = result
		return result
	}
}
