package main

import (
	_ "embed"
	"fmt"
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
	if exampleResult1 != 8 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 2286 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	configuration := CubeSet{
		RedCount:   12,
		GreenCount: 13,
		BlueCount:  14,
	}
	lines := strings.Split(input, "\n")
	sum := 0
	for _, line := range lines {
		game := parseGame(line)
		if game.IsPossibleWith(configuration) {
			sum += game.id
		}
	}
	return sum
}

func part2(input string) int {
	lines := strings.Split(input, "\n")
	sum := 0
	for _, line := range lines {
		game := parseGame(line)
		sum += game.powerOfTheMinimumSetOfCubes()
	}
	return sum
}

func parseGame(line string) Game {
	re := regexp.MustCompile("Game (\\d+): (.*)")
	submatches := re.FindStringSubmatch(line)
	gameId, _ := strconv.Atoi(submatches[1])
	rawCubeSets := strings.Split(submatches[2], ";")
	return Game{
		id:       gameId,
		cubeSets: mapSlice(rawCubeSets, parseCubeSet),
	}
}

func parseCubeSet(raw string) CubeSet {
	rawCounts := strings.Split(strings.Trim(raw, " "), ",")
	cubeSet := CubeSet{}
	for _, rawCount := range rawCounts {
		pair := strings.Split(strings.Trim(rawCount, " "), " ")
		count, _ := strconv.Atoi(pair[0])
		switch pair[1] {
		case "red":
			cubeSet.RedCount = count
		case "green":
			cubeSet.GreenCount = count
		case "blue":
			cubeSet.BlueCount = count
		}
	}
	return cubeSet
}

func mapSlice[T any, M any](a []T, f func(T) M) []M {
	n := make([]M, len(a))
	for i, e := range a {
		n[i] = f(e)
	}
	return n
}

type Game struct {
	id       int
	cubeSets []CubeSet
}

type CubeSet struct {
	RedCount   int
	GreenCount int
	BlueCount  int
}

func (g Game) String() string {
	return fmt.Sprintf("%d", g)
}

func (g Game) IsPossibleWith(configuration CubeSet) bool {
	for _, cubeSet := range g.cubeSets {
		if cubeSet.RedCount > configuration.RedCount ||
			cubeSet.GreenCount > configuration.GreenCount ||
			cubeSet.BlueCount > configuration.BlueCount {
			return false
		}
	}
	return true
}

func (g Game) powerOfTheMinimumSetOfCubes() int {
	minSetOfCubes := CubeSet{}
	for _, cubeSet := range g.cubeSets {
		if cubeSet.RedCount > minSetOfCubes.RedCount {
			minSetOfCubes.RedCount = cubeSet.RedCount
		}
		if cubeSet.GreenCount > minSetOfCubes.GreenCount {
			minSetOfCubes.GreenCount = cubeSet.GreenCount
		}
		if cubeSet.BlueCount > minSetOfCubes.BlueCount {
			minSetOfCubes.BlueCount = cubeSet.BlueCount
		}
	}
	return minSetOfCubes.RedCount * minSetOfCubes.GreenCount * minSetOfCubes.BlueCount
}
