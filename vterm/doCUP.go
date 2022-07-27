package vterm

/*
  CUP
  Cursor Position

  CSI Ps ; Ps H
Cursor Position [row;column] (default = [1,1]) (CUP).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	"testing"
)

func doCUP(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	cx, cy := term.cursor()
	cx = min(params[1]-1, term.width()-1)
	cy = min(params[0]-1, term.height()-1)
	term.cursorTo(cx, cy)
	return emptyZone
}

func testCUP(t *testing.T) {
	testImpl(
		"CUP",
		"Params default to 1, 1",
		[]int{},
		testState{
			visualScreen: `a b c
			               d e f
			               g^h i`,
		},
		testState{
			visualScreen: `^a b c
			                d e f
			                g h i`,
		},
		emptyZone,
		t,
	)
	testImpl(
		"CUP",
		"With params",
		[]int{2, 2},
		testState{
			visualScreen: `^a b c
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
		"CUP",
		"To bottom-right if params too high",
		[]int{9, 9},
		testState{
			visualScreen: `^a b c
			                d e f
			                g h i`,
		},
		testState{
			visualScreen: `a b c
			               d e f
			               g h^i`,
		},
		emptyZone,
		t,
	)
}
