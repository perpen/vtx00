package vterm

/*
  CUF
  Cursor Forward Ps Times (default = 1)

  CSI Ps C  Cursor Forward Ps Times (default = 1) (CUF).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	"testing"
)

func doCUF(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	n := params[0]
	cx, cy := term.cursor()
	cx = min(cx+n, term.width()-1)
	term.cursorTo(cx, cy)
	return emptyZone
}

func testCUF(t *testing.T) {
	testImpl(
		"CUF",
		"Cursor moves forward by 1 if no param",
		[]int{},
		testState{
			visualScreen: `a b c
			              ^d e f
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
		"CUF",
		"Cursor moves forward by 2",
		[]int{2},
		testState{
			visualScreen: `a b c
			              ^d e f
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
		"CUF",
		"Cursor to rightest column if param too high",
		[]int{9},
		testState{
			visualScreen: `a b c
			              ^d e f
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
