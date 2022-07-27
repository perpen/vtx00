package vterm

/*
  DECKPNM
  Normal Keypad

  ESC >     Normal Keypad (DECKPNM).
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"testing"
)

func doDECKPNM(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	log.Warn("not implemented")
	return emptyZone
}

func testDECKPNM(t *testing.T) {
}
