package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
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

	var magicPhraseAsIntegers = MagicPhraseToIntegers(magicPhrase)
	for _, letter := range magicPhraseAsIntegers {
		if letter < 0 || letter > 26 {
			fmt.Println("Invalid input")
			return
		}
	}

	var instructions = Solution(magicPhraseAsIntegers)

	fmt.Println(instructions)
}

func MagicPhraseToIntegers(magicPhrase string) []int {
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

func Solution(magicPhrase []int) string {
	var repeats = CompressMagicPhrase(magicPhrase)
	var state State
	var instructions string
	if len(repeats) < 100 {
		state = NewState()
	} else {
		state, instructions = ShuffledState()
	}

	// TODO: remove repeats with smaller count and longer letters

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

// A letter sequence and the number of times it is repeated.
type Repeat struct {
	letters []Letter
	count   int
}

// Returns the repeats of letters length up to `n`.
func NextRepeats(letters []Letter, n int) []Repeat {
	// zones count - 2 = 28
	n = min(n, len(letters), 28)
	var repeats = make([]Repeat, n)
	for l := 1; l <= n; l++ {
		var seq = letters[:l]
		var count = 1
		for j := l; j < len(letters); j += l {
			if !reflect.DeepEqual(seq, letters[j:j+l]) {
				break
			}
			count++
		}
		repeats[l] = Repeat{letters: seq, count: count}
	}
	return repeats
}

type State struct {
	zones      []Letter
	currentPos int
}

func NewState() State {
	return State{zones: make([]Letter, 30), currentPos: 0}
}

func ShuffledState() (State, Spells) {
	var state = NewState()
	for i := 0; i < len(state.zones)/2; i++ {
		state.zones[i*2] = 26
		state.zones[i*2+1] = 13
	}
	state.zones[0] = 0
	return state, "+[<+<++]"
}

////////////////////
// Simple spells
////////////////////

// Returns the state that Blub moves to `pos` and the spells to Move.
func Move(state *State, pos int) (State, Spells) {
	// left to right
	var offset = pos - state.currentPos
	if offset < 0 {
		offset = offset + len(state.zones)
	}
	var newState = State{zones: state.zones, currentPos: pos}
	if offset <= len(state.zones)/2 {
		return newState, strings.Repeat(">", offset)
	} else {
		return newState, strings.Repeat("<", len(state.zones)-offset)
	}
}

// Returns the state that Blub rolls the letter on the zone to `letter` and the spells to Roll.
func Roll(state *State, letter Letter) (State, Spells) {
	var baseLetter = state.zones[state.currentPos]
	var letterOffset = letter - baseLetter
	if letterOffset < 0 {
		letterOffset = letterOffset + LETTERS
	}

	var newZones = make([]Letter, len(state.zones))
	copy(newZones, state.zones)
	newZones[state.currentPos] = letter
	var newState = State{zones: newZones, currentPos: state.currentPos}

	if letterOffset <= LETTERS/2 {
		return newState, strings.Repeat("+", letterOffset)
	} else {
		return newState, strings.Repeat("-", LETTERS-letterOffset)
	}
}

////////////////////
// Combined spells
////////////////////

// Returns the state that `seq` is arranged at `seqStart` and the spells to arrange.
func arrange(state *State, seq []Letter, moveStart int, moveDir, seqDir int) (State, Spells) {
	// We assume seq is short, so we use arranging left to right approach simply.

	// move to seqStart from currentPos
	var currentState, spells = Move(state, moveStart)

	// roll letter
	var rollSpells Spells
	if seqDir > 0 {
		currentState, rollSpells = Roll(&currentState, seq[0])
	} else {
		currentState, rollSpells = Roll(&currentState, seq[len(seq)-1])
	}
	spells += rollSpells

	for i := 1; i < len(seq); i++ {
		var targetLetter Letter
		if seqDir > 0 {
			targetLetter = seq[i]
		} else {
			targetLetter = seq[len(seq)-1-i]
		}

		// move to next
		currentState, moveSpells := Move(&currentState, currentState.currentPos+moveDir)
		spells += moveSpells

		// roll letter
		currentState, rollSpells := Roll(&currentState, targetLetter)
		spells += rollSpells
	}

	return currentState, spells
}

////////////////////
// Trigger
////////////////////

// Spells to trigger the current letter.
// e.g. >+[<.>+]< for A * 26
func triggerSpells(state *State, triggerCount int) string {
	// TODO: trigger with loop each letter using loop
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
