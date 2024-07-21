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

	var instructions string
	var currentLetter = 0
	for _, repeat := range repeats {
		instructions += roleSpells(currentLetter, repeat.letter)
		instructions += triggerSpells(repeat.count)

		currentLetter = repeat.letter
	}

	return instructions
}

type Letter = int

// A letter and the number of times it is repeated.
type Repeat struct {
	letter Letter
	count  int
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
// Roll
////////////////////

const LETTERS = 27

func roleSpells(current Letter, target Letter) string {
	var rolePlusCount = calcRollPlusCount(current, target)
	var roleMinusCount = calcRollMinusCount(current, target)

	if rolePlusCount < roleMinusCount {
		return calcRolePlusSpells(current, target)
	} else {
		return calcRoleMinusSpells(current, target)
	}
}

// Calculate the number of steps to roll the current letter to the target letter using only +.
// e.g. ++++ for A to E
func calcRollPlusCount(current Letter, target Letter) int {
	var offset = target - current
	if offset < 0 {
		offset = LETTERS + offset
	}
	return offset
}

// Calculate the spells to roll the current letter to the target letter using only +.
// e.g. ++++ for A to E
func calcRolePlusSpells(current Letter, target Letter) string {
	var offset = target - current
	if offset < 0 {
		offset = LETTERS + offset
	}
	return strings.Repeat("+", offset)
}

// Calculate the number of steps to roll the current letter to the target letter using only -.
// e.g. ---- for E to A
func calcRollMinusCount(current Letter, target Letter) int {
	var offset = current - target
	if offset < 0 {
		offset = LETTERS + offset
	}
	return offset
}

// Calculate the spells to roll the current letter to the target letter using only -.
// e.g. ---- for E to A
func calcRoleMinusSpells(current Letter, target Letter) string {
	var offset = current - target
	if offset < 0 {
		offset = LETTERS + offset
	}
	return strings.Repeat("-", offset)
}

////////////////////
// Trigger
////////////////////

// Spells to trigger the current letter.
func triggerSpells(triggerCount int) string {
	var triggerNaiveCount = calcTriggerNaiveCount(triggerCount)
	var triggerLoopCount = calcTriggerLoopCount(triggerCount)

	if triggerNaiveCount < triggerLoopCount {
		return calcTriggerNaiveSpells(triggerCount)
	} else {
		return calcTriggerLoopSpells(triggerCount)
	}
}

// Calculate the number of steps to trigger the current letter.
// e.g. ... for A * 3
func calcTriggerNaiveCount(triggerCount int) int {
	return triggerCount
}

// Calculate the spells to trigger the current letter.
// e.g. ... for AAA
func calcTriggerNaiveSpells(triggerCount int) string {
	return strings.Repeat(".", triggerCount)
}

// Calculate the number of steps to trigger the current letter using loop.
// e.g. >+[<.>+]< for A * 26
func calcTriggerLoopCount(triggerCount int) int {
	var loopCount = triggerCount / 26
	var restCount = triggerCount % 26

	// >
	var count = 1
	// +[<.>+]
	count += 7 * loopCount
	// <
	count += 1
	// ...
	count += restCount

	return count
}

// Calculate the spells to trigger the current letter using loop.
// e.g. >+[<.>+]< for A * 26
func calcTriggerLoopSpells(triggerCount int) string {
	var loopCount = triggerCount / 26
	var restCount = triggerCount % 26

	var spells = ">"
	spells += strings.Repeat("+[<.>+]", loopCount)
	spells += "<"
	spells += strings.Repeat(".", restCount)

	return spells
}

// Debug print
func ePrintln(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
}
