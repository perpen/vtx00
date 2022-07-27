package vterm

/*
  SGR
  Character Attributes

  CSI Pm m  Character Attributes (SGR).
  Ps = 0  -> Normal (default).
  Ps = 1  -> Bold.
  Ps = 2  -> Faint, decreased intensity (ISO 6429).
  Ps = 3  -> Italicized (ISO 6429).
  Ps = 4  -> Underlined.
  Ps = 5  -> Blink (appears as Bold in X11R6 xterm).
  Ps = 7  -> Inverse.
  Ps = 8  -> Invisible, i.e., hidden (VT300).
  Ps = 9  -> Crossed-out characters (ISO 6429).
  Ps = 2 1  -> Doubly-underlined (ISO 6429).
  Ps = 2 2  -> Normal (neither bold nor faint).
  Ps = 2 3  -> Not italicized (ISO 6429).
  Ps = 2 4  -> Not underlined.
  Ps = 2 5  -> Steady (not blinking).
  Ps = 2 7  -> Positive (not inverse).
  Ps = 2 8  -> Visible, i.e., not hidden (VT300).
  Ps = 2 9  -> Not crossed-out (ISO 6429).
  Ps = 3 0  -> Set foreground color to Black.
  Ps = 3 1  -> Set foreground color to Red.
  Ps = 3 2  -> Set foreground color to Green.
  Ps = 3 3  -> Set foreground color to Yellow.
  Ps = 3 4  -> Set foreground color to Blue.
  Ps = 3 5  -> Set foreground color to Magenta.
  Ps = 3 6  -> Set foreground color to Cyan.
  Ps = 3 7  -> Set foreground color to White.
  Ps = 3 9  -> Set foreground color to default (original).
  Ps = 4 0  -> Set background color to Black.
  Ps = 4 1  -> Set background color to Red.
  Ps = 4 2  -> Set background color to Green.
  Ps = 4 3  -> Set background color to Yellow.
  Ps = 4 4  -> Set background color to Blue.
  Ps = 4 5  -> Set background color to Magenta.
  Ps = 4 6  -> Set background color to Cyan.
  Ps = 4 7  -> Set background color to White.
  Ps = 4 9  -> Set background color to default (original).

If 16-color support is compiled, the following apply.  Assume
that xterm's resources are set so that the ISO color codes are
the first 8 of a set of 16.  Then the aixterm colors are the
bright versions of the ISO colors:
  Ps = 9 0  -> Set foreground color to Black.
  Ps = 9 1  -> Set foreground color to Red.
  Ps = 9 2  -> Set foreground color to Green.
  Ps = 9 3  -> Set foreground color to Yellow.
  Ps = 9 4  -> Set foreground color to Blue.
  Ps = 9 5  -> Set foreground color to Magenta.
  Ps = 9 6  -> Set foreground color to Cyan.
  Ps = 9 7  -> Set foreground color to White.
  Ps = 1 0 0  -> Set background color to Black.
  Ps = 1 0 1  -> Set background color to Red.
  Ps = 1 0 2  -> Set background color to Green.
  Ps = 1 0 3  -> Set background color to Yellow.
  Ps = 1 0 4  -> Set background color to Blue.
  Ps = 1 0 5  -> Set background color to Magenta.
  Ps = 1 0 6  -> Set background color to Cyan.
  Ps = 1 0 7  -> Set background color to White.

If xterm is compiled with the 16-color support disabled, it
supports the following, from rxvt:
  Ps = 1 0 0  -> Set foreground and background color to
default.

Xterm maintains a color palette whose entries are identified
by an index beginning with zero.  If 88- or 256-color support
is compiled, the following apply:
o   All parameters are decimal integers.
o   RGB values range from zero (0) to 255.
o   ISO-8613-6 has been interpreted in more than one way;
    xterm allows the semicolons separating the subparameters
    in this control to be replaced by colons (but after the
    first colon, colons must be used).

These ISO-8613-6 controls are supported:
  Pm = 3 8 ; 2 ; Pi; Pr; Pg; Pb -> Set foreground color to the
closest match in xterm's palette for the given RGB Pr/Pg/Pb.
The color space identifier Pi is ignored.
  Pm = 3 8 ; 5 ; Ps -> Set foreground color to Ps.
  Pm = 4 8 ; 2 ; Pi; Pr; Pg; Pb -> Set background color to the
closest match in xterm's palette for the given RGB Pr/Pg/Pb.
The color space identifier Pi is ignored.
  Pm = 4 8 ; 5 ; Ps -> Set background color to Ps.

This variation on ISO-8613-6 is supported for compatibility
with KDE konsole:
  Pm = 3 8 ; 2 ; Pr; Pg; Pb -> Set foreground color to the
closest match in xterm's palette for the given RGB Pr/Pg/Pb.
  Pm = 4 8 ; 2 ; Pr; Pg; Pb -> Set background color to the
closest match in xterm's palette for the given RGB Pr/Pg/Pb.

If xterm is compiled with direct-color support, and the
resource directColor is true, then rather than choosing the
closest match, xterm asks the X server to directly render a
given color.
*/

