package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// func main() {
// 	var scanner = bufio.NewScanner(os.Stdin)
// 	// magicPhrase contains between 1 and 500 characters.
// 	scanner.Buffer(make([]byte, 500), 500)
// 	scanner.Scan()
// 	var magicPhrase = scanner.Text()

// 	if len(magicPhrase) == 0 {
// 		fmt.Println()
// 		return
// 	}

// 	var magicPhraseAsIntegers = MagicPhraseToIntegers(magicPhrase)
// 	for _, letter := range magicPhraseAsIntegers {
// 		if letter < 0 || letter > 26 {
// 			fmt.Println("Invalid input")
// 			return
// 		}
// 	}

// 	var instructions = Solution(magicPhraseAsIntegers)

// 	fmt.Println(instructions)
// }

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

// func Solution(magicPhrase []int) string {
// 	var repeats = CompressMagicPhrase(magicPhrase)
// 	var state State
// 	var instructions string
// 	if len(repeats) < 100 {
// 		state = NewState()
// 	} else {
// 		state, instructions = ShuffledState()
// 	}

// 	// TODO: remove repeats with smaller count and longer letters

// 	for _, repeat := range repeats {
// 		instructions += reachSpells(&state, repeat.letter)
// 		instructions += triggerSpells(&state, repeat.count)
// 	}

// 	return instructions
// }

////////////////////
// Data types
////////////////////

type Letter = int

const LETTERS = 27

// Returns the offset when rolling forward and backward.
func letterOffsets(a, b Letter) (int, int) {
	var offset = b - a
	if offset < 0 {
		offset = offset + LETTERS
	}
	return offset, LETTERS - offset
}

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

////////////////////
// State type
////////////////////

type State struct {
	zones      []Letter
	currentPos int
}

func NewState() State {
	return State{zones: make([]Letter, 30), currentPos: 0}
}

func (state *State) Clone() State {
	var newZones = make([]Letter, len(state.zones))
	copy(newZones, state.zones)
	return State{zones: newZones, currentPos: state.currentPos}
}

func (state *State) CurrentLetter() Letter {
	return state.zones[state.currentPos]
}

func (state *State) CurrentZone() *Letter {
	return &state.zones[state.currentPos]
}

