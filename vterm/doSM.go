package vterm

/*
  SM
  Set Mode

  CSI Pm h  Set Mode (SM).
  Ps = 2  -> Keyboard Action Mode (AM).
  Ps = 4  -> Insert Mode (IRM).
  Ps = 1 2  -> Send/receive (SRM).
  Ps = 2 0  -> Automatic Newline (LNM).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"testing"
)

func doSM(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	log.Warn("SM not implemented")
	return emptyZone
}

func testSM(t *testing.T) {
	/*
		testImpl(
			"SM",
			"SM test",
			[]int{},
			testState{
				visualScreen: `a^b c
				               d e f
				               g h i`,
				settings: defaultSettings,
			},
			testState{
				visualScreen: `a^b c
				               d e f
				               g h i`,
				settings: defaultSettings,
			},
			t,
		)
	*/
}
