package vterm

import (
	"fmt"
	"io"

	"github.com/perpen/vtx00/vparser"
	"github.com/perpen/vtx00/vscreen"
	log "github.com/sirupsen/logrus"
)

// FIXME use it
type grid interface {
	Set(x, y, n int, ch rune, style vscreen.Style) error
	Clear(x, y, n int) error
	Shift(x, y, n, delta int) error
	Resize(w, h int) error
	Print(s string, style vscreen.Style) error
}

type TermSettings struct {
	title          string
	style          vscreen.Style
	savedCursor    vscreen.Pair
	regtop, regbot int
	DECCKM         bool
}

var defaultSettings = TermSettings{
	title: "TITLE",
	style: vscreen.Style{
		Bold:       false,
		Reverse:    false,
		Italics:    false,
		Underlined: false,
		Fg:         "white",
		Bg:         "black",
	},
	savedCursor: vpair(0, 0),
	regtop:      -1,
	regbot:      -1,
	DECCKM:      false,
}

type termState struct {
	settings TermSettings
	Screen   vscreen.Screen
}

type Term struct {
	termState
	fineDamage bool
	pty        io.Writer
	parser     *vparser.Parser
	parserRcvr parserRcvr
	damageChan chan Damage
}

func NewTerm(pty io.Writer, damageChan chan Damage) *Term {
	t := new(Term)
	t.pty = pty
	t.parserRcvr = parserRcvr{t}
	bindings := vparser.NewBindings(AllControls)
	// FIXME reuse the same bindings from everywhere, ie create a
	// specialised newParser() function.
	t.parser = vparser.NewParser(bindings, t.parserRcvr)
	t.damageChan = damageChan

	t.settings.savedCursor = vpair(-1, -1)
	t.settings.regtop, t.settings.regbot = -1, -1

	return t
}

func (t *Term) pushDamage(zone Rect) {
	dmg := Damage{
		Term: t,
		X:    zone.x,
		Y:    zone.y,
		W:    zone.w,
		H:    zone.h,
	}
	log.Infof("vterm.pushDamage: zone=%v", zone)
	t.damageChan <- dmg
}

func (t *Term) Write(p []byte) (n int, err error) {
	// log.Debugln("Term.Write:", p)
	t.parser.Parse(p)
	if !t.fineDamage {
		t.pushDamage(Rect{0, 0, t.width(), t.height()})
	}
	return len(p), nil
}

func (t *Term) newDamage(x, y, w, h int, cells []*vscreen.Cell) Damage {
	dmg := Damage{t, x, y, w, h}
	return dmg
}

func (t *Term) width() int {
	return t.Screen.Width()
}

func (t *Term) height() int {
	return t.Screen.Height()
}

// Refresh .
func (t *Term) Refresh() {
	t.damageChan <- t.tmpFullDamage()
}

func (t *Term) tmpFullDamage() Damage {
	w, h := t.width(), t.height()

	dmgCells := make([]*vscreen.Cell, 0, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			cell := t.Screen.CellAt(x, y)
			if cell == nil {
				log.Fatalf("tmpFullDamage invalid coords: %v, %v", x, y)
			}
			dmgCells = append(dmgCells, cell)
		}
	}

	//log.Debugln("Term.tmpFullDamage: dmgCells[0]:", dmgCells[0])

	return t.newDamage(0, 0, w, h, dmgCells)
}

// Resize .
func (t *Term) Resize(w, h int) error {
	log.Debugln("vterm.Term.Resize:", w, h)
	if t.Screen.Width() == 0 {
		screen, err := vscreen.NewScreen(w, h, 100)
		if err != nil {
			return fmt.Errorf("Term.Resize: %v", err)
		}
		t.Screen = screen
	} else {
		t.Screen.Resize(w, h)
	}
	// FIXME - try reuse existing region, although good programs should
	// handle the winch and adjust it themselves
	t.settings.regtop = 0
	t.settings.regbot = t.Screen.Height() - 1
	return nil
}

func (t *Term) cursor() (x, y int) {
	return t.Screen.Cursor()
}

func (t *Term) cursorTo(x, y int) {
	t.Screen.CursorTo(x, y)
}

func (t *Term) print(str string, style vscreen.Style) (x, y, w, h int) {
	return t.Screen.Print(str, style)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func vpair(x, y int) vscreen.Pair {
	return vscreen.Pair{X: x, Y: y}
}
