package main

import (
	_ "embed"
	"log"
	"regexp"
	"strconv"
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
	if part1(example1) != 142 {
		log.Fatal("Part 1 wrong")
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult := part2(example2)
	if exampleResult != 281 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int64 {
	lines := strings.Split(input, "\n")
	sum := int64(0)
	for _, line := range lines {
		firstDigit := string(line[strings.IndexFunc(line, unicode.IsDigit)])
		lastDigit := string(line[strings.LastIndexFunc(line, unicode.IsDigit)])
		twoDigitNumber, _ := strconv.ParseInt(firstDigit+lastDigit, 0, 0)
		sum += twoDigitNumber
	}
	return sum
}

func part2(input string) int64 {
	lines := strings.Split(input, "\n")
	sum := int64(0)
	for _, line := range lines {
		re := regexp.MustCompile("[0-9]|one|two|three|four|five|six|seven|eight|nine|zero")
		reReverse := regexp.MustCompile("[0-9]|orez|enin|thgie|neves|xis|evif|ruof|eerht|owt|eno")
		firstDigit := re.FindString(line)
		lastDigit := Reverse(reReverse.FindString(Reverse(line)))
		twoDigitNumber, _ := strconv.ParseInt(ReplaceLetters(firstDigit)+ReplaceLetters(lastDigit), 0, 0)
		println(twoDigitNumber)
		sum += twoDigitNumber
	}
	return sum
}

func ReplaceLetters(s string) string {
	switch s {
	case "one":
		return "1"
	case "two":
		return "2"
	case "three":
		return "3"
	case "four":
		return "4"
	case "five":
		return "5"
	case "six":
		return "6"
	case "seven":
		return "7"
	case "eight":
		return "8"
	case "nine":
		return "9"
	case "zero":
		return "0"
	default:
		return s
	}
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
