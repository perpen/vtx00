package vscreen

import (
	"fmt"

	runewidth "github.com/mattn/go-runewidth"
	log "github.com/sirupsen/logrus"
)

type Line struct {
	cells   []Cell
	wrapped bool
}

type Screen struct {
	lines          []*Line
	top            int // first line of live screen
	scrolledTop    int // first line of scrolled screen
	w, h           int // size
	cx, cy         int // cursor
	scrollbackSize int // max number of cells used by scrollback buffer
	total          int // height + scrollback
}

// scrollbackSize is max number of cells to be used by scrollback buffer
func NewScreen(w, h, scrollbackSize int) (Screen, error) {
	log.Debugln("vterm.NewScreen:", w, h)
	if scrollbackSize < 0 {
		log.Warnf("invalid scrollback size: %v", scrollbackSize)
		scrollbackSize = 0
	}
	scr := Screen{
		scrollbackSize: scrollbackSize,
	}
	err := scr.Resize(w, h)
	if err != nil {
		return Screen{}, err
	}
	return scr, nil
}

func (scr *Screen) newLine(w int) *Line {
	return &Line{
		cells: make([]Cell, w),
	}
}

func (scr *Screen) Width() int {
	return scr.w
}

func (scr *Screen) Height() int {
	return scr.h
}

// var names ending with 2 relate to the new screen
func (scr *Screen) Resize(w2, h2 int) error {
	log.Debugln("Screen.Resize:", w2, h2)
	if w2 == scr.w && h2 == scr.h {
		return nil
	}
	if w2 < 1 || h2 < 1 {
		return fmt.Errorf("invalid screen size: %v, %v", w2, h2)
	}
	total2 := h2 + scr.scrollbackSize/w2
	cx2, cy2 := 0, 0
	if len(scr.lines) == 0 {
		// Not a resize, the screen is brand new
		scr.lines = scr.newLines(total2, w2)
	} else {
		log.Infof("Screen.Resize: scr.total=%v total=%v h=%v", scr.total, total2, h2)

		i2 := 0 // cell index, spanning all lines
		lines2 := scr.newLines(total2, w2)

		ul, ucx, ucy := scr.unwrappedLines()
		log.Infof("Screen.Resize: ul=%v ucx=%v ucy=%v", ul, ucx, ucy)
		linesCount2 := 0
		for uy := 0; uy < len(ul); uy++ {
			prevLineInd := 0
			for x := 0; x < len(ul[uy]); x++ {
				c := ul[uy][x]
				x2 := i2 % w2
				y2 := (i2 / w2) % total2

				if ucx == x && ucy == uy {
					cx2 = x2
					cy2 = y2
				}

				line2 := lines2[y2]
				// log.Infof("Screen.Resize: i2=%v ch=%q", i2, c.Ch)
				overflow := x2 + c.width - w2
				if overflow > 0 {
					// Does not fit
					if c.width >= scr.w {
						// Doesn't fit, even on full screen width
						c.unprintable = c.Ch
						c.Ch = '?'
					} else {
						// Fill rest of line with empty cells, we'll print on next line
						i2 += overflow
						y2 = (y2 + 1) % total2
					}
				}
				log.Infof("Screen.Resize: i2=%v ch=%q at %v, %v", i2, c.Ch, x2, y2)
				line2.cells[x2] = c
				line2.wrapped = true
				if c.width == 0 {
					// blank cell
					i2++
				} else {
					i2 += c.width
				}
				y2 = (i2 / w2) % total2
				if y2 > prevLineInd {
					linesCount2++
					prevLineInd = y2
				}
			}

			y2 := ((i2 - 1) / w2) % total2
			lines2[y2].wrapped = false
			ll := len(ul[uy])
			if ll < w2 {
				i2 += w2 - (ll % w2)
			}
		}

		var top2 int
		if linesCount2 <= h2 {
			top2 = 0
		} else {
			top2 = linesCount2 % h2
		}

		scr.lines = lines2
		scr.top = top2
	}
	scr.total = total2
	scr.w, scr.h = w2, h2
	scr.cx, scr.cy = cx2, cy2
	log.Infof("Screen.Resize done")
	return nil
}

