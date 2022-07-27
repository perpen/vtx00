package vterm

/*
  CUD
  Cursor Down Ps Times (default = 1)

  CSI Ps B  Cursor Down Ps Times (default = 1) (CUD).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	"testing"
)

func doCUD(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	n := params[0]
	cx, cy := term.cursor()
	cy = min(cy+n, term.height()-1)
	term.cursorTo(cx, cy)
	return emptyZone
}

func testCUD(t *testing.T) {
	testImpl(
		"CUD",
		"Cursor moves down by 1 if no params",
		[]int{},
		testState{
			visualScreen: `a^b c
			               d e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		emptyZone,
		t,
	)
	testImpl(
		"CUD",
		"Cursor moves down by 2",
		[]int{2},
		testState{
			visualScreen: `a^b c
			               d e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               d e f
			               g^h i`,
		},
		emptyZone,
		t,
	)
	testImpl(
		"CUD",
		"Cursor goes to lowest row if param too high",
		[]int{9},
		testState{
			visualScreen: `a^b c
			               d e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               d e f
			               g^h i`,
		},
		emptyZone,
		t,
	)
}
