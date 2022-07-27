package vterm

/*
  DECRST
  DEC Private Mode Reset (DECRST).

  CSI ? Pm l
DEC Private Mode Reset (DECRST).
  Ps = 1  -> Normal Cursor Keys (DECCKM).
  Ps = 2  -> Designate VT52 mode (DECANM).
  Ps = 3  -> 80 Column Mode (DECCOLM).
  Ps = 4  -> Jump (Fast) Scroll (DECSCLM).
  Ps = 5  -> Normal Video (DECSCNM).
  Ps = 6  -> Normal Cursor Mode (DECOM).
  Ps = 7  -> No Wraparound Mode (DECAWM).
  Ps = 8  -> No Auto-repeat Keys (DECARM).
  Ps = 9  -> Don't send Mouse X & Y on button press.
  Ps = 1 0  -> Hide toolbar (rxvt).
  Ps = 1 2  -> Stop Blinking Cursor (AT&T 610).
  Ps = 1 3  -> Disable Blinking Cursor (reset only via
resource or menu).
  Ps = 1 4  -> Disable XOR of Blinking Cursor control sequence
and menu.
  Ps = 1 8  -> Don't print form feed (DECPFF).
  Ps = 1 9  -> Limit print to scrolling region (DECPEX).
  Ps = 2 5  -> Hide Cursor (DECTCEM).
  Ps = 3 0  -> Don't show scrollbar (rxvt).
  Ps = 3 5  -> Disable font-shifting functions (rxvt).
  Ps = 4 0  -> Disallow 80 -> 132 Mode.
  Ps = 4 1  -> No more(1) fix (see curses(3) resource).
  Ps = 4 2  -> Disable National Replacement Character sets
(DECNRCM).
  Ps = 4 4  -> Turn Off Margin Bell.
  Ps = 4 5  -> No Reverse-wraparound Mode.
  Ps = 4 6  -> Stop Logging.  (This is normally disabled by a
compile-time option).
  Ps = 4 7  -> Use Normal Screen Buffer.
  Ps = 6 6  -> Numeric keypad (DECNKM).
  Ps = 6 7  -> Backarrow key sends delete (DECBKM).
  Ps = 6 9  -> Disable left and right margin mode (DECLRMM),
VT420 and up.
  Ps = 9 5  -> Clear screen when DECCOLM is set/reset (DEC-
NCSM), VT510 and up.
  Ps = 1 0 0 0  -> Don't send Mouse X & Y on button press and
release.  See the section Mouse Tracking.
  Ps = 1 0 0 1  -> Don't use Hilite Mouse Tracking.
  Ps = 1 0 0 2  -> Don't use Cell Motion Mouse Tracking.
  Ps = 1 0 0 3  -> Don't use All Motion Mouse Tracking.
  Ps = 1 0 0 4  -> Don't send FocusIn/FocusOut events.
  Ps = 1 0 0 5  -> Disable UTF-8 Mouse Mode.
  Ps = 1 0 0 6  -> Disable SGR Mouse Mode.
  Ps = 1 0 0 7  -> Disable Alternate Scroll Mode, i.e., the
alternateScroll resource.
  Ps = 1 0 1 0  -> Don't scroll to bottom on tty output
(rxvt).
  Ps = 1 0 1 1  -> Don't scroll to bottom on key press (rxvt).
  Ps = 1 0 1 5  -> Disable urxvt Mouse Mode.
  Ps = 1 0 3 4  -> Don't interpret "meta" key.  (This disables
the eightBitInput resource).
  Ps = 1 0 3 5  -> Disable special modifiers for Alt and Num-
Lock keys.  (This disables the numLock resource).
  Ps = 1 0 3 6  -> Don't send ESC  when Meta modifies a key.
(This disables the metaSendsEscape resource).
  Ps = 1 0 3 7  -> Send VT220 Remove from the editing-keypad
Delete key.
  Ps = 1 0 3 9  -> Don't send ESC  when Alt modifies a key.
(This disables the altSendsEscape resource).
  Ps = 1 0 4 0  -> Do not keep selection when not highlighted.
(This disables the keepSelection resource).
  Ps = 1 0 4 1  -> Use the PRIMARY selection.  (This disables
the selectToClipboard resource).
  Ps = 1 0 4 2  -> Disable Urgency window manager hint when
Control-G is received.  (This disables the bellIsUrgent
resource).
  Ps = 1 0 4 3  -> Disable raising of the window when Control-
G is received.  (This disables the popOnBell resource).
  Ps = 1 0 4 6  -> Disable switching to/from Alternate Screen
Buffer.  This works for terminfo-based systems, updating the
titeInhibit resource.  If currently using the Alternate Screen
Buffer, xterm switches to the Normal Screen Buffer.
  Ps = 1 0 4 7  -> Use Normal Screen Buffer, clearing screen
first if in the Alternate Screen Buffer.  (This may be dis-
abled by the titeInhibit resource).
  Ps = 1 0 4 8  -> Restore cursor as in DECRC.  (This may be
disabled by the titeInhibit resource).
  Ps = 1 0 4 9  -> Use Normal Screen Buffer and restore cursor
as in DECRC.  (This may be disabled by the titeInhibit
resource).  This combines the effects of the 1 0 4 7  and 1 0
4 8  modes.  Use this with terminfo-based applications rather
than the 4 7  mode.
  Ps = 1 0 5 0  -> Reset terminfo/termcap function-key mode.
  Ps = 1 0 5 1  -> Reset Sun function-key mode.
  Ps = 1 0 5 2  -> Reset HP function-key mode.
  Ps = 1 0 5 3  -> Reset SCO function-key mode.
  Ps = 1 0 6 0  -> Reset legacy keyboard emulation (i.e,
X11R6).
  Ps = 1 0 6 1  -> Reset keyboard emulation to Sun/PC style.
  Ps = 2 0 0 4  -> Reset bracketed paste mode.
*/

// edits above this line will be lost if code is generated again

import (
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"testing"
)

func doDECRST(spec *vparser.ControlSpec, params []int, term *Term) Rect {
	decSetOrRst(spec, params, term, false)
	return emptyZone
}

func decSetOrRst(spec *vparser.ControlSpec, params []int, term *Term, isSet bool) {
	if len(params) != 1 {
		log.Warnln("DECSET: invalid number of params:", params)
		return
	}
	mode := params[0]
	switch mode {
	case 1: // Application Cursor Keys (DECCKM)
		term.settings.DECCKM = isSet
	default:
		name := "DECSET"
		if !isSet {
			name = "DECRST"
		}
		log.Warnf("%v mode not implemented", name)
	}
}

func testDECRST(t *testing.T) {
	testImpl(
		"DECRST",
		"Application Cursor Keys (DECCKM)",
		[]int{1},
		testState{
			settings: TermSettings{DECCKM: true},
		},
		testState{
			settings: TermSettings{DECCKM: false},
		},
		emptyZone,
		t,
	)
}
