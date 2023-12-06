package main

import (
	_ "embed"
	"log"
	"regexp"
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
	if exampleResult1 != 288 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 71503 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	races := parseRaces(input)
	result := 1
	for _, race := range races {
		result *= race.countWaysToWin()
	}
	return result
}

func part2(input string) int {
	race := parseRace(input)
	return race.countWaysToWin()
}

func parseRaces(input string) []Race {
	re := regexp.MustCompile("\\s+")
	timeLine, _ := strings.CutPrefix(re.ReplaceAllString(strings.Split(input, "\n")[0], " "), "Time: ")
	distanceLine, _ := strings.CutPrefix(re.ReplaceAllString(strings.Split(input, "\n")[1], " "), "Distance: ")
	timeStrings := strings.Split(timeLine, " ")
	distanceStrings := strings.Split(distanceLine, " ")
	races := make([]Race, len(timeStrings))
	for i := range races {
		time, _ := strconv.Atoi(timeStrings[i])
		distance, _ := strconv.Atoi(distanceStrings[i])
		races[i] = Race{
			timeAllowed:    time,
			recordDistance: distance,
		}
	}
	return races
}

func parseRace(input string) Race {
	re := regexp.MustCompile("\\s+")
	timeString, _ := strings.CutPrefix(re.ReplaceAllString(strings.Split(input, "\n")[0], ""), "Time:")
	distanceString, _ := strings.CutPrefix(re.ReplaceAllString(strings.Split(input, "\n")[1], ""), "Distance:")
	time, _ := strconv.Atoi(timeString)
	distance, _ := strconv.Atoi(distanceString)
	return Race{
		timeAllowed:    time,
		recordDistance: distance,
	}
}

type Race struct {
	timeAllowed    int
	recordDistance int
}

func (r Race) countWaysToWin() int {
	winCount := 0
	for i := 0; i < r.timeAllowed; i++ {
		distance := (r.timeAllowed - i) * i
		if distance > r.recordDistance {
			winCount++
		}
	}
	return winCount
}
