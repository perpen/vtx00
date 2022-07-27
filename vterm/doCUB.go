package vterm

/*
  CUB
  Cursor Backward Ps Times (default = 1)

  CSI Ps D  Cursor Backward Ps Times (default = 1) (CUB).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	"testing"
)

func doCUB(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	n := params[0]
	cx, cy := term.cursor()
	cx = max(cx-n, 0)
	term.cursorTo(cx, cy)
	return emptyZone
}

func testCUB(t *testing.T) {
	testImpl(
		"CUB",
		"Without param, go left by 1",
		[]int{},
		testState{
			visualScreen: `a b c
			               d e^f
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
		"CUB",
		"Go left by 2",
		[]int{2},
		testState{
			visualScreen: `a b c
			               d e^f
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
		"CUB",
		"Go to leftest column if param too high",
		[]int{9},
		testState{
			visualScreen: `a b c
			               d e^f
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
