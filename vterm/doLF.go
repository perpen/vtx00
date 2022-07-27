package vterm

/*
  LF
  LF

  Line Feed or New Line (NL).  (LF  is Ctrl-J).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"testing"
)

func doLF(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	cx, cy := term.cursor()
	w, h := term.width(), term.height()
	if cy == term.settings.regbot {
		log.Infof("doLF: regtop=%v regbot=%v",
			term.settings.regtop, term.settings.regbot)
		regHeight := term.settings.regbot - term.settings.regtop + 1
		term.Screen.MoveLines(term.settings.regtop+1, regHeight-1, -1)
		return Rect{0, term.settings.regtop, w, regHeight}
	} else {
		cy = min(h-1, cy+1)
		term.cursorTo(cx, cy)
		return emptyZone
	}
}

func testLF(t *testing.T) {
	testImpl(
		"LF",
		"Move cursor down",
		[]int{},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
		},
		testState{
			visualScreen: `a b c
			               d e f
			               g^h i`,
		},
		emptyZone,
		t,
	)

	testImpl(
		"LF",
		"Scroll up if on last line",
		[]int{},
		testState{
			visualScreen: `a b c
			               d e f
			               g^h i`,
		},
		testState{
			visualScreen: `d e f
			               g h i
			               .^. .`,
		},
		Rect{0, 0, 3, 3},
		t,
	)

	testImpl(
		"LF",
		"Scroll up region if on last line",
		[]int{},
		testState{
			visualScreen: `a b c
			               d e f
			               g^h i
			               j k l`,
			settings: TermSettings{
				regtop: 1,
				regbot: 2,
			},
		},
		testState{
			visualScreen: `a b c
			               g h i
			               .^. .
			               j k l`,
			settings: TermSettings{
				regtop: 1,
				regbot: 2,
			},
		},
		Rect{0, 1, 3, 2},
		t,
	)

	testImpl(
		"LF",
		"Move cursor down",
		[]int{},
		testState{
			visualScreen: `a^b c
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
}
