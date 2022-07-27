package vterm

/*
  RM
  Reset Mode

  CSI Pm l  Reset Mode (RM).
  Ps = 2  -> Keyboard Action Mode (AM).
  Ps = 4  -> Replace Mode (IRM).
  Ps = 1 2  -> Send/receive (SRM).
  Ps = 2 0  -> Normal Linefeed (LNM).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"testing"
)

func doRM(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	log.Warn("RM not implemented")
	return emptyZone
}

func testRM(t *testing.T) {
	/*
		testImpl(
			"RM",
			"RM test",
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
