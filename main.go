package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var scanner = bufio.NewScanner(os.Stdin)
	// magicPhrase contains between 1 and 500 characters.
	scanner.Buffer(make([]byte, 500), 500)
	scanner.Scan()
	var magicPhrase = scanner.Text()

	if len(magicPhrase) == 0 {
		fmt.Println()
		return
	}

	var magicPhraseAsIntegers = magicPhraseToIntegers(magicPhrase)
	for _, letter := range magicPhraseAsIntegers {
		if letter < 0 || letter > 26 {
			fmt.Println("Invalid input")
			return
		}
	}

	var instructions = solution(magicPhraseAsIntegers)

	fmt.Println(instructions)
}

func magicPhraseToIntegers(magicPhrase string) []int {
	var magicPhraseAsIntegers []int
	for _, char := range magicPhrase {
		if char == ' ' {
			magicPhraseAsIntegers = append(magicPhraseAsIntegers, 0)
		} else {
			magicPhraseAsIntegers = append(magicPhraseAsIntegers, int(char)-int('A')+1)
		}
	}
	return magicPhraseAsIntegers
}

// Currently only uses the first and second zones
// TODO: use the entire zones
func solution(magicPhrase []int) string {
	var repeats = compressMagicPhrase(magicPhrase)
	var state State
	var instructions string
	if len(repeats) < 100 {
		state = newState()
	} else {
		state, instructions = shuffled()
	}

	for _, repeat := range repeats {
		instructions += reachSpells(&state, repeat.letter)
		instructions += triggerSpells(&state, repeat.count)
	}

	return instructions
}

////////////////////
// Data types
////////////////////

type Letter = int

const LETTERS = 27

type Spells = string

// A letter and the number of times it is repeated.
type Repeat struct {
	letter Letter
	count  int
}

type State struct {
	zones      []Letter
	currentPos int
}

func newState() State {
	return State{zones: make([]Letter, 30), currentPos: 0}
}

func shuffled() (State, Spells) {
	var state = newState()
	for i := 0; i < len(state.zones)/2; i++ {
		state.zones[i*2] = 26
		state.zones[i*2+1] = 13
	}
	state.zones[0] = 0
	return state, "+[<+<++]"
}

func compressMagicPhrase(magicPhrase []Letter) []Repeat {
	var compressedMagicPhrase []Repeat
	var currentLetter = magicPhrase[0]
	var currentCount = 1

	for _, letter := range magicPhrase[1:] {
		if letter == currentLetter {
			currentCount++
		} else {
			compressedMagicPhrase = append(compressedMagicPhrase, Repeat{letter: currentLetter, count: currentCount})
			currentLetter = letter
			currentCount = 1
		}
	}
	compressedMagicPhrase = append(compressedMagicPhrase, Repeat{letter: currentLetter, count: currentCount})

	return compressedMagicPhrase
}

////////////////////
// Reach
////////////////////

// Calculate the spells to reach the target letter and update state.
func reachSpells(state *State, letter Letter) Spells {
	var bestPos = 0
	var bestCount = LETTERS
	var bestPosDir = 1
	var bestPosOffset = 0
	var bestLetterDir = 1
	var bestLetterOffset = 0
	for pos := 0; pos < len(state.zones); pos++ {
		var currentPos = state.currentPos

		// move to pos from currentPos
		var offset = pos - currentPos
		var posDir = 1
		var letterDir = 1
		if offset < 0 {
			offset = len(state.zones) + offset
		}
		if offset > len(state.zones)/2 {
			posDir = -1
			offset = len(state.zones) - offset
		}

		// roll letter
		var baseLetter = state.zones[pos]
		var letterOffset = letter - baseLetter
		if letterOffset < 0 {
			letterOffset = LETTERS + letterOffset
		}

		if letterOffset > LETTERS/2 {
			letterDir = -1
			letterOffset = LETTERS - letterOffset
		}

		var count = offset + letterOffset
		if count < bestCount {
			bestPos = pos
			bestCount = count
			bestPosDir = posDir
			bestPosOffset = offset
			bestLetterDir = letterDir
			bestLetterOffset = letterOffset
		}
	}

	var spells Spells
	if bestPosDir > 0 {
		spells += strings.Repeat(">", bestPosOffset)
	} else {
		spells += strings.Repeat("<", bestPosOffset)
	}
	if bestLetterDir > 0 {
		spells += strings.Repeat("+", bestLetterOffset)
	} else {
		spells += strings.Repeat("-", bestLetterOffset)
	}

	state.currentPos = bestPos
	state.zones[state.currentPos] = letter

	return spells
}

////////////////////
// Trigger
////////////////////

// Spells to trigger the current letter.
// e.g. >+[<.>+]< for A * 26
func triggerSpells(state *State, triggerCount int) string {
	var bestCount = triggerCount
	var bestPosDir = 1
	var bestPos = 0
	var bestLetterDir = 1
	var bestLoopCount = 0
	var bestRestCount = 0

	for _, posDir := range []int{1, -1} {
		var pos = state.currentPos + posDir
		if pos < 0 {
			pos = pos + len(state.zones)
		} else if pos >= len(state.zones) {
			pos = pos - len(state.zones)
		}

		var initLetter = state.zones[pos]

		var letterDirs []int
		if initLetter != 0 {
			letterDirs = []int{1, -1}
		} else {
			letterDirs = []int{0}
		}
		for _, letterDir := range letterDirs {
			var initCount int
			if letterDir > 0 {
				initCount = LETTERS - initLetter
			} else if letterDir < 0 {
				initCount = initLetter
			} else {
				initCount = 0
			}

			if triggerCount-initCount < 0 {
				continue
			}
			var loopCount = (triggerCount - initCount) / 26
			var restCount = (triggerCount - initCount) % 26

			// >
			var count = 1
			if letterDir != 0 {
				// [<.>+]
				count += 6
			}
			// +[<.>+]
			count += 6 * loopCount
			// <
			count += 1
			// ...
			count += restCount

			if count < bestCount {
				bestCount = count
				bestPosDir = posDir
				bestPos = pos
				bestLetterDir = letterDir
				bestLoopCount = loopCount
				bestRestCount = restCount
			}
		}
	}

	if bestCount == triggerCount {
		return strings.Repeat(".", triggerCount)
	}

	var spells string
	var f string
	var b string
	if bestPosDir > 0 {
		f = ">"
		b = "<"
	} else {
		f = "<"
		b = ">"
	}
	spells += f
	// [<.>+]
	if bestLetterDir > 0 {
		spells += "[" + b + "." + f + "+" + "]"
		state.zones[bestPos] = 0
	} else if bestLetterDir < 0 {
		spells += "[" + b + "." + f + "-" + "]"
		state.zones[bestPos] = 0
	}
	// +[<.>+]
	spells += strings.Repeat("+"+"["+b+"."+f+"+"+"]", bestLoopCount)
	if bestLoopCount > 0 {
		state.zones[bestPos] = 0
	}
	spells += b
	spells += strings.Repeat(".", bestRestCount)

	return spells
}

// Debug print
func ePrintln(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
}