// edits above this line will be lost if code is generated again
import (
	"fmt"
	"github.com/perpen/vtx00/vparser"
	"github.com/perpen/vtx00/vscreen"
	log "github.com/sirupsen/logrus"
	"testing"
)

var colors8 = []string{
	"black",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",
}

func doSGR(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	for _, mode := range params {
		implemented := true
		switch {
		case mode == 0: // Normal (default).
			term.settings.style.Bold = false
			term.settings.style.Reverse = false
			term.settings.style.Italics = false
			term.settings.style.Underlined = false
			term.settings.style.Fg = "white" // FIXME, should be default
			term.settings.style.Bg = "black" // FIXME, should be default
		case mode == 1: // Bold.
			term.settings.style.Bold = true
		case mode == 4: // Underlined.
			term.settings.style.Underlined = true
		case mode == 7: // Reverse.
			term.settings.style.Reverse = true
		case 30 <= mode && mode < 38: // Set foreground color
			term.settings.style.Fg = lookupColor8(mode - 30)
		case mode == 39: // Set foreground color to default
			term.settings.style.Fg = "white" // FIXME
		case 40 <= mode && mode < 48: // Set background color
			term.settings.style.Bg = lookupColor8(mode - 40)
		case mode == 49: // Set background color
			term.settings.style.Bg = "black" // FIXME, should be default
		case 90 <= mode && mode < 98: // Set foreground color
			term.settings.style.Fg = lookupColor8(mode - 90)
		case 100 <= mode && mode < 108: // Set background color
			term.settings.style.Bg = lookupColor8(mode - 100)
		default:
			implemented = false
		}
		if implemented {
			log.Debugf("SGR mode %v -> %v", mode, spec.Selections[mode])
		} else {
			log.Warn("SGR mode not implemented")
		}
	}
	return emptyZone
}

func lookupColor8(i int) string {
	return colors8[i]
}

func testSGR(t *testing.T) {
	testImpl(
		"SGR",
		"0 Normal (default)",
		[]int{},
		testState{
			settings: TermSettings{
				style: vscreen.Style{
					Bold:       true,
					Reverse:    true,
					Italics:    true,
					Underlined: true,
					Fg:         "red",
					Bg:         "yellow",
				},
			},
		},
		testState{},
		emptyZone,
		t,
	)
	testImpl(
		"SGR",
		"1 Bold",
		[]int{1},
		testState{},
		testState{
			settings: TermSettings{
				style: vscreen.Style{
					Bold: true,
				},
			},
		},
		emptyZone,
		t,
	)
	testImpl(
		"SGR",
		"4 Underlined",
		[]int{4},
		testState{},
		testState{
			settings: TermSettings{
				style: vscreen.Style{
					Underlined: true,
				},
			},
		},
		emptyZone,
		t,
	)
	testImpl(
		"SGR",
		"7 Reverse",
		[]int{7},
		testState{},
		testState{
			settings: TermSettings{
				style: vscreen.Style{
					Reverse: true,
				},
			},
		},
		emptyZone,
		t,
	)
	for i, color := range colors8 {
		for _, n := range []int{30, 90} {
			testImpl(
				"SGR",
				fmt.Sprintf("%v %v foreground", n+i, color),
				[]int{n + i},
				testState{},
				testState{
					settings: TermSettings{
						style: vscreen.Style{
							Fg: color,
						},
					},
				},
				emptyZone,
				t,
			)
			testImpl(
				"SGR",
				fmt.Sprintf("%v %v background", n+10+i, color),
				[]int{n + 10 + i},
				testState{},
				testState{
					settings: TermSettings{
						style: vscreen.Style{
							Bg: color,
						},
					},
				},
				emptyZone,
				t,
			)
		}
	}
}
