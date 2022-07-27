package vterm

/*
  XDESG1
  Designate G1 Character Set (ISO 2022, VT100).

  ESC ) C   Designate G1 Character Set (ISO 2022, VT100). (HFD15)
The same character sets apply as for ESC ( C.
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"testing"
)

func doXDESG1(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	log.Infof("XDESG1 %q not implemented", params[0])
	return emptyZone
}

func testXDESG1(t *testing.T) {
	/*
		testImpl(
			"XDESG1",
			"XDESG1 test",
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
