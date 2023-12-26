package main

import (
	_ "embed"
	"log"
	"regexp"
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
	exampleResult1 := part1(example1)
	if exampleResult1 != 32000000 {
		log.Fatalf("Part 1 wrong; acutal: %d\n", exampleResult1)
	}
	log.Printf("Part 1: %d\n", part1(input))
	log.Printf("Part 2: %d\n", part2(input))
}

func part1(input string) int {
	modules := parseModules(input)
	for i := 0; i < 1000; i++ {
		pushButton(modules)
	}

	lowPulseCount := 1000
	highPulseCount := 0
	for _, module := range modules {
		lowPulseCount += module.LowPulseCount()
		highPulseCount += module.HighPulseCount()
	}
	return lowPulseCount * highPulseCount
}

// module which ouputs rx is &
// 1. check each of its inputs individually how many button presses it takes for a high pulse
// 2. take least common multiple of those counts (high pulses occur at fixed intervals)
func part2(input string) int {
	modules := parseModules(input)
	re := regexp.MustCompile("(..) -> rx")
	targetModuleName := re.FindStringSubmatch(input)[1]
	targetModule := modules[targetModuleName].(*ConjunctionModule)

	targetModuleInputs := maps.Keys(targetModule.currentInputs)

	buttonPressCount := 0
	counts := map[string]int{}
	for {
		buttonPressCount++
		pushButton(modules)
		for _, name := range targetModuleInputs {
			if _, ok := counts[name]; !ok && modules[name].HighPulseCount() == 1 {
				counts[name] = buttonPressCount
			}
		}
		if len(counts) == len(targetModuleInputs) {
			break
		}
	}

	return lcm(maps.Values(counts)...)
}

func pushButton(modules map[string]Module) {
	toProcess := []string{"broadcaster"}
	for len(toProcess) != 0 {
		nextModName := toProcess[0]
		toProcess = toProcess[1:]

		if nextModule, ok := modules[nextModName]; ok {
			destinations := nextModule.ProcessPulse(modules)
			toProcess = append(toProcess, destinations...)
		}
	}
}

func parseModules(s string) map[string]Module {
	rawModules := strings.Split(s, "\n")
	modules := make(map[string]Module)
	moduleInputs := make(map[string][]string)

	for _, rawModule := range rawModules {
		module := parseModule(rawModule, moduleInputs)
		modules[module.Name()] = module
	}

	for name, inputs := range moduleInputs {
		if module, ok := modules[name]; ok {
			module.WireInputs(inputs)
		}
	}

	return modules
}

func parseModule(s string, moduleInputs map[string][]string) Module {
	moduleName := strings.Split(strings.ReplaceAll(s, " ", ""), "->")[0]
	destinations := strings.Split(strings.Split(strings.ReplaceAll(s, " ", ""), "->")[1], ",")
	moduleType := moduleName[0]

	var result Module
	switch moduleType {
	case '%':
		moduleName = moduleName[1:]
		result = &FlipFlopModule{
			name:          moduleName,
			currentState:  false,
			currentInputs: []bool{},
			destinations:  destinations,
		}
	case '&':
		moduleName = moduleName[1:]
		result = &ConjunctionModule{
			name:          moduleName,
			currentInputs: make(map[string]bool),
			destinations:  destinations,
		}
	default:
		result = &BroadcasterModule{
			name:         moduleName,
			currentInput: false,
			destinations: destinations,
		}
	}

	for _, dest := range destinations {
		if _, ok := moduleInputs[dest]; !ok {
			moduleInputs[dest] = make([]string, 0)
		}
		moduleInputs[dest] = append(moduleInputs[dest], moduleName)
	}

	return result
}

type Module interface {
	Name() string
	SendPulse(from string, pulse bool)
	ProcessPulse(modules map[string]Module) []string
	WireInputs(inputs []string)
	LowPulseCount() int
	HighPulseCount() int
}

type BroadcasterModule struct {
	name           string
	currentInput   bool
	destinations   []string
	lowPulseCount  int
	highPulseCount int
}