// Returns the letter at `pos` position.
// Supports `-len(state.zones) <= pos < 2 * len(state.zones)`.
func (state *State) Zone(pos int) Letter {
	if pos < 0 {
		pos = pos + len(state.zones)
	} else if pos >= len(state.zones) {
		pos = pos - len(state.zones)
	}
	return state.zones[pos]
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
// Area type
////////////////////

// Directional area within a zone
type Area struct {
	start  int
	dir    int
	length int
}

func (area *Area) At(state *State, i int) Letter {
	return state.Zone(area.start + i*area.dir)
}

////////////////////
// Simple spells
////////////////////

// Move to `pos` position.
// Returns the state after the action and the spells to act.
func (state *State) MoveTo(pos int) (State, Spells) {
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

// Move forward in the `dir` direction.
// Returns the state after the action and the spells to act.
func (state *State) MoveForward(dir int) (State, Spells) {
	var spells Spells
	if dir > 0 {
		spells = ">"
	} else {
		spells = "<"
	}

	var pos = state.currentPos + dir
	if pos < 0 {
		pos = pos + len(state.zones)
	} else if pos >= len(state.zones) {
		pos = pos - len(state.zones)
	}

	var nextState = State{zones: state.zones, currentPos: pos}
	return nextState, spells
}

// Move backward in the `dir` direction.
// Returns the state after the action and the spells to act.
func (state *State) MoveBackward(dir int) (State, Spells) {
	return state.MoveForward(-dir)
}

// Roll to the `letter`.
// Returns the state after the action and the spells to act.
func RollTo(state *State, letter Letter) (State, Spells) {
	var baseLetter = state.CurrentLetter()
	var offsetForward, offsetBackward = letterOffsets(baseLetter, letter)

	var newState = state.Clone()
	*state.CurrentZone() = letter

	// Roll in the shorter direction
	if offsetForward <= offsetBackward {
		return newState, strings.Repeat("+", offsetForward)
	} else {
		return newState, strings.Repeat("-", offsetBackward)
	}
}

////////////////////
// Combined spells
////////////////////

// Sequence type for alignment
type Sequence struct {
	letters []Letter
	dir     int
}

func NewSequence(letters []Letter, dir int) Sequence {
	return Sequence{letters: letters, dir: dir}
}

func (seq *Sequence) Len() int {
	return len(seq.letters)
}

func (seq *Sequence) At(i int) Letter {
	if seq.dir > 0 {
		return seq.letters[i]
	} else {
		return seq.letters[len(seq.letters)-1-i]
	}
}

// Roll the characters in `area` to match `seq`.
func Align(state *State, area Area, seq Sequence) (State, Spells) {
	// We assume seq is short, so we use arranging left to right approach simply.

	var currentState = *state
	var spells Spells

	var nextPos = area.start
	for i := 0; i < seq.Len(); i++ {
		// early return if arranged
		if match(&currentState, area, seq) {
			return currentState, spells
		}

		var targetLetter = seq.At(i)

		// move to next
		currentState, moveSpells := currentState.MoveTo(nextPos)
		spells += moveSpells
		nextPos += area.dir

		// roll letter
		currentState, rollSpells := RollTo(&currentState, targetLetter)
		spells += rollSpells
	}

	return currentState, spells
}

// Returns whether the characters in `area` match `seq`
func match(state *State, area Area, seq Sequence) bool {
	for i := 0; i < seq.Len(); i++ {
		var targetLetter = seq.At(i)
		var zoneLetter = area.At(state, i)
		if zoneLetter != targetLetter {
			return false
		}
	}
	return true
}

////////////////////
// Trigger
////////////////////

// Trigger the characters in area count times.
func Trigger(state *State, area Area, count int) (State, Spells) {
	// TODO: align "abbc" as "abc" and trigger ".>..>."

	// ループカウンタがらんだむな初期値

	var currentState, spells = state.Move(start)

}

func triggerRepeatWithLoop(state *State, start int, dir int, length int) (State, Spells) {
	// TODO: double loop
	if currentState.Zone(currentState.currentPos) != 0 {
		var maxRoleDir int
		if currentState.Zone(currentState.currentPos) > 13 {
			maxRoleDir = -1
		} else {
			maxRoleDir = 1
		}
	}
}

func triggerSequenceWithLoop(state *State, dir int, length int, rollDir int) (State, Spells) {
	var spells Spells = "["
	var currentState = *state

	spells += triggerAndReturn(&currentState, dir, length)
}

// Trigger in the direction of dir for length times from start, and return to the original position
func triggerWithoutLoop(state *State, dir int, length int) (State, Spells) {
	var currentState = *state
	var spells Spells

	// move to start
	currentState, additionalSpells := currentState.Move(start)
	spells += additionalSpells

	for i := 0; i < length; i++ {
		currentState, additionalSpells = triggerAndReturn(&currentState, dir, length)
		spells += additionalSpells
	}

	return currentState, spells
}

// Trigger in the direction of dir for length times, and return to the original position
func triggerAndReturn(state *State, dir int, length int) (State, Spells) {
	// TODO: palindrome
	var currentState = *state
	var spells Spells = "."

	for i := 1; i < length; i++ {
		// move to next
		var additionalSpells Spells
		currentState, additionalSpells = Move(&currentState, currentState.currentPos+dir)
		spells += additionalSpells

		// trigger
		spells += "."
	}

	for i := 1; i < length; i++ {
		// move to previous
		var additionalSpells Spells
		currentState, additionalSpells = Move(&currentState, currentState.currentPos-dir)
		spells += additionalSpells
	}

	return currentState, spells
}

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
