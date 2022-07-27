package vterm

/*
  ED
  Erase in Display

  CSI Ps J  Erase in Display (ED).
  Ps = 0  -> Erase Below (default)
  Ps = 1  -> Erase Above
  Ps = 2  -> Erase All
  Ps = 3  -> Erase Saved Lines (xterm)

From Ecma-048:
If the DEVICE COMPONENT SELECT MODE (DCSM) is set to PRESENTATION, ED causes some or
all character positions of the active page (the page which contains the active presentation position in the
presentation component) to be put into the erased st ate, depending on the parameter values:
0 the active presentation position and the character positions up to the end of the page are put into the
   erased state
1 the character positions from the beginning of the page up to and including the active presentation
   position are put into the erased state
2 all character positions of the page are put into the erased state
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"testing"
)

func doED(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	cx, cy := term.cursor()
	w, h := term.width(), term.height()
	scr := term.Screen
	var rect Rect
	valid := true

	mode := params[0]
	switch mode {
	case 0: // Erase Below
		scr.ClearOnLine(cx, cy, w-cx)
		if cy < h-1 {
			scr.ClearLines(cy+1, h-1-cy)
		}
		rect = Rect{0, cy, w, h - cy}
	case 1: // Erase Above
		scr.ClearOnLine(0, cy, cx+1)
		if cy > 0 {
			scr.ClearLines(0, cy)
		}
		rect = Rect{0, 0, w, h - cy}
	case 2: // Erase All
		scr.ClearLines(0, h)
		rect = Rect{0, 0, w, h}
	case 3: // Erase Saved Lines (xterm)
		log.Warn("ED mode not implemented")
		valid = false
	default:
		log.Warn("ED: invalid mode")
		valid = false
	}

	if !valid {
		return emptyZone
	}
	return rect
}

func testED(t *testing.T) {
	testImpl(
		"ED",
		"Erase below if no param",
		[]int{},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               d^. .
			               . . .`,
		},
		Rect{0, 1, 3, 2},
		t,
	)
	testImpl(
		"ED",
		"Erase below if param 0",
		[]int{0},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               d^. .
			               . . .`,
		},
		Rect{0, 1, 3, 2},
		t,
	)
	testImpl(
		"ED",
		"Erase above if param 1",
		[]int{1},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `. . .
			               .^. f
			               g h i`,
		},
		Rect{0, 0, 3, 2},
		t,
	)
	testImpl(
		"ED",
		"Erase all if param 2",
		[]int{2},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `. . .
			               .^. .
			               . . .`,
		},
		Rect{0, 0, 3, 3},
		t,
	)
}