func (b *BroadcasterModule) Name() string                      { return b.name }
func (b *BroadcasterModule) WireInputs(inputs []string)        {}
func (b *BroadcasterModule) SendPulse(from string, pulse bool) { b.currentInput = pulse }
func (b *BroadcasterModule) ProcessPulse(modules map[string]Module) []string {
	outputPulse := b.currentInput
	for _, dest := range b.destinations {
		if outputPulse {
			b.highPulseCount++
		} else {
			b.lowPulseCount++
		}
		if module, ok := modules[dest]; ok {
			module.SendPulse(b.name, outputPulse)
		}
	}
	return b.destinations
}
func (b *BroadcasterModule) LowPulseCount() int  { return b.lowPulseCount }
func (b *BroadcasterModule) HighPulseCount() int { return b.highPulseCount }

type FlipFlopModule struct {
	name           string
	currentState   bool
	currentInputs  []bool
	destinations   []string
	lowPulseCount  int
	highPulseCount int
}

func (b *FlipFlopModule) Name() string               { return b.name }
func (b *FlipFlopModule) WireInputs(inputs []string) {}
func (b *FlipFlopModule) SendPulse(from string, pulse bool) {
	b.currentInputs = append(b.currentInputs, pulse)
}
func (b *FlipFlopModule) ProcessPulse(modules map[string]Module) []string {
	currentInput := b.currentInputs[0]
	b.currentInputs = b.currentInputs[1:]
	if !currentInput {
		b.currentState = !b.currentState
		outputPulse := b.currentState
		for _, dest := range b.destinations {
			if outputPulse {
				b.highPulseCount++
			} else {
				b.lowPulseCount++
			}
			if module, ok := modules[dest]; ok {
				module.SendPulse(b.name, outputPulse)
			}
		}
		return b.destinations
	}
	return []string{}
}
func (b *FlipFlopModule) LowPulseCount() int  { return b.lowPulseCount }
func (b *FlipFlopModule) HighPulseCount() int { return b.highPulseCount }

type ConjunctionModule struct {
	name           string
	currentInputs  map[string]bool
	destinations   []string
	lowPulseCount  int
	highPulseCount int
}

func (b *ConjunctionModule) Name() string { return b.name }
func (b *ConjunctionModule) WireInputs(inputs []string) {
	for _, inpModule := range inputs {
		b.currentInputs[inpModule] = false
	}
}
func (b *ConjunctionModule) SendPulse(from string, pulse bool) { b.currentInputs[from] = pulse }
func (b *ConjunctionModule) ProcessPulse(modules map[string]Module) []string {
	outputPulse := true
	if allHigh(b.currentInputs) {
		outputPulse = false
	}
	for _, dest := range b.destinations {
		if outputPulse {
			b.highPulseCount++
		} else {
			b.lowPulseCount++
		}
		if module, ok := modules[dest]; ok {
			module.SendPulse(b.name, outputPulse)
		}
	}
	return b.destinations
}
func (b *ConjunctionModule) LowPulseCount() int  { return b.lowPulseCount }
func (b *ConjunctionModule) HighPulseCount() int { return b.highPulseCount }

func allHigh(inputs map[string]bool) bool {
	result := true
	for _, input := range inputs {
		if !input {
			return false
		}
	}
	return result
}

// greatest common divisor
func gcd(nums ...int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	a := nums[0]
	if len(nums) == 2 {
		b := nums[1]
		if a == 0 {
			return b
		}
		for b != 0 {
			if a > b {
				a = a - b
			} else {
				b = b - a
			}
		}
		return a
	}

	overallGcd := gcd(a, nums[1])
	for _, num := range nums[2:] {
		overallGcd = gcd(overallGcd, num)
	}
	return overallGcd
}

// least common multiple
func lcm(nums ...int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	a := nums[0]
	if len(nums) == 2 {
		b := nums[1]
		return abs((a / gcd(a, b)) * b)
	}

	overallLcm := lcm(a, nums[1])
	for _, num := range nums[2:] {
		overallLcm = lcm(overallLcm, num)
	}
	return overallLcm
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
