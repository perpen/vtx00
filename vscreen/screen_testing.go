package vscreen

import (
	"bytes"
	"container/list"
	"fmt"
	"regexp"
	"strings"

	runewidth "github.com/mattn/go-runewidth"
	log "github.com/sirupsen/logrus"
)

// Returns a string showing the differences between the screens
// Empty string if they are equal
func ScreenDiff(screen1, screen2 Screen) string {
	spacing := 3
	columnWidths := make([]int, screen2.w)

	// Compare screens
	sameScreen, sameCursor := true, true
	for y := 0; y < screen2.h; y++ {
		for x := 0; x < screen2.w; x++ {
			cell1 := screen1.CellAt(x, y)
			if cell1 == nil {
				log.Fatalf("ScreenDiff: illegal coords: %v, %v", x, y)
			}
			cell2 := screen2.CellAt(x, y)
			if cell2 == nil {
				log.Fatalf("ScreenDiff: illegal coords: %v, %v", x, y)
			}
			colWidth := columnWidths[x]
			colWidth = max(colWidth, runewidth.RuneWidth(cell1.Ch))
			colWidth = max(colWidth, runewidth.RuneWidth(cell2.Ch))

			if cell1.Ch != cell2.Ch {
				// log.Infof("ScreenDiff: diff ch")
				sameScreen = false
			}
		}
		if screen1.line(y).wrapped != screen2.line(y).wrapped {
			// log.Infof("ScreenDiff: diff wrapped %v", y)
			sameScreen = false
		}
	}
	if screen1.cx != screen2.cx || screen1.cy != screen2.cy {
		sameCursor = false
	}

	var buf bytes.Buffer
	if !sameScreen {
		buf.WriteString("screens (expected, actual):\n")
		for y := 0; y < screen2.h; y++ {
			buf.WriteString("    ")
			for screen_num, screen := range []Screen{screen1, screen2} {
				for x := 0; x < screen2.w; x++ {
					if screen.cx == x {
						var prefix rune
						if screen.cy == y {
							prefix = '^'
						} else {
							prefix = ' '
						}
						buf.WriteRune(prefix)
					}
					cell := screen.CellAt(x, y)
					ch := cell.Ch
					switch ch {
					case 0:
						ch = '.'
					case ' ':
						ch = '_'
					}
					n, _ := buf.WriteRune(ch)
					for i := 0; i < columnWidths[x]-n; i++ {
						buf.WriteString(" ")
					}
				}
				if screen.w == screen.cx && y == screen.cy {
					buf.WriteString("^")
				} else {
					buf.WriteString(" ")
				}
				if screen.lines[y].wrapped {
					buf.WriteString("|")
				} else {
					buf.WriteString(" ")
				}
				if screen_num == 0 {
					for i := 0; i < spacing; i++ {
						buf.WriteString(" ")
					}
				}
			}
			buf.WriteString("\n")
		}
	}
	if !(sameCursor) {
		s := fmt.Sprintf("cursor: expected %v, %v, was %v, %v",
			screen1.cx, screen1.cy, screen2.cx, screen2.cy)
		buf.WriteString(s)
		buf.WriteString("\n")
	}

	return buf.String()
}

// Creates a screen from the string representation of it used in the tests.
// Each cell is separated by any number of spaces.
// A null char is represented by '.', and a space by '_'
// The position of the cursor is indicated by following a rune with '^'
func MakeScreenFromVisual(str string) (Screen, error) {
	if len(str) == 0 {
		return Screen{}, nil
	}
	oops := func(msg string) (Screen, error) {
		trimmed := regexp.MustCompile(`(?m)(^\s+|\s+$)`).ReplaceAllString(str, "")
		return Screen{}, fmt.Errorf("invalid screen: %v\n%v", msg, trimmed)
	}

	lines := regexp.MustCompile(`\s*\n\s*`).Split(strings.TrimSpace(str), -1)
	height := len(lines)
	if height == 0 {
		return oops("empty screen")
	}

	type runeWithCoord struct {
		r    rune
		x, y int
	}

	// Will be populated by the next for loop, then used to print
	// into the screen.
	runesWithCoords := list.New()
	wrappedLines := make(map[int]bool)

	width := 0
	cx, cy := -1, -1
	lastRuneWidth := -1
	for y, line := range lines {
		x := 0
		wrappedLines[y] = false
		for _, r := range line {
			switch r {
			case '^':
				if cx != -1 {
					return oops("more than 1 cursor")
				}
				cx, cy = x, y
			case ' ', '\t':
				// skip
			case '|':
				wrappedLines[y] = true
			default:
				runesWithCoords.PushBack(runeWithCoord{r, x, y})
				lastRuneWidth = runewidth.RuneWidth(r)
				x += lastRuneWidth
			}
		}
		if width > 0 && x != width {
			return oops("inconsistent lines lengths")
		}
		width = x
	}
	if width == 0 {
		return oops("empty rows")
	}
	if cx == -1 {
		return oops("missing cursor")
	}

	screen, err := NewScreen(width, height, 0)
	if err != nil {
		log.Fatal(err)
	}
	for elt := runesWithCoords.Front(); elt != nil; elt = elt.Next() {
		posRune := elt.Value.(runeWithCoord)
		switch posRune.r {
		case '.':
			continue
		case '_':
			posRune.r = ' '
		}
		screen.CursorTo(posRune.x, posRune.y)
		screen.Print(string(posRune.r), Style{})
	}
	for y, wrap := range wrappedLines {
		screen.setWrapped(y, wrap)
	}

	screen.CursorTo(cx, cy)
	return screen, nil
}
