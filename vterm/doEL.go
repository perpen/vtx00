package vterm

/*
  EL
  Erase in Line

  CSI Ps K  Erase in Line (EL).
  Ps = 0  -> Erase to Right (default).
  Ps = 1  -> Erase to Left.
  Ps = 2  -> Erase All.

From Ecma-048
If the DEVICE COMPONENT SELECT MODE (DCSM) is set to PRESENTATION, EL causes som e or
all character positions of the active line (the line which contains the active presentation position in the
presentation component) to be put into the erased st ate, depending on the parameter values:
0 the active presentation position and the character positions up to the end of the line are put into the
   erased state
1 the character positions from the beginning of the line up to and including the active presentation
   position are put into the erased state
2 all character positions of the lineare put into the erased state
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"testing"
)

func doEL(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	cx, cy := term.cursor()
	w, _ := term.width(), term.height()
	scr := term.Screen
	var beg, end int
	valid := true

	mode := params[0]
	switch mode {
	case 0: // from cursor to eol
		beg = cx
		end = w-1
	case 1: // from bol to cursor
		beg = 0
		end = cx
	case 2: // whole line
		beg = 0
		end = w-1
	case 3:
		log.Warn("EL mode not implemented")
		valid = false
	default:
		log.Warn("EL: invalid mode")
		valid = false
	}

	if !valid {
		return emptyZone
	}
	scr.ClearOnLine(beg, cy, end-beg+1)
	return Rect{beg, cy, end-beg+1, 1}
}

func testEL(t *testing.T) {
	testImpl(
		"EL",
		"Erase to right by default",
		[]int{},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               d^. .
			               g h i`,
		},
		Rect{1, 1, 2, 1},
		t,
	)
	testImpl(
		"EL",
		"Erase to right",
		[]int{},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               d^. .
			               g h i`,
		},
		Rect{1, 1, 2, 1},
		t,
	)
	testImpl(
		"EL",
		"Erase to left",
		[]int{1},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               .^. f
			               g h i`,
		},
		Rect{0, 1, 2, 1},
		t,
	)
	testImpl(
		"EL",
		"Erase all",
		[]int{2},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               .^. .
			               g h i`,
		},
		Rect{0, 1, 3, 1},
		t,
	)
}
