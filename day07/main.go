package main

import (
	_ "embed"
	"log"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed input.txt
var input string

func main() {
	exampleResult1 := part1(example2)
	if exampleResult1 != 6440 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 5905 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	lines := strings.Split(input, "\n")
	hands := make([]Hand, len(lines))
	for i, line := range lines {
		hands[i] = parseHand(line)
	}
	sort.Sort(ByCards(hands))

	totalWinnings := 0
	for i, hand := range hands {
		rank := i + 1
		totalWinnings += (rank * hand.bidAmount)
	}
	return totalWinnings
}

func part2(input string) int {
	lines := strings.Split(input, "\n")
	hands := make([]Hand2, len(lines))
	for i, line := range lines {
		hands[i] = parseHand2(line)
	}
	sort.Sort(ByCards2(hands))

	totalWinnings := 0
	for i, hand := range hands {
		rank := i + 1
		totalWinnings += (rank * hand.bidAmount)
	}
	return totalWinnings
}

type Card byte

func (c Card) Compare(other Card) int {
	const ascLabels string = "23456789TJQKA"
	return strings.Index(ascLabels, string(c)) - strings.Index(ascLabels, string(other))
}

func parseHand(s string) Hand {
	bidAmount, _ := strconv.Atoi(strings.Split(s, " ")[1])
	cards := strings.Split(s, " ")[0]
	return Hand{
		cards:     cards,
		bidAmount: bidAmount,
		handType:  handType(cards),
	}
}

func handType(cards string) int {
	cardCount := make(map[rune]int)
	for _, char := range cards {
		cardCount[char]++
	}
	cardCounts := maps.Values(cardCount)
	sort.Slice(cardCounts, func(i, j int) bool { return cardCounts[i] > cardCounts[j] })
	if cardCounts[0] == 5 {
		return 7
	}
	if cardCounts[0] == 4 {
		return 6
	}
	if cardCounts[0] == 3 && cardCounts[1] == 2 {
		return 5
	}
	if cardCounts[0] == 3 {
		return 4
	}
	if cardCounts[0] == 2 && cardCounts[1] == 2 {
		return 3
	}
	if cardCounts[0] == 2 {
		return 2
	}
	return 1
}

type Hand struct {
	cards     string
	bidAmount int
	handType  int
}

type ByCards []Hand

// Implement the sort.Interface interface
func (a ByCards) Len() int { return len(a) }
func (a ByCards) Less(i, j int) bool {
	if a[i].handType != a[j].handType {
		return a[i].handType < a[j].handType
	}
	for k := 0; k < 5; k++ {
		compare := Card(a[i].cards[k]).Compare(Card(a[j].cards[k]))
		if compare != 0 {
			return compare < 0
		}
	}
	return false
}
func (a ByCards) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// part2

type Card2 byte

func (c Card2) Compare(other Card2) int {
	const ascLabels string = "J23456789TQKA"
	return strings.Index(ascLabels, string(c)) - strings.Index(ascLabels, string(other))
}

func parseHand2(s string) Hand2 {
	bidAmount, _ := strconv.Atoi(strings.Split(s, " ")[1])
	cards := strings.Split(s, " ")[0]
	return Hand2{
		cards:     cards,
		bidAmount: bidAmount,
		handType:  handType2(cards),
	}
}

func handType2(cards string) int {
	cardCount := make(map[rune]int)
	for _, char := range cards {
		cardCount[char]++
	}

	cardCounts := maps.Values(cardCount)
	sort.Slice(cardCounts, func(i, j int) bool { return cardCounts[i] > cardCounts[j] })
	if cardCounts[0] == 5 {
		return 7
	}
	if cardCounts[0] == 4 {
		if cardCount['J'] > 0 {
			return 7
		}
		return 6
	}
	if cardCounts[0] == 3 && cardCounts[1] == 2 {
		if cardCount['J'] > 0 {
			return 7
		}
		return 5
	}
	if cardCounts[0] == 3 {
		if cardCount['J'] > 0 {
			return 6
		}
		return 4
	}
	if cardCounts[0] == 2 && cardCounts[1] == 2 {
		if cardCount['J'] == 2 {
			return 6
		}
		if cardCount['J'] == 1 {
			return 5
		}
		return 3
	}
	if cardCounts[0] == 2 {
		if cardCount['J'] > 0 {
			return 4
		}
		return 2
	}
	if cardCount['J'] == 1 {
		return 2
	}
	return 1
}

type Hand2 struct {
	cards     string
	bidAmount int
	handType  int
}

type ByCards2 []Hand2

// Implement the sort.Interface interface
func (a ByCards2) Len() int { return len(a) }
func (a ByCards2) Less(i, j int) bool {
	if a[i].handType != a[j].handType {
		return a[i].handType < a[j].handType
	}
	for k := 0; k < 5; k++ {
		compare := Card2(a[i].cards[k]).Compare(Card2(a[j].cards[k]))
		if compare != 0 {
			return compare < 0
		}
	}
	return false
}
func (a ByCards2) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
