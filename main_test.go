// We cannot import main package, so we place tests in the same package.
package main

import (
	"testing"
)

func TestMove(t *testing.T) {
	var tests = []struct {
		name       string
		currentPos int
		nextPos    int
		spells     string
	}{
		{"left to right", 0, 3, ">>>"},
		{"left to right circle", 28, 0, ">>"},
		{"right to left", 4, 1, "<<<"},
		{"right to left circle", 1, 29, "<<"},
	}

	for _, tc := range tests {
		var state = NewState()
		state.currentPos = tc.currentPos
		var newState, spells = Move(&state, tc.nextPos)
		if newState.currentPos != tc.nextPos {
			t.Errorf("Move(&state, %d) = %d; want %d", tc.nextPos, newState.currentPos, tc.nextPos)
		}
		if spells != tc.spells {
			t.Errorf("Move(&state, %d) = %s; want %s", tc.nextPos, spells, tc.spells)
		}
	}
}
