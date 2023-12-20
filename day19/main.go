package main

import (
	_ "embed"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	start := time.Now()
	workflows := parseWorkflows(strings.Split(input, "\n\n")[0])
	rawParts := strings.Split(strings.Split(input, "\n\n")[1], "\n")

	sum := 0
	for _, rawPart := range rawParts {
		part := parsePart(rawPart)
		if workflows.get("in").accepts(part) {
			sum += part.sumRatings()
		}
	}
	elapsed := time.Since(start)
	println("part1:", elapsed.String())
	return sum
}

func part2(input string) int {
	start := time.Now()
	maxValue := 4000
	workflows := parseWorkflows(strings.Split(input, "\n\n")[0])
	combinationsCount := workflows.get("in").countCombinations("A", workflows, CategoryRatings{
		IntRange{1, maxValue + 1},
		IntRange{1, maxValue + 1},
		IntRange{1, maxValue + 1},
		IntRange{1, maxValue + 1},
	})
	elapsed := time.Since(start)
	println("part2:", elapsed.String())
	return combinationsCount
}

func (this Workflow) accepts(part Part) bool {
	if this.name == "A" {
		return true
	}
	if this.name == "R" {
		return false
	}
	for _, rule := range this.rules {
		if rule.accepts(part) {
			return rule.next.accepts(part)
		}
	}
	return false
}

func (this Rule) accepts(part Part) bool {
	i := strings.Index("xmas", this.category)
	switch this.op {
	case ">":
		return part.ratings[i] > this.value
	case "<":
		return part.ratings[i] < this.value
	default:
		return true
	}
}

type Part struct {
	ratings []int
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

func (p Part) sumRatings() int {
	sum := 0
	for i := range p.ratings {
		sum += p.ratings[i]
	}
	return sum
}

func parseWorkflows(s string) Workflows {
	rawWorkflows := strings.Split(s, "\n")
	workflows := make(map[string]*Workflow)
	for _, rawWorkflow := range rawWorkflows {
		parseWorkflow(rawWorkflow, workflows)
	}
	return workflows
}

func parseWorkflow(s string, workflows Workflows) {
	s = s[:len(s)-1]
	name := strings.Split(s, "{")[0]
	rawRules := strings.Split(strings.Split(s, "{")[1], ",")
	rules := make([]Rule, len(rawRules))
	for i := range rawRules {
		rules[i] = parseRule(rawRules[i], workflows)
	}
	workflows.update(name, Workflow{name, rules})
}

func parseRule(s string, workflows Workflows) Rule {
	if !strings.Contains(s, ":") {
		return Rule{
			Condition: Condition{
				category: "",
				op:       "",
				value:    -1,
			},
			next: workflows.get(s),
		}
	}

	re := regexp.MustCompile("([xmas])([<>])([0-9]+)\\:([a-zA-Z]+)")
	submatches := re.FindStringSubmatch(s)
	category := submatches[1]
	operation := submatches[2]
	value, _ := strconv.Atoi(submatches[3])
	nextWorkflow := submatches[4]

	return Rule{
		Condition: Condition{
			category: category,
			op:       operation,
			value:    value,
		},
		next: workflows.get(nextWorkflow),
	}
}

type Workflows map[string]*Workflow

func (this Workflows) get(name string) *Workflow {
	if _, ok := this[name]; !ok {
		this[name] = &Workflow{
			name: name,
		}
	}
	return this[name]
}

func (this Workflows) update(name string, workflow Workflow) {
	if _, ok := this[name]; !ok {
		this[name] = &Workflow{
			name: name,
		}
	}
	this[name].name = workflow.name
	this[name].rules = workflow.rules
}

type Workflow struct {
	name  string
	rules []Rule
}

func (this Workflow) countCombinations(target string, workflows Workflows, acceptsRatings CategoryRatings) int {
	result := 0
	prevCondition := Condition{}
	for _, rule := range this.rules {
		acceptsRatings = acceptsRatings.restrictReverse(prevCondition)
		currentRuleAcceptRanges := acceptsRatings.restrict(rule.Condition)
		if rule.next.name == target {
			result += currentRuleAcceptRanges.combinations()
		} else {
			result += rule.next.countCombinations(target, workflows, currentRuleAcceptRanges)
		}
		prevCondition = rule.Condition
	}
	return result
}

type CategoryRatings []IntRange

func (this CategoryRatings) String() string {
	builder := strings.Builder{}
	for _, r := range this {
		builder.WriteString(fmt.Sprintf("[%d,%d)", r.start, r.end))
	}
	return builder.String()
}

func (this CategoryRatings) combinations() int {
	result := 1
	for i := range this {
		result *= this[i].Size()
	}
	return result
}

func (this CategoryRatings) restrict(c Condition) CategoryRatings {
	newRatings := make(CategoryRatings, len(this))
	copy(newRatings, this)
	categoryIndex := strings.Index("xmas", c.category)
	switch c.op {
	case ">":
		newRatings[categoryIndex].start = c.value + 1
	case "<":
		newRatings[categoryIndex].end = c.value
	}
	return newRatings
}

func (this CategoryRatings) restrictReverse(c Condition) CategoryRatings {
	newRatings := make(CategoryRatings, len(this))
	copy(newRatings, this)
	categoryIndex := strings.Index("xmas", c.category)
	switch c.op {
	case ">":
		newRatings[categoryIndex].end = c.value + 1
	case "<":
		newRatings[categoryIndex].start = c.value
	}
	return newRatings
}

type Rule struct {
	Condition
	next *Workflow
}

type Condition struct {
	op       string
	value    int
	category string
}

type IntRange struct {
	start, end int
}

func (this IntRange) Size() int { return this.end - this.start }