func (scr *Screen) unwrappedLines() (cells [][]Cell, cx, cy int) {
	lines := make([][]Cell, 0)
	line := make([]Cell, 0)
	// var lineWidths []int
	var ccx, ccy int
	for abs := 0; abs < scr.total; abs++ {
		i := (scr.top + abs) % scr.total
		row := scr.lines[i]
		blanks := 0
		for x := 0; x < scr.w; x++ {
			c := row.cells[x]
			log.Infof("Screen.unwrappedLines: %q", c.Ch)
			if c.unprintable != 0 {
				c.Ch = c.unprintable
				c.unprintable = 0
			}
			if c.width == 0 {
				blanks++
			} else {
				// log.Infof("Screen.unwrappedLines: blanks in line=%v", blanks)
				for i := 0; i < blanks; i++ {
					line = append(line, Cell{})
				}
				line = append(line, c)
				blanks = 0
			}
			log.Infof("Screen.unwrappedLines: x, y=%v, %v blanks=%v",
				x, len(line)-1, blanks)
			if x+blanks == scr.cx && scr.cy == i {
			// if x == scr.cx && scr.cy == i {
				ccx = x
				ccy = i
				log.Infof("Screen.unwrappedLines: ccx, ccy=%v, %v ################",
					ccx, ccy)
			}
		}
		if row.wrapped {
			// Preserve trailing blanks on wrapped row
			log.Infof("Screen.unwrappedLines: blanks at end=%v", blanks)
			for i := 0; i < blanks; i++ {
				line = append(line, Cell{})
			}
		}
		if !row.wrapped || abs == scr.total-1 {
			// TODO remove trailing blanks, in case of multiple wrapped blank lines
			// for i := len(line) - 1; i >= 0 && line[i].width > 0; i-- {}
			var lastCharIndex int
			for lastCharIndex = len(line) - 1; lastCharIndex >= 0; lastCharIndex-- {
				// log.Infof("Screen.unwrappedLines: line[%v]=%q", lastCharIndex, line[lastCharIndex].Ch)
				if line[lastCharIndex].width > 0 {
					break
				}
			}
			// log.Infof("Screen.unwrappedLines: last char=%v", i)
			if lastCharIndex > 0 {
				// log.Infof("Screen.unwrappedLines: pre len=%v", len(line))
				line = line[0 : lastCharIndex+1]
				// log.Infof("Screen.unwrappedLines: post len=%v", len(line))
			}
			lines = append(lines, line)
			line = make([]Cell, 0)
		}
	}
	return lines, ccx, ccy
}

// FIXME could be inner function of Resize()
func (scr *Screen) newLines(total, w int) []*Line {
	lines := make([]*Line, total)
	for i, _ := range lines {
		lines[i] = scr.newLine(w)
	}
	return lines
}

func (scr *Screen) Scroll(lines int) {
	maxScroll := scr.total - scr.h
	if lines > maxScroll {
		log.Warnf("Screen.Scroll: cannot scroll by %v lines, max is %v",
			lines, maxScroll)
		return
	}
	scr.scrolledTop = (scr.top + scr.total - lines) % scr.total
}

func (scr *Screen) Cursor() (x, y int) {
	return scr.cx, scr.cy
}

func (scr *Screen) CursorTo(x, y int) {
	log.Debug("Screen.CursorTo:", x, y)
	scr.cx, scr.cy = x, y
}

func (scr *Screen) setWrapped(y int, b bool) {
	// log.Infof("Screen.setWrapped(%v, %v)", y, b)
	scr.line(y).wrapped = b
}

func (scr *Screen) isWrapped(y int) bool {
	return scr.line(y).wrapped
}

func (scr *Screen) line(y int) *Line {
	return scr.lines[scr.lineIndex(y)]
}

func (scr *Screen) lineIndex(y int) int {
	// log.Infof("Screen.lineIndex(%v)=%v top=%v", y, (scr.top + y) % scr.total, scr.top)
	return (scr.top + y) % scr.total
}

func (scr *Screen) validCoords(x, y int) bool {
	// log.Infof("Screen.validCoords: %v, %v size: %v, %v", x, y, scr.w, scr.h)
	return ordered(-1, x, scr.w) && ordered(-1, y, scr.h)
}

func (scr *Screen) CellAt(x, y int) *Cell {
	if !scr.validCoords(x, y) {
		return nil
	}
	return scr.cellAt(x, y)
}

func (scr *Screen) cellAt(x, y int) *Cell {
	// log.Infof("Screen.cellAt: %v, %v", x, y)
	return &scr.line(y).cells[x]
}

func (scr *Screen) setCellAt(x, y int, ch rune, width int, style Style) {
	c := scr.cellAt(x, y)
	c.set(ch, width, style)
}

