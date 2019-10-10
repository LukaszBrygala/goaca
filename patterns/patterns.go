package patterns

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type (
	// Pattern mimics to some extent JS micromatch library and can be used to
	// represent VEG fingerprint patterns.
	Pattern struct {
		steps []patternStep
	}

	patternStep interface {
		// match trims input prefix matching current patternStep. Returns boolean
		// flag indicating if beggining of the input matches current patternStep and
		// trimmed (remaining) input string.
		match(input string) (ok bool, remainingInput string)
	}

	exactString string

	star struct {
		until byte
	}

	numberRange struct {
		min int
		max int
	}

	options []exactString
)

// Match checks if given input string matches Pattern.
func (p Pattern) Match(input string) bool {
	for _, step := range p.steps {
		ok, remainingInput := step.match(input)
		if !ok {
			return false
		}
		input = remainingInput
	}
	return len(input) == 0
}

func (step exactString) match(input string) (bool, string) {
	if strings.HasPrefix(input, string(step)) {
		return true, input[len(step):]
	}
	return false, input
}

func (step star) match(input string) (bool, string) {
	if step.until == 0 {
		return true, ""
	}
	idx := strings.IndexByte(input, step.until)
	if idx > -1 {
		return true, input[idx+1:]
	}
	return false, input
}

func (step numberRange) match(input string) (bool, string) {
	endIdx := indexNoDigit(input)
	numStr := input
	if endIdx > -1 {
		numStr = numStr[:endIdx]
	}
	if len(numStr) == 0 {
		return false, input
	}
	num, _ := strconv.Atoi(numStr)
	if num < step.min || num > step.max {
		return false, input
	}
	if endIdx == -1 {
		return true, ""
	}
	return true, input[endIdx:]
}

func (step options) match(input string) (bool, string) {
	for _, opt := range step {
		ok, remaining := opt.match(input)
		if ok {
			return true, remaining
		}
	}
	return false, input
}

// MustCompile compiles give string to Pattern instance or panics in case of
// malformed input.
func MustCompile(s string) Pattern {
	pattern, err := Compile(s)
	if err != nil {
		panic(err)
	}
	return pattern
}

// Compile compiles given string to Pattern instance or returns error in case of
// malformed input.
func Compile(s string) (Pattern, error) {
	var steps []patternStep
	for i := 0; i < len(s); {
		var (
			step patternStep
			err  error
		)
		switch s[i] {
		case '*':
			step, i, err = parseStar(s, i+1)
		case '{':
			step, i, err = parseNumberRange(s, i+1)
		case '(':
			step, i, err = parseOptions(s, i+1)
		default:
			step, i, err = parseExactStr(s, i)
		}
		if err != nil {
			return Pattern{}, err
		}
		steps = append(steps, step)
	}
	return Pattern{steps}, nil
}

func parseStar(s string, i int) (patternStep, int, error) {
	step := star{}
	if i < len(s) {
		if isSpecial(s[i]) {
			return nil, 0, parsingError(s, "\"*\" should not be directly before special character %q", s[i])
		}
		step.until = s[i]
		i++
	}
	return step, i, nil
}

func parseNumberRange(s string, i int) (patternStep, int, error) {
	step := numberRange{}

	num, i, err := parseNumber(s, i)
	if err != nil {
		return nil, 0, err
	}
	step.min = num

	if s[i] != '.' || len(s) == i+1 || s[i+1] != '.' {
		return nil, 0, parsingError(s, "invalid number range")
	}
	i += 2

	num, i, err = parseNumber(s, i)
	if err != nil {
		return nil, 0, err
	}
	step.max = num

	if s[i] != '}' {
		return nil, 0, parsingError(s, "invalid number range")
	}
	i++

	return step, i, nil
}

func parseNumber(s string, i int) (int, int, error) {
	digitsCount := indexNoDigit(s[i:])
	if digitsCount == -1 {
		return 0, 0, parsingError(s, "invalid number range")
	}
	numStr := s[i : i+digitsCount]
	num, _ := strconv.Atoi(numStr)
	return num, i + digitsCount, nil
}

func parseOptions(s string, i int) (patternStep, int, error) {
	step := options{}
	for i < len(s) && s[i] != ')' {
		startIdx := i
		for i < len(s) && s[i] != '|' && s[i] != ')' {
			i++
		}
		option := exactString(s[startIdx:i])
		if option == "" {
			return nil, 0, parsingError(s, "first option must not be empty")
		}
		step = append(step, option)
		if s[i] == '|' {
			i++
		}
	}
	if i == len(s) {
		return nil, 0, parsingError(s, "options miss closing ')'")
	}
	i++
	return step, i, nil
}

func parseExactStr(s string, i int) (patternStep, int, error) {
	startIdx := i
	for i < len(s) && s[i] != '(' && s[i] != '*' && s[i] != '{' {
		i++
	}
	step := exactString(s[startIdx:i])
	return step, i, nil
}

func indexNoDigit(s string) int {
	return strings.IndexFunc(s, func(r rune) bool {
		return !unicode.IsDigit(r)
	})
}

func isSpecial(b byte) bool {
	const specialChars = "*(){}|"
	return strings.IndexByte(specialChars, b) != -1
}

func parsingError(input, messageFormat string, args ...interface{}) error {
	message := fmt.Sprintf(messageFormat, args...)
	return fmt.Errorf("invalid pattern \"%v\": %v", input, message)
}
