package main

import (
	"container/list"
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
	if exampleResult1 != 1320 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 145 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	steps := strings.Split(input, ",")
	sum := 0
	for _, step := range steps {
		sum += hash(step)
	}
	return sum
}

func part2(input string) int {
	steps := strings.Split(input, ",")
	boxes := make([]Box, 256)
	for i := range boxes {
		boxes[i] = Box{number: i, lensSlots: list.New()}
	}
	lenses := make(map[string]*list.Element)
	for _, step := range steps {
		label, isRemoveOperation := strings.CutSuffix(step, "-")
		if isRemoveOperation {
			if lensToRemove, ok := lenses[label]; ok {
				boxes[hash(label)].lensSlots.Remove(lensToRemove)
				delete(lenses, label)
			}
		} else {
			label = strings.Split(step, "=")[0]
			focalLength, _ := strconv.Atoi(strings.Split(step, "=")[1])
			if lensElem, ok := lenses[label]; ok {
				lensElem.Value = Lens{label, focalLength}
			} else {
				lensElem := boxes[hash(label)].lensSlots.PushBack(Lens{label, focalLength})
				lenses[label] = lensElem
			}
		}
	}
	totalFocusingPower := 0
	for _, box := range boxes {
		lensElem := box.lensSlots.Front()
		for slot := 1; slot <= box.lensSlots.Len(); slot++ {
			focusingPower := (1 + box.number) * slot * lensElem.Value.(Lens).focalLength
			totalFocusingPower += focusingPower
			lensElem = lensElem.Next()
		}
	}
	return totalFocusingPower
}

type Box struct {
	number    int
	lensSlots *list.List
}

type Lens struct {
	label       string
	focalLength int
}

var hashmap map[string]int = make(map[string]int)

func hash(s string) int {
	if result, ok := hashmap[s]; ok {
		return result
	}

	currentValue := 0
	for _, char := range s {
		currentValue += int(char)
		currentValue *= 17
		currentValue %= 256
	}
	hashmap[s] = currentValue
	return currentValue
}
