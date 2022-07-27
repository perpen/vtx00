package vterm

/*
  RI
  REVERSE LINE FEED

  From Ecma-048:
If the DEVICE COMPONENT SELECT MODE (DCSM) is set to PRESENTATION, RI causes the
active presentation position to be moved in the presentation component to the corresponding character
position of the preceding line.
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	"testing"
)

func doRI(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	//log.Debugln("doRI: region:", term.settings.regtop, "to", term.settings.regbot)
	//log.Debugln("doRI: cursor:", *cursor)
	cx, cy := term.cursor()
	w, _ := term.width(), term.height()
	if cy == term.settings.regtop {
		scr := term.Screen
        regHeight := term.settings.regbot - term.settings.regtop + 1
		scr.MoveLines(cy, regHeight - 1, 1)
		return Rect{0, term.settings.regtop, w, regHeight}
	} else {
		cy = max(0, cy-1)
		term.cursorTo(cx, cy)
		return emptyZone
	}
}

func testRI(t *testing.T) {
	testImpl(
		"RI",
		"Move cursor up",
		[]int{},
		testState{
			visualScreen: `a b c
			               d e f
			               g^h i`,
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
		"RI",
		"From top line, scroll down",
		[]int{},
		testState{
			visualScreen: `a^b c
			               d e f
			               g h i`,
		},
		testState{
			visualScreen: `.^. .
		                   a b c
			               d e f`,
		},
        Rect{0, 0, 3, 3},
		t,
	)
	testImpl(
		"RI",
		"From region top line, scroll region down",
		[]int{},
		testState{
			visualScreen: `a^b c
			               d e f
			               g h i`,
			settings: TermSettings{
				regtop: 0,
				regbot: 1,
			},
		},
		testState{
			visualScreen: `.^. .
		                   a b c
			               g h i`,
			settings: TermSettings{
				regtop: 0,
				regbot: 1,
			},
		},
        Rect{0, 0, 3, 2},
		t,
	)
	testImpl(
		"RI",
		"If outside region, xterm does a CUU",
		[]int{},
		testState{
			visualScreen: `a b c
			               d e f
			               g^h i`,
			settings: TermSettings{
				regtop: 0,
				regbot: 1,
			},
		},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
			settings: TermSettings{
				regtop: 0,
				regbot: 1,
			},
		},
        emptyZone,
		t,
	)
}
