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
	exampleResult1 := part1(example2)
	if exampleResult1 != 35 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 46 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	sections := strings.Split(input, "\n\n")
	maps := make([]Map, len(sections)-1)
	for i := range maps {
		maps[i] = parseMap(sections[i+1])
	}
	almanac := Almanac{
		maps: maps,
	}

	seedsWithouPrefix, _ := strings.CutPrefix(sections[0], "seeds: ")
	seedsString := strings.Split(seedsWithouPrefix, " ")
	seeds := make([]int, len(seedsString))
	for i, s := range seedsString {
		seeds[i], _ = strconv.Atoi(s)
	}

	minLocation := almanac.Location(seeds[3])
	for _, seed := range seeds[1:] {
		location := almanac.Location(seed)
		if location < minLocation {
			minLocation = location
		}
	}

	return minLocation
}

func part2(input string) int {
	sections := strings.Split(input, "\n\n")
	maps := make([]Map, len(sections)-1)
	for i := range maps {
		maps[i] = parseMap(sections[i+1])
	}
	almanac := Almanac{
		maps: maps,
	}

	seedsWithouPrefix, _ := strings.CutPrefix(sections[0], "seeds: ")
	seedsString := strings.Split(seedsWithouPrefix, " ")
	seeds := make([]int, len(seedsString))
	for i, s := range seedsString {
		seeds[i], _ = strconv.Atoi(s)
	}

	minLocation := almanac.Location(seeds[0])
	for i, seedStart := range seeds {
		if i%2 == 0 {
			for j := 0; j < seeds[i+1]; j++ {
				location := almanac.Location(seedStart + j)
				if location < minLocation {
					minLocation = location
				}
			}
			println("progress", i/2+1, "/", len(seeds)/2)
		}
	}

	return minLocation
}

type Almanac struct {
	maps []Map
}

func (a Almanac) Location(seed int) int {
	temp := seed
	for _, _map := range a.maps {
		temp = _map.Map(temp)
	}
	return temp
}

func parseMap(s string) Map {
	lines := strings.Split(s, "\n")[1:]
	ranges := make([]Range, len(lines))
	for i, line := range lines {
		sourceRangeStart, _ := strconv.Atoi(strings.Split(line, " ")[1])
		destinationRangeStart, _ := strconv.Atoi(strings.Split(line, " ")[0])
		length, _ := strconv.Atoi(strings.Split(line, " ")[2])
		ranges[i] = Range{
			sourceRangeStart:      sourceRangeStart,
			destinationRangeStart: destinationRangeStart,
			length:                length,
		}
	}
	return Map{
		ranges: ranges,
	}
}

type Map struct {
	ranges []Range
}

func (m Map) Map(source int) int {
	for _, r := range m.ranges {
		if r.CanMap(source) {
			return r.Map(source)
		}
	}
	return source
}

type Range struct {
	sourceRangeStart      int
	destinationRangeStart int
	length                int
}

func (r Range) CanMap(source int) bool {
	return r.sourceRangeStart <= source && source < r.sourceRangeStart+r.length
}

func (r Range) Map(source int) int {
	return source + (r.destinationRangeStart - r.sourceRangeStart)
}
