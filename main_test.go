// We cannot import main package, so we place tests in the same package.
package main

import (
	"testing"
)

func TestMove(t *testing.T) {
	var state = NewState()
	var newState, spells = Move(&state, 3)
	if newState.currentPos != 3 {
		t.Errorf("Move(&state, 3) = %d; want 3", newState.currentPos)
	}
	if spells != ">>>" {
		t.Errorf("Move(&state, 3) = %s; want >>>", spells)
	}
}
