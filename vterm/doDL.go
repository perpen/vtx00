package vterm

/*
  DL
  Delete Ps Line(s)

  CSI Ps M  Delete Ps Line(s) (default = 1) (DL).

From Ecma-048:
If the DEVICE COMPONENT SELECT MODE (DCSM) is set to PRESENTATION, DL causes the
contents of the active line (the line that contains the active presentation position) and, depending on the
setting of the LINE EDITING MODE (VEM), the contents of the n-1 preceding or following lines to be
removed from the presentation component, where n equals the value of Pn. The resulting gap is closed
by shifting the contents of a number of adjacent lines towa rds the active line. At the other end of the
shifted part, n lines are put into the erased state.
The active presentation position is moved to the line home position in the active line. The line home
position is established by the parameter value of SET LINE HOME (SLH). If the TABULATION STOP
MODE (TSM ) is set to SINGLE, character tabulation stops are cleared in the lines that are put into the
erased state.
The extent of the shifted part is established by SELECT EDITING EXTENT (SEE).
Any occurrences of the start or end of a selected area, the start or nd of a qualified area, or a tabulation
stop in the shifted part, are also shifted.
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"testing"
)

func doDL(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	_, cy := term.cursor()
	w, h := term.width(), term.height()
	scr := term.Screen
	linesToBottom := h - cy
	count := min(params[0], linesToBottom)
	// log.Info("doDL: count: ", count, " linesToBottom: ", linesToBottom)
	if count == 0 {
		log.Print("NOOO")
		return emptyZone
	}
	if linesToBottom == count {
		// log.Info("doDL: clearing ", 0, cy, count)
		scr.ClearLines(cy, count)
		return Rect{0, cy, w, count}
	}

	linesToMove := min(h - (cy + count), linesToBottom)
	scr.MoveLines(cy+count, linesToMove, -count)
	scr.ClearLines(cy+linesToMove, count)
	return Rect{0, cy, w, count+1}
}

func testDL(t *testing.T) {
	testImpl(
		"DL",
		"Without param, delete current line",
		[]int{},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
			settings: defaultSettings,
		},
		testState{
			visualScreen: `a b c
			               g^h i
			               . . .`,
			settings: defaultSettings,
		},
		Rect{0, 1, 3, 2},
		t,
	)
	testImpl(
		"DL",
		"Param 0, do nothing",
		[]int{0},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
			settings: defaultSettings,
		},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
			settings: defaultSettings,
		},
		emptyZone,
		t,
	)
	testImpl(
		"DL",
		"Param 2, delete current and next lines",
		[]int{2},
		testState{
			visualScreen: `a^b c
			               d e f
			               g h i`,
			settings: defaultSettings,
		},
		testState{
			visualScreen: `g^h i
			               . . .
			               . . .`,
			settings: defaultSettings,
		},
		Rect{0, 0, 3, 3},
		t,
	)
	testImpl(
		"DL",
		"Param 2 from middle",
		[]int{2},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i
			               j k l`,
			settings: defaultSettings,
		},
		testState{
			visualScreen: `a b c
			               j^k l
			               . . .
			               . . .`,
			settings: defaultSettings,
		},
		Rect{0, 1, 3, 3},
		t,
	)
	testImpl(
		"DL",
		"Param 2 to end",
		[]int{2},
		testState{
			visualScreen: `a b c
			              ^d e f
			               g h i`,
			settings: defaultSettings,
		},
		testState{
			visualScreen: `a b c
			              ^. . .
			               . . .`,
			settings: defaultSettings,
		},
		Rect{0, 1, 3, 2},
		t,
	)
}
