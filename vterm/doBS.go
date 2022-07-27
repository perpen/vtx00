package vterm

/*
  BS
  BS

  Backspace (Ctrl-H).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	"testing"
)

func doBS(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	cx, cy := term.cursor()
	cx = max(cx-1, 0)
	term.cursorTo(cx, cy)
	return emptyZone
}

func testBS(t *testing.T) {
	testImpl(
		"BS",
		"Cursor goes left",
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
	testImpl(
		"BS",
		"Cursor unchanged if on first column already",
		[]int{},
		testState{
			visualScreen: `a b c
			              ^d e f
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
