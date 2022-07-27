package vterm

/*
  BEL
  BEL

  Bell (Ctrl-G).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	"testing"
)

func doBEL(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	// term.settings.DECCKM = true
	// term.settings.style.Fg = "yellow"
	return emptyZone
}

func testBEL(t *testing.T) {
	// FIXME should set a flag? or counter?
	testImpl(
		"BEL",
		"no effect",
		[]int{},
		testState{},
		testState{},
		emptyZone,
		t,
	)
}
