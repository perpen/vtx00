package vterm

/*
  CR
  CR

  Carriage Return (Ctrl-M).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	"testing"
)

func doCR(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	_, cy := term.cursor()
	term.cursorTo(0, cy)
	return emptyZone
}

func testCR(t *testing.T) {
	testImpl(
		"CR",
		"Cursor moves to leftest column",
		[]int{},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			              ^d e f
			               g h i`,
		},
		emptyZone,
		t,
	)
}
