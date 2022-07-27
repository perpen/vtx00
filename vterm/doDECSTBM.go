package vterm

/*
  DECSTBM
  Set Scrolling Region

  CSI Ps ; Ps r
Set Scrolling Region [top;bottom] (default = full size of win-
dow) (DECSTBM).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"testing"
)

func doDECSTBM(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	settings := &term.settings
	regtop := params[0] - 1
	regbot := term.height() - 1
	if len(params) > 1 {
		regbot = params[1] - 1
	}
	if regtop >= regbot {
		log.Warnf("DECSTBM: invalid region requested: %v, %v", regtop, regbot)
		// FIXME - try set something?
		return emptyZone
	}
	settings.regtop = min(regtop, term.height()-1)
	settings.regbot = min(regbot, term.height()-1)
	// log.Info("DECSTBM: set scrolling region to: ",
	// 	settings.regtop, settings.regbot)
	term.cursorTo(0, 0)
	return emptyZone
}

func testDECSTBM(t *testing.T) {
	testImpl(
		"DECSTBM",
		"With params",
		[]int{2, 3},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
			settings: TermSettings{
				regtop: 0,
				regbot: 1,
			},
		},
		testState{
			visualScreen: `^a b c
			                d e f
			                g h i`,
			settings: TermSettings{
				regtop: 1,
				regbot: 2,
			},
		},
		emptyZone,
		t,
	)
	testImpl(
		"DECSTBM",
		"Handle params if too high",
		[]int{1, 9},
		testState{
			visualScreen: `a b c
			               d^e f
			               g h i`,
			settings: TermSettings{
				regtop: 0,
				regbot: 1,
			},
		},
		testState{
			visualScreen: `^a b c
			                d e f
			                g h i`,
			settings: TermSettings{
				regtop: 0,
				regbot: 2,
			},
		},
		emptyZone,
		t,
	)
	testImpl(
		"DECSTBM",
		"With invalid params - region height must be > 1",
		[]int{1, 1},
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
}
