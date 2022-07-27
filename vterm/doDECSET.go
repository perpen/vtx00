package vterm

/*
  DECSET
  DEC Private Mode Set

  CSI ? Pm h
DEC Private Mode Set (DECSET).
  Ps = 1  -> Application Cursor Keys (DECCKM).
  Ps = 2  -> Designate USASCII for character sets G0-G3
(DECANM), and set VT100 mode.
  Ps = 3  -> 132 Column Mode (DECCOLM).
  Ps = 4  -> Smooth (Slow) Scroll (DECSCLM).
  Ps = 5  -> Reverse Video (DECSCNM).
  Ps = 6  -> Origin Mode (DECOM).
  Ps = 7  -> Wraparound Mode (DECAWM).
  Ps = 8  -> Auto-repeat Keys (DECARM).
  Ps = 9  -> Send Mouse X & Y on button press.  See the sec-
tion Mouse Tracking.  This is the X10 xterm mouse protocol.
  Ps = 1 0  -> Show toolbar (rxvt).
  Ps = 1 2  -> Start Blinking Cursor (AT&T 610).
  Ps = 1 3  -> Start Blinking Cursor (set only via resource or
menu).
  Ps = 1 4  -> Enable XOR of Blinking Cursor control sequence
and menu.
  Ps = 1 8  -> Print form feed (DECPFF).
  Ps = 1 9  -> Set print extent to full screen (DECPEX).
  Ps = 2 5  -> Show Cursor (DECTCEM).
  Ps = 3 0  -> Show scrollbar (rxvt).
  Ps = 3 5  -> Enable font-shifting functions (rxvt).
  Ps = 3 8  -> Enter Tektronix Mode (DECTEK).
  Ps = 4 0  -> Allow 80 -> 132 Mode.
  Ps = 4 1  -> more(1) fix (see curses(3) resource).
  Ps = 4 2  -> Enable National Replacement Character sets
(DECNRCM).
  Ps = 4 4  -> Turn On Margin Bell.
  Ps = 4 5  -> Reverse-wraparound Mode.
  Ps = 4 6  -> Start Logging.  This is normally disabled by a
compile-time option.
  Ps = 4 7  -> Use Alternate Screen Buffer.  (This may be dis-
abled by the titeInhibit resource).
  Ps = 6 6  -> Application keypad (DECNKM).
  Ps = 6 7  -> Backarrow key sends backspace (DECBKM).
  Ps = 6 9  -> Enable left and right margin mode (DECLRMM),
VT420 and up.
  Ps = 9 5  -> Do not clear screen when DECCOLM is set/reset
(DECNCSM), VT510 and up.
  Ps = 1 0 0 0  -> Send Mouse X & Y on button press and
release.  See the section Mouse Tracking.  This is the X11
xterm mouse protocol.
  Ps = 1 0 0 1  -> Use Hilite Mouse Tracking.
  Ps = 1 0 0 2  -> Use Cell Motion Mouse Tracking.
  Ps = 1 0 0 3  -> Use All Motion Mouse Tracking.
  Ps = 1 0 0 4  -> Send FocusIn/FocusOut events.
  Ps = 1 0 0 5  -> Enable UTF-8 Mouse Mode.
  Ps = 1 0 0 6  -> Enable SGR Mouse Mode.
  Ps = 1 0 0 7  -> Enable Alternate Scroll Mode, i.e., the
alternateScroll resource.
  Ps = 1 0 1 0  -> Scroll to bottom on tty output (rxvt).
  Ps = 1 0 1 1  -> Scroll to bottom on key press (rxvt).
  Ps = 1 0 1 5  -> Enable urxvt Mouse Mode.
  Ps = 1 0 3 4  -> Interpret "meta" key, sets eighth bit.
(enables the eightBitInput resource).
  Ps = 1 0 3 5  -> Enable special modifiers for Alt and Num-
Lock keys.  (This enables the numLock resource).
  Ps = 1 0 3 6  -> Send ESC   when Meta modifies a key.  (This
enables the metaSendsEscape resource).
  Ps = 1 0 3 7  -> Send DEL from the editing-keypad Delete
key.
  Ps = 1 0 3 9  -> Send ESC  when Alt modifies a key.  (This
enables the altSendsEscape resource).
  Ps = 1 0 4 0  -> Keep selection even if not highlighted.
(This enables the keepSelection resource).
  Ps = 1 0 4 1  -> Use the CLIPBOARD selection.  (This enables
the selectToClipboard resource).
  Ps = 1 0 4 2  -> Enable Urgency window manager hint when
Control-G is received.  (This enables the bellIsUrgent
resource).
  Ps = 1 0 4 3  -> Enable raising of the window when Control-G
is received.  (enables the popOnBell resource).
  Ps = 1 0 4 4  -> Reuse the most recent data copied to CLIP-
BOARD.  (This enables the keepClipboard resource).
  Ps = 1 0 4 6  -> Enable switching to/from Alternate Screen
Buffer.  This works for terminfo-based systems, updating the
titeInhibit resource.
  Ps = 1 0 4 7  -> Use Alternate Screen Buffer.  (This may be
disabled by the titeInhibit resource).
  Ps = 1 0 4 8  -> Save cursor as in DECSC.  (This may be dis-
abled by the titeInhibit resource).
  Ps = 1 0 4 9  -> Save cursor as in DECSC and use Alternate
Screen Buffer, clearing it first.  (This may be disabled by
the titeInhibit resource).  This combines the effects of the 1
0 4 7  and 1 0 4 8  modes.  Use this with terminfo-based
applications rather than the 4 7  mode.
  Ps = 1 0 5 0  -> Set terminfo/termcap function-key mode.
  Ps = 1 0 5 1  -> Set Sun function-key mode.
  Ps = 1 0 5 2  -> Set HP function-key mode.
  Ps = 1 0 5 3  -> Set SCO function-key mode.
  Ps = 1 0 6 0  -> Set legacy keyboard emulation (i.e, X11R6).
  Ps = 1 0 6 1  -> Set VT220 keyboard emulation.
  Ps = 2 0 0 4  -> Set bracketed paste mode.
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	"testing"
)

func doDECSET(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	decSetOrRst(spec, params, term, true)
	return emptyZone
}

func testDECSET(t *testing.T) {
	testImpl(
		"DECSET",
		"Application Cursor Keys (DECCKM)",
		[]int{1},
		testState{
			settings: TermSettings{DECCKM: false},
		},
		testState{
			settings: TermSettings{DECCKM: true},
		},
		emptyZone,
		t,
	)
}
