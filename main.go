package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	// magicPhrase contains between 1 and 500 characters.
	scanner.Buffer(make([]byte, 500), 500)
	scanner.Scan()
	magicPhrase := scanner.Text()

	magicPhraseAsIntegers := magicPhraseToIntegers(magicPhrase)
	instructions := solution(magicPhraseAsIntegers)

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

func solution(magicPhrase []int) string {
	var instructions string

	for _, letter := range magicPhrase {
		if letter == 0 {
			instructions += "."
		} else if letter <= 13 {
			for i := 0; i < letter; i++ {
				instructions += "+"
			}
			instructions += "."
			for i := 0; i < letter; i++ {
				instructions += "-"
			}
		} else {
			for i := 0; i < 26-letter+1; i++ {
				instructions += "-"
			}
			instructions += "."
			for i := 0; i < 26-letter+1; i++ {
				instructions += "+"
			}
		}
	}

	return instructions
}
