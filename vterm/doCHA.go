package vterm

/*
  CHA
  Cursor Character Absolute  [column] (default = [row,1])

  CSI Ps G  Cursor Character Absolute  [column] (default = [row,1]) (CHA).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	"testing"
)

func doCHA(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	cx, cy := term.cursor()
	cx = min(params[0]-1, term.width()-1)
	term.cursorTo(cx, cy)
	return emptyZone
}

func testCHA(t *testing.T) {
	testImpl(
		"CHA",
		"Cursor moves to specified column",
		[]int{3},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               d e^f
			               g h i`,
		},
		emptyZone,
		t,
	)
	testImpl(
		"CHA",
		"Cursor moves to rightest column if param too high",
		[]int{9},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               d e^f
			               g h i`,
		},
		emptyZone,
		t,
	)
}
