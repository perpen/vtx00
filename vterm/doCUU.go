package vterm

/*
  CUU
  Cursor Up Ps Times (default = 1)

  CSI Ps A  Cursor Up Ps Times (default = 1) (CUU).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	"testing"
)

func doCUU(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	n := params[0]
	cx, cy := term.cursor()
	cy = max(cy-n, 0)
	term.cursorTo(cx, cy)
	return emptyZone
}

func testCUU(t *testing.T) {
	testImpl(
		"CUU",
		"Don't move if already at the top",
		[]int{},
		testState{
			visualScreen: `a^b c
			               d e f
			               g h i`,
		},
		testState{
			visualScreen: `a^b c
			               d e f
			               g h i`,
		},
		emptyZone,
		t,
	)

	testImpl(
		"CUU",
		"Move cursor up with default 1",
		[]int{},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a^b c
			               d e f
			               g h i`,
		},
		emptyZone,
		t,
	)

	testImpl(
		"CUU",
		"Move cursor up by 2",
		[]int{2},
		testState{
			visualScreen: `a b c
			               d e f
			               g^h i`,
		},
		testState{
			visualScreen: `a^b c
			               d e f
			               g h i`,
		},
		emptyZone,
		t,
	)

	testImpl(
		"CUU",
		"Move cursor up by too much",
		[]int{5},
		testState{
			visualScreen: `a b c
			               d e f
			               g^h i`,
		},
		testState{
			visualScreen: `a^b c
			               d e f
			               g h i`,
		},
		emptyZone,
		t,
	)
}