func (scr *Screen) ClearLines(y, count int) {
	// log.Infof("Screen.ClearLines(%v, %v)", y, count)
	for i := 0; i < min(count, scr.h-y+0); i++ {
		scr.clearLine(y + i)
	}
}

func (scr *Screen) clearLine(y int) {
	scr.ClearOnLine(0, y, scr.w)
}

// Clear specified number of cells on line, starting from coords
func (scr *Screen) ClearOnLine(x, y, count int) {
	// log.Infof("Screen.ClearOnLine(%v, %v, %v)", x, y, count)
	line := scr.line(y)
	for i := 0; i < count; i++ {
		// log.Infof("Screen.clearFrom(): i=%v", i)
		line.cells[x+i].set(0, 0, defaultStyle)
	}
}

// Moves n lines from line y by delta number of lines (may be negative)
func (scr *Screen) MoveLines(y, n, delta int) {
	// log.Infof("Screen.MoveLines(%v, %v, %v)", y, n, delta)
	if n == 0 || delta == 0 {
		return
	}

	moveLine := func(i int) {
		// log.Infof("Screen.MoveLines(): loop i=%v", i)
		srcIndex := scr.lineIndex(i)
		// log.Infof("Screen.MoveLines(): i=%v i+delta=%v", i, i+delta)
		if ordered(-1, i+delta, scr.h) {
			// log.Infof("Screen.MoveLines(): i=%v i+delta=%v IN", i, i+delta)
			tgtIndex := scr.lineIndex(i + delta)
			savedTgt := scr.lines[tgtIndex]
			scr.lines[tgtIndex] = scr.lines[srcIndex]
			scr.lines[srcIndex] = savedTgt
		}
		if i < scr.h {
			// log.Infof("Screen.MoveLines(): clearing %v", i)
			scr.clearLine(i)
		}
	}

	if delta > 0 {
		for i := y + n - 1; i >= y; i-- {
			moveLine(i)
		}
	} else {
		for i := y; i < y+n; i++ {
			moveLine(i)
		}
	}
}

// Print adds the chars from the current cursor position, advances cursor.
func (scr *Screen) Print(s string, style Style) (x, y, w, h int) {
	// advance := func(n int) {
	// 	over := scr.cx + n - scr.w
	// 	if over > 1 {

	// 	} else {
	// 		scr.cx += n
	// 	}
	// }

	cx1, cy1 := scr.Cursor()
	rolled := false

	for _, ch := range s {
		// log.Infof("ch=%q", ch)
		if scr.cy == scr.h-1 && scr.cx == scr.w {
			scr.roll()
			rolled = true
		}

		chWidth := runewidth.RuneWidth(ch)
		if chWidth > scr.w || chWidth == 0 {
			ch = '?'
			chWidth = 1
		}
		log.Infof("Screen.print: ch: %q at: %v, %v", ch, scr.cx, scr.cy)

		// If char doesn't fit at position then move to next line before printing
		overflow := scr.cx + chWidth - scr.w
		// log.Infof("overflow=%v", overflow)
		if overflow > 0 {
			blanks := scr.w - scr.cx
			// log.Infof("blanks=%v, top=%v", blanks, scr.top)
			scr.ClearOnLine(scr.cx, scr.cy, blanks)
			scr.setWrapped(scr.cy, true)
			scr.cx = 0
			scr.cy++
			scr.setWrapped(scr.cy, false)
		}

		// Put the char
		// log.Infof("Screen.print2: ch: %q at: %v, %v", ch, scr.cx, scr.cy)
		scr.setCellAt(scr.cx, scr.cy, ch, chWidth, style)
		scr.cx++

		// If char is wide, mark following cells as shadows
		if chWidth > 1 {
			shadowLen := chWidth - 1
			// log.Infof("shadowLen=%v", shadowLen)
			scr.ClearOnLine(scr.cx, scr.cy, shadowLen)
			scr.cx += shadowLen
		}
	}

	if rolled {
		return 0, 0, scr.w, scr.h
	} else {
		if scr.cy == cy1 {
			return cx1, cy1, scr.cx - cx1, 1
		} else {
			return 0, cy1, scr.w, scr.cy - cx1 + 1
		}
	}
}

func (scr *Screen) roll() {
	// log.Infoln("Screen.roll")
	scr.cy--
	scr.top = (scr.top + 1) % scr.total
	scr.ClearOnLine(0, scr.h-1, scr.w)
	scr.setWrapped(scr.h-1, false)
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ordered(a, b, c int) bool {
	return a < b && b < c
}
