package main

import (
	_ "embed"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed input.txt
var input string

func main() {
	exampleResult1 := part1(example1)
	if exampleResult1 != 5 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 7 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	startTime := time.Now()
	snapshot := parseSnapshot(input)
	snapshot.settleBricks()
	result := snapshot.disintegrateCount()
	println("part1:", time.Since(startTime).String())
	return result
}

func part2(input string) int {
	startTime := time.Now()
	snapshot := parseSnapshot(input)
	snapshot.settleBricks()
	result := snapshot.fallingBricksCount()
	println("part2:", time.Since(startTime).String())
	return result
}

func parseSnapshot(s string) Snapshot {
	brickId = 0
	rawBricks := strings.Split(s, "\n")
	bricks := Bricks{}

	for _, rawBrick := range rawBricks {
		brick := parseBrick(rawBrick)
		bricks = append(bricks, &brick)
	}

	return Snapshot{bricks}
}

var brickId int = 0

func parseBrick(s string) Brick {
	brickId++
	coords := strings.Split(strings.Replace(s, "~", ",", 1), ",")
	return Brick{
		brickId,
		parseVector3D(coords[0:3]),
		parseVector3D(coords[3:6]),
	}
}

type Snapshot struct {
	bricks Bricks
}

func (this *Snapshot) settleBricks() {
	this.bricks.sort()
	for i, brick := range this.bricks {
		// bricks are sorted by z coordinate
		// check last 100 bricks because in worst case on the xy-plane could be 100 1x1x1 bricks and the current brick could be 100x100x1
		bricksBelow := this.bricks[max(0, i-100):i]
		brick.settleOn(bricksBelow)
	}
}

func (this *Brick) settleOn(bricksBelow Bricks) {
	for {
		this.move(down)
		if this.intersectsAny(bricksBelow) || this.start[z] < 1 {
			this.move(up)
			break
		}
	}
}

func (this Brick) intersectsAny(bricks []*Brick) bool {
	for _, brick := range bricks {
		if this.isIntersecting(brick) {
			return true
		}
	}
	return false
}

func (this *Brick) move(direction Vector3D) {
	this.start = this.start.move(direction)
	this.end = this.end.move(direction)
}

func (this Snapshot) disintegrateCount() int {
	criticalBrickIds, _ := this.criticalBrickIds()
	return len(this.bricks) - len(criticalBrickIds)
}

func (this Snapshot) fallingBricksCount() int {
	criticalBrickIds, supportedBy := this.criticalBrickIds()
	this.bricks.sort()

	sum := 0
	fallingBrickIds := map[int]bool{}
	for _, id := range criticalBrickIds {
		fallingBrickIds[id] = true
		for _, brick := range this.bricks {
			if containsAll(fallingBrickIds, supportedBy[brick.id]) {
				fallingBrickIds[brick.id] = true
			}
		}
		sum += len(fallingBrickIds) - 1
		fallingBrickIds = map[int]bool{}
	}

	return sum
}

func containsAll(a, b map[int]bool) bool {
	if len(b) == 0 {
		return false
	}
	for key := range b {
		if _, ok := a[key]; !ok {
			return false
		}
	}
	return true
}

// brick is critical if it supports another brick and is the only one supporting that other brick
func (this Snapshot) criticalBrickIds() ([]int, map[int]map[int]bool) {
	blocks := make([][][]int, 10)
	for x := range blocks {
		blocks[x] = make([][]int, 10)
		for y := range blocks[x] {
			blocks[x][y] = make([]int, 400)
		}
	}

	for _, brick := range this.bricks {
		for _, block := range brick.blocks() {
			blocks[block[x]][block[y]][block[z]] = brick.id
		}
	}

	isSupportedBy := map[int]map[int]bool{}
	for i := 1; i <= 1600; i++ {
		isSupportedBy[i] = map[int]bool{}
	}
	for _, brick := range this.bricks {
		for _, block := range brick.blocksOnXYPlane() {
			blockIdBelow := blocks[block[x]][block[y]][block[z]-1]
			if blockIdBelow != 0 {
				isSupportedBy[brick.id][blockIdBelow] = true
			}
		}
	}

	criticalBrickIds := map[int]bool{}
	for _, supportedBy := range isSupportedBy {
		supportedByIds := maps.Keys(supportedBy)
		if len(supportedByIds) == 1 {
			criticalBrickIds[supportedByIds[0]] = true
		}
	}

	return maps.Keys(criticalBrickIds), isSupportedBy
}

type Bricks []*Brick

func (this Bricks) sort() {
	slices.SortFunc(this, func(a, b *Brick) int { return a.start[z] - b.start[z] })
}

type Brick struct {
	id         int
	start, end Vector3D
}

func (this Brick) blocks() []Vector3D {
	direction := this.end.subtract(this.start)
	for xi := range direction {
		if direction[xi] != 0 {
			direction[xi] /= direction[xi]
		}
	}

	result := []Vector3D{}
	current := this.start
	for !current.equals(this.end) {
		result = append(result, current)
		current = current.move(direction)
	}
	return append(result, current)
}

func (this Brick) blocksOnXYPlane() []Vector3D {
	direction := this.end.subtract(this.start)
	for xi := range direction {
		if direction[xi] != 0 {
			direction[xi] /= direction[xi]
		}
	}
	direction[z] = 0
	if direction.equals(Vector3D{0, 0, 0}) {
		return []Vector3D{this.start}
	}

	result := []Vector3D{}
	current := this.start
	for !current.equals(this.end) {
		result = append(result, current)
		current = current.move(direction)
	}
	return append(result, current)
}

func (this Brick) isIntersecting(other *Brick) bool {
	for _, block1 := range this.blocks() {
		for _, block2 := range other.blocks() {
			if block1.equals(block2) {
				return true
			}
		}
	}
	return false
}

type Vector3D []int

var x int = 0
var y int = 1
var z int = 2

var down Vector3D = Vector3D{0, 0, -1}
var up Vector3D = Vector3D{0, 0, 1}

func (this Vector3D) equals(other Vector3D) bool {
	result := true
	for i := range this {
		if this[i] != other[i] {
			return false
		}
	}
	return result
}

func (this Vector3D) move(direction Vector3D) Vector3D {
	newVector := make(Vector3D, len(this))
	for i := range this {
		newVector[i] = this[i] + direction[i]
	}
	return newVector
}

func (this Vector3D) subtract(other Vector3D) Vector3D {
	newVector := make(Vector3D, len(this))
	for i := range this {
		newVector[i] = this[i] - other[i]
	}
	return newVector
}

func parseVector3D(s []string) Vector3D {
	result := make([]int, len(s))
	var err error
	for i := range s {
		result[i], err = strconv.Atoi(s[i])
	}
	if err != nil {
		log.Fatalln("failed parsing string")
	}
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
