package vterm

/*
  IL
  Insert Ps Line(s) (default = 1)

  CSI Ps L  Insert Ps Line(s) (default = 1) (IL).

From Ecma-048:
If the DEVICE COMPONENT SELECT MODE (DCSM) is set to PRESENTATION, IL is used to
prepare the insertion of n lines, by putting into the erased st ate in the presentation component the active
line (the line that contains the active presentation position) and, depending on the setting of the LINE
EDITING MODE (VEM), the n-1 preceding or following lines, where n equals the value of Pn. The
previous contents of the active line and of adjacent lines are shifted away from the active line. The
contents of n lines at the other end of the shi fted part are rem oved. The active presentation position is
moved to the line home position in the active line. The line home position is established by the
parameter value of SET LINE HOME (SLH).
The extent of the shifted part is established by SELECT EDITING EXTENT (SEE).
Any occurrences of the start or end of a selected area, the start or nd of a qualified area, or a tabulation
stop in the shifted part, are also shifted.
If the TABULATION STOP MODE (TSM ) is set to SINGLE, character tabulation stops are cleared in
the lines that are put in to the erased state.
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"testing"
)

func doIL(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	n := params[0]
	cx, cy := term.cursor()
	w, _ := term.width(), term.height()
	settings := term.settings
	scr := term.Screen

	if offRegion(cy, term) {
		log.Warn("IL: called from off-region: cursor:", cx, cy,
			" region:", settings.regtop, settings.regbot)
		return emptyZone
	}
	bot := settings.regbot
	regHeight := bot - settings.regtop
	scr.MoveLines(cy, regHeight-n+1, n)
	scr.ClearLines(cy, min(n, bot - cy+1))
	return Rect{0, cy, w, bot - cy + 1}
}

// FIXME move to some util file, or maybe to Term
func offRegion(y int, term *Term) bool {
	return y < term.settings.regtop || y > term.settings.regbot
}

func testIL(t *testing.T) {
	testImpl(
		"IL",
		"No param, insert 1 line",
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
		"IL",
		"Insert 1 line from top",
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
		"IL",
		"Insert 1 line from middle",
		[]int{1},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
		                   .^. .
			               d e f`,
		},
		Rect{0, 1, 3, 2},
		t,
	)
	testImpl(
		"IL",
		"Insert many lines",
		[]int{3},
		testState{
			visualScreen: `a b c
			               d e f
			               g^h i`,
		},
		testState{
			visualScreen: `a b c
			               d e f
			               .^. .`,
		},
		Rect{0, 2, 3, 1},
		t,
	)
	testImpl(
		"IL",
		"Insert 2 lines from scroll region",
		[]int{2},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i
			               j k l
			               m n o`,
			settings: TermSettings{
				regtop: 1,
				regbot: 3,
			},
		},
		testState{
			visualScreen: `a b c
			               .^. .
			               . . .
			               d e f
			               m n o`,
			settings: TermSettings{
				regtop: 1,
				regbot: 3,
			},
		},
		Rect{0, 1, 3, 3},
		t,
	)
	testImpl(
		"IL",
		"Insert many lines",
		[]int{9},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               .^. .
			               . . .`,
		},
		Rect{0, 1, 3, 2},
		t,
	)
	testImpl(
		"IL",
		"Insert many lines from scroll region",
		[]int{9},
		testState{
			visualScreen: `a b c
			               d e f
			               g^h i
			               j k l
			               m n o`,
			settings: TermSettings{
				regtop: 1,
				regbot: 3,
			},
		},
		testState{
			visualScreen: `a b c
			               d e f
			               .^. .
			               . . .
			               m n o`,
			settings: TermSettings{
				regtop: 1,
				regbot: 3,
			},
		},
		Rect{0, 2, 3, 2},
		t,
	)
	testImpl(
		"IL",
		"Do nothing if cursor off the scrolling region",
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
			               d e f
			               g^h i`,
			settings: TermSettings{
				regtop: 0,
				regbot: 1,
			},
		},
		emptyZone,
		t,
	)
}
