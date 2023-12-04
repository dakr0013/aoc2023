package main

import (
	_ "embed"
	"log"
	"math"
	"slices"
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
	if exampleResult1 != 13 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 30 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	lines := strings.Split(input, "\n")
	sum := 0
	for _, line := range lines {
		card := parseCard(line)
		sum += card.worth()
	}
	return sum
}

func part2(input string) int {
	lines := strings.Split(input, "\n")
	cardCounts := make([]int, len(lines))
	for i := range cardCounts {
		cardCounts[i] = 1
	}
	sum := 0
	for i, line := range lines {
		card := parseCard(line)
		winCount := card.winCount()
		copiedCards := cardCounts[i+1 : i+1+winCount]
		for j, _ := range copiedCards {
			copiedCards[j] += cardCounts[i]
		}
		sum += cardCounts[i]
	}
	return sum
}

func parseCard(s string) Card {
	cardId, _ := strconv.Atoi(strings.Trim(strings.Split(s[5:], ":")[0], " "))
	allNumbers := strings.Trim(strings.Split(s[5:], ":")[1], " ")
	winningNumbers := strings.Trim(strings.Split(allNumbers, "|")[0], " ")
	numbers := strings.Trim(strings.Split(allNumbers, "|")[1], " ")

	return Card{
		id:             cardId,
		winningNumbers: parseNumbers(winningNumbers),
		numbers:        parseNumbers(numbers),
	}
}

func parseNumbers(s string) []int {
	stringNumbers := strings.Split(strings.ReplaceAll(s, "  ", " "), " ")
	numbers := make([]int, len(stringNumbers))
	for i, num := range stringNumbers {
		numbers[i], _ = strconv.Atoi(num)
	}
	return numbers
}

type Card struct {
	id             int
	winningNumbers []int
	numbers        []int
}

func (c Card) winCount() int {
	winCount := 0
	for _, num := range c.numbers {
		if slices.Contains(c.winningNumbers, num) {
			winCount++
		}
	}
	return winCount
}

func (c Card) worth() int {
	winCount := c.winCount()
	if winCount > 0 {
		return int(math.Pow(float64(2), float64(winCount-1)))
	} else {
		return 0
	}
}
