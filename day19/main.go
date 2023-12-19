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
	exampleResult1 := part1(example1)
	if exampleResult1 != 19114 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))

	exampleResult2 := part2(example2)
	if exampleResult2 != 167409079868000 {
		log.Fatalf("Part 2 wrong; acutal: %d\n", exampleResult2)
	}
	log.Printf("Part 2: %d\n", part2(input))

}

func part1(input string) int {
	start, workflows, parts := parseInput(input)
	overallSumRatings := 0
	for _, part := range parts {
		if start.isAccepted(workflows, part) {
			overallSumRatings += part.sumRatings()
		}
	}
	return overallSumRatings
}

func part2(input string) int {
	lines := strings.Split(input, "\n")
	sum := len(lines)
	return sum
}

func parseInput(s string) (start Workflow, workflows map[string]Workflow, parts []Part) {
	rawWorkflows := strings.Split(strings.Split(s, "\n\n")[0], "\n")
	rawParts := strings.Split(strings.Split(s, "\n\n")[1], "\n")
	parts = make([]Part, len(rawParts))
	for i := range rawParts {
		parts[i] = parsePart(rawParts[i])
	}

	workflows = make(map[string]Workflow)
	for _, rawWorkflow := range rawWorkflows {
		workflowName, workflow := parseWorkflow(rawWorkflow)
		workflows[workflowName] = workflow
		if workflowName == "in" {
			start = workflow
		}
	}

	return
}

func parseWorkflow(s string) (string, Workflow) {
	s = s[:len(s)-1]
	name := strings.Split(s, "{")[0]
	rawRules := strings.Split(strings.Split(s, "{")[1], ",")
	rules := make([]Rule, len(rawRules))
	for i := range rawRules {
		rules[i] = parseRule(rawRules[i])
	}
	return name, Workflow{rules}
}

func parseRule(s string) Rule {
	if !strings.Contains(s, ":") {
		return func(part Part) (string, bool) { return s, true }
	}
	re := regexp.MustCompile("([xmas])([<>])([0-9]+)\\:([a-zA-Z]+)")
	submatches := re.FindStringSubmatch(s)
	ratingIndex := strings.Index("xmas", submatches[1])
	operation := submatches[2]
	value, _ := strconv.Atoi(submatches[3])
	nextWorkflow := submatches[4]

	if operation == ">" {
		return func(part Part) (string, bool) {
			if part.ratings[ratingIndex] > value {
				return nextWorkflow, true
			}
			return "", false
		}
	}
	return func(part Part) (string, bool) {
		if part.ratings[ratingIndex] < value {
			return nextWorkflow, true
		}
		return "", false
	}
}

func parsePart(s string) Part {
	rawRatings := strings.Split(strings.Trim(s, "{}"), ",")
	ratings := make([]int, len(rawRatings))
	for i := range rawRatings {
		rating, _ := strconv.Atoi(strings.Split(rawRatings[i], "=")[1])
		ratings[i] = rating
	}
	return Part{ratings}
}

type Workflow struct {
	rules []Rule
}

func (start Workflow) isAccepted(workflows map[string]Workflow, p Part) bool {
	currentWorkflow := start
	for {
		next := currentWorkflow.nextWorkflow(p)
		if next == "A" || next == "R" {
			return next == "A"
		}
		currentWorkflow = workflows[next]
	}
}

func (w Workflow) nextWorkflow(p Part) string {
	for _, rule := range w.rules {
		if nextWorkflow, ok := rule(p); ok {
			return nextWorkflow
		}
	}
	log.Fatalln("should never happen")
	return ""
}

type Rule func(part Part) (string, bool)

type Part struct {
	ratings []int
}

func (p Part) sumRatings() int {
	sum := 0
	for i := range p.ratings {
		sum += p.ratings[i]
	}
	return sum
}
