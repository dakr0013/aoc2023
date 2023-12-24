package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"math"
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
	exampleResult1 := part1(example1, 7, 27)
	if exampleResult1 != 2 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input, 200000000000000, 400000000000000))

	exampleResult2 := part2(example2)
	if exampleResult2 != 123 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string, from, to float64) int {
	lines := parseLines(input, true)
	return countIntersections(lines, from, to)
}

func part2(input string) int {
	lines := strings.Split(input, "\n")
	sum := len(lines)
	return sum
}

func countIntersections(lines []Line, from, to float64) int {
	count := 0
	for i := 0; i < len(lines)-1; i++ {
		for j := i + 1; j < len(lines); j++ {
			f := lines[i]
			g := lines[j]
			println("--")
			println(f.String())
			println(g.String())
			if f.isParallelTo(g) {
				println("lines are parallel")
				if f.isIdentical(g) {
					println("lines are identical")
					t := (f.p.x - g.p.x) / (g.v.x - f.v.x)
					if t > 0 {
						intersection := f.call(t)
						if intersection.isWithin(from, to) {
							count++
							fmt.Printf("Paths will cross inside test area %s\n", intersection)
						} else {
							fmt.Printf("Paths will cross ouside test area %s\n", intersection)
						}
					} else if !f.call(t).equals(g.call(t)) {
						fmt.Printf("Paths will never cross\n")
					} else if t < 0 {
						fmt.Printf("Paths crossed in the past\n")
					}
				}
				continue
			}
			isSkew, t1, t2 := f.isSkew(g)
			if isSkew {
				println("lines are skew")
				continue
			}
			if t1 >= 0 && t2 >= 0 {
				intersection := f.call(t1)
				if intersection.isWithin(from, to) {
					count++
					fmt.Printf("Paths will cross inside test area %s\n", intersection)
				} else {
					fmt.Printf("Paths will cross ouside test area %s\n", intersection)
				}
			} else {
				fmt.Printf("Paths crossed in the past\n")
			}
		}
	}
	return count
}

func parseLines(s string, ignoreZ bool) []Line {
	rawLines := strings.Split(s, "\n")
	lines := make([]Line, len(rawLines))
	for i := range rawLines {
		lines[i] = parseLine(rawLines[i], ignoreZ)
	}
	return lines
}

func parseLine(s string, ignoreZ bool) Line {
	coords := strings.Split(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "@", ","), ",")
	px, _ := strconv.ParseFloat(coords[0], 64)
	py, _ := strconv.ParseFloat(coords[1], 64)
	pz, _ := strconv.ParseFloat(coords[2], 64)
	vx, _ := strconv.ParseFloat(coords[3], 64)
	vy, _ := strconv.ParseFloat(coords[4], 64)
	vz, _ := strconv.ParseFloat(coords[5], 64)

	if ignoreZ {
		pz = 0
		vz = 0
	}

	return Line{
		Vector3D{px, py, pz},
		Vector3D{vx, vy, vz},
	}
}

type Line struct {
	p Vector3D
	v Vector3D
}

func (f Line) String() string {
	return fmt.Sprintf("%s + λ * %s", f.p, f.v)
}

func (f Line) call(lambda float64) Vector3D {
	return f.p.add(f.v.scalarMul(lambda))
}

func (f Line) isParallelTo(g Line) bool {
	return f.v.normalized().equals(g.v.normalized()) ||
		f.v.normalized().equals(g.v.normalized().scalarMul(-1))
}

func (f Line) isIdentical(g Line) bool {
	lambda := (f.p.x - g.p.x) / g.v.x
	return math.Abs(f.p.y-(g.p.y+lambda*g.v.y)) < 1e-9 && math.Abs(f.p.z-(g.p.z+lambda*g.v.z)) < 1e-9
}

func (f Line) isSkew(g Line) (isSkew bool, lambda float64, my float64) {
	A := [][]float64{
		{f.v.x, -g.v.x, g.p.x - f.p.x},
		{f.v.y, -g.v.y, g.p.y - f.p.y},
	}
	result, err := solveLinearSystem(A)
	if err != nil {
		return false, 0, 0
	}
	lambda = result[0]
	my = result[1]
	isSkew = math.Abs((f.p.z+lambda*f.v.z)-(g.p.z+my*g.v.z)) > 1e-9
	return
}

type Vector3D struct {
	x, y, z float64
}

func (v Vector3D) String() string {
	return fmt.Sprintf("(%.f,%.f,%.f)", v.x, v.y, v.z)
}

func (v1 Vector3D) equals(v2 Vector3D) bool {
	return math.Abs(v1.x-v2.x) < 1e-9 &&
		math.Abs(v1.y-v2.y) < 1e-9 &&
		math.Abs(v1.z-v2.z) < 1e-9
}

func (v Vector3D) isWithin(from, to float64) bool {
	return from <= v.x && v.x <= to &&
		from <= v.y && v.y <= to
}

func (v1 Vector3D) add(v2 Vector3D) Vector3D {
	return Vector3D{
		v1.x + v2.x,
		v1.y + v2.y,
		v1.z + v2.z,
	}
}

func (v Vector3D) scalarMul(a float64) Vector3D {
	return Vector3D{
		v.x * a,
		v.y * a,
		v.z * a,
	}
}

func (v Vector3D) normalized() Vector3D {
	magnitude := math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
	return Vector3D{
		v.x / magnitude,
		v.y / magnitude,
		v.z / magnitude,
	}
}

func solveLinearSystem(A [][]float64) ([]float64, error) {
	n := len(A)

	// Step 2: Apply Gauss Elimination on Matrix A
	for i := 0; i < n-1; i++ {
		if A[i][i] == 0 {
			return nil, errors.New("Cannot solve")
		}

		for j := i + 1; j < n; j++ {
			ratio := A[j][i] / A[i][i]
			for k := 0; k < n+1; k++ {
				A[j][k] = A[j][k] - ratio*A[i][k]
			}
		}
	}

	// Step 3: Obtaining Solution by Back Substitution
	var X = make([]float64, n)
	X[n-1] = A[n-1][n] / A[n-1][n-1]
	for i := n - 2; i >= 0; i-- {
		X[i] = A[i][n]
		for j := i + 1; j < n; j++ {
			X[i] = X[i] - A[i][j]*X[j]
		}
		X[i] = X[i] / A[i][i]
	}

	return X, nil
}
