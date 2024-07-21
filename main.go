package main

import (
	"bufio"
	"fmt"
	"os"
)

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	scanner.Scan()
	magicPhrase := scanner.Text()
	_ = magicPhrase // to avoid unused error

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Println("+.>-.") // Write action to stdout
}
