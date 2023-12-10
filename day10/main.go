package main

import (
	_ "embed"
	"log"
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
	if exampleResult1 != 8 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 10 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	lines := strings.Split(input, "\n")
	tiles := make([]string, len(lines)+2)
	tiles[0] = strings.Repeat(".", len(lines[0])+2)
	tiles[len(tiles)-1] = strings.Repeat(".", len(lines[0])+2)
	for i := range lines {
		tiles[i+1] = "." + lines[i] + "."
	}

	iStart, jStart := findStart(tiles)
	distance := 1
	iPrev1, jPrev1 := iStart, jStart
	iPrev2, jPrev2 := iStart, jStart
	iCurr1, jCurr1 := findFirst(tiles, iStart, jStart, true)
	iCurr2, jCurr2 := findFirst(tiles, iStart, jStart, false)
	for {
		if iCurr1 == iCurr2 && jCurr1 == jCurr2 {
			return distance
		}
		if iCurr1 == iStart && jCurr1 == jStart {
			log.Fatalln("looped")
		}

		iNext1, jNext1 := findNext(tiles, iCurr1, jCurr1, iPrev1, jPrev1)
		iNext2, jNext2 := findNext(tiles, iCurr2, jCurr2, iPrev2, jPrev2)
		iPrev1, jPrev1 = iCurr1, jCurr1
		iPrev2, jPrev2 = iCurr2, jCurr2
		iCurr1, jCurr1 = iNext1, jNext1
		iCurr2, jCurr2 = iNext2, jNext2
		distance++
	}
}

func part2(input string) int {
	lines := strings.Split(input, "\n")
	tiles := make([]string, len(lines)+2)
	tiles[0] = strings.Repeat(".", len(lines[0])+2)
	tiles[len(tiles)-1] = strings.Repeat(".", len(lines[0])+2)
	for i := range lines {
		tiles[i+1] = "." + lines[i] + "."
	}
	isLoop := make([][]bool, len(tiles))
	for i := range isLoop {
		isLoop[i] = make([]bool, len(tiles[0]))
	}

	iStart, jStart := findStart(tiles)
	isLoop[iStart][jStart] = true
	iPrev1, jPrev1 := iStart, jStart
	iPrev2, jPrev2 := iStart, jStart
	iCurr1, jCurr1 := findFirst(tiles, iStart, jStart, true)
	iCurr2, jCurr2 := findFirst(tiles, iStart, jStart, false)
	startPipeType := '.'
	switch {
	case iCurr1 != iStart && iCurr2 != iStart:
		startPipeType = '|'
	case jCurr1 != jStart && jCurr2 != jStart:
		startPipeType = '-'
	case jCurr1 == jStart+1 && iCurr2 == iStart-1:
		startPipeType = 'L'
	case jCurr1 == jStart-1 && iCurr2 == iStart-1:
		startPipeType = 'J'
	case iCurr1 == iStart+1 && jCurr2 == jStart-1:
		startPipeType = '7'
	case jCurr1 == jStart+1 && iCurr2 == iStart+1:
		startPipeType = 'F'
	}
	tiles[iStart] = strings.Replace(tiles[iStart], "S", string(startPipeType), 1)
	for {
		isLoop[iCurr1][jCurr1] = true
		isLoop[iCurr2][jCurr2] = true
		if iCurr1 == iCurr2 && jCurr1 == jCurr2 {
			break
		}
		if iCurr1 == iStart && jCurr1 == jStart {
			log.Fatalln("looped")
		}

		iNext1, jNext1 := findNext(tiles, iCurr1, jCurr1, iPrev1, jPrev1)
		iNext2, jNext2 := findNext(tiles, iCurr2, jCurr2, iPrev2, jPrev2)
		iPrev1, jPrev1 = iCurr1, jCurr1
		iPrev2, jPrev2 = iCurr2, jCurr2
		iCurr1, jCurr1 = iNext1, jNext1
		iCurr2, jCurr2 = iNext2, jNext2
	}

	tilesCount := 0
	isInLoop := false
	prevLoopBorder := '.'
	for i := range tiles {
		isInLoop = false
		for j, tile := range tiles[i] {
			if isLoop[i][j] {
				switch tile {
				case '|':
					fallthrough
				case 'L':
					fallthrough
				case 'F':
					isInLoop = !isInLoop
					prevLoopBorder = tile
				case 'J':
					if prevLoopBorder != 'F' {
						isInLoop = !isInLoop
					}
					prevLoopBorder = tile
				case '7':
					if prevLoopBorder != 'L' {
						isInLoop = !isInLoop
					}
					prevLoopBorder = tile
				}
			} else if isInLoop {
				tilesCount++
			}
		}
	}

	return tilesCount
}

func findStart(tiles []string) (int, int) {
	for i := range tiles {
		for j := range tiles[i] {
			if tiles[i][j] == 'S' {
				return i, j
			}
		}
	}
	log.Fatalln("Start not found")
	return -1, -1
}

func findFirst(tiles []string, i, j int, direction bool) (int, int) {
	if direction {
		if strings.Contains("-J7", string(tiles[i][j+1])) {
			return i, j + 1
		}
		if strings.Contains("|JL", string(tiles[i+1][j])) {
			return i + 1, j
		}
		if strings.Contains("-LF", string(tiles[i][j-1])) {
			return i, j - 1
		}
		if strings.Contains("|7F", string(tiles[i-1][j])) {
			return i - 1, j
		}
	} else {
		if strings.Contains("|7F", string(tiles[i-1][j])) {
			return i - 1, j
		}
		if strings.Contains("-LF", string(tiles[i][j-1])) {
			return i, j - 1
		}
		if strings.Contains("|JL", string(tiles[i+1][j])) {
			return i + 1, j
		}
		if strings.Contains("-J7", string(tiles[i][j+1])) {
			return i, j + 1
		}
	}
	log.Fatalln("error findFirst")
	return -1, -1
}

func findNext(tiles []string, i, j, iPrev, jPrev int) (int, int) {
	switch tiles[i][j] {
	case '|':
		if i-1 == iPrev && j == jPrev {
			return i + 1, j
		} else {
			return i - 1, j
		}
	case '-':
		if i == iPrev && j-1 == jPrev {
			return i, j + 1
		} else {
			return i, j - 1
		}
	case 'L':
		if i-1 == iPrev && j == jPrev {
			return i, j + 1
		} else {
			return i - 1, j
		}
	case 'J':
		if i-1 == iPrev && j == jPrev {
			return i, j - 1
		} else {
			return i - 1, j
		}
	case '7':
		if i == iPrev && j-1 == jPrev {
			return i + 1, j
		} else {
			return i, j - 1
		}
	case 'F':
		if i == iPrev && j+1 == jPrev {
			return i + 1, j
		} else {
			return i, j + 1
		}
	}
	log.Fatalln("error findNext")
	return -1, -1
}
