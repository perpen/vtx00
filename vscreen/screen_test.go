package vscreen

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

// func TestScreenPrint(t *testing.T) {
// 	screenPrintTest := func(testName, s,
// 		visualInitial, visualExpected string) {

// 		// log.Infof("#### testing Print %v", testName)
// 		scr := makeScreenFromVisualForTest(visualInitial, testName)
// 		expected := makeScreenFromVisualForTest(visualExpected, testName)

// 		// log.Print("screenPrintTest: expected cursor: ", expected.Cursor())
// 		scr.Print(s, Style{})
// 		// log.Print("screenPrintTest: actual   cursor: ", scr.Cursor())

// 		diff := ScreenDiff(expected, scr)
// 		if len(diff) > 0 {
// 			t.Errorf("failure:\n%s\n%s", testName, diff)
// 		}
// 	}

// 	screenPrintTest(
// 		"print char on tiny screen",
// 		"a",
// 		`^. `,
// 		` a^`,
// 	)
// 	screenPrintTest(
// 		"print at end of tiny screen",
// 		"a",
// 		` .^`,
// 		` a^`,
// 	)
// 	screenPrintTest(
// 		"print wide char on tiny screen",
// 		"👦",
// 		`^. `,
// 		` ?^`,
// 	)
// 	screenPrintTest(
// 		"print wide char on narrow screen",
// 		"👦a",
// 		`^.
// 		  .
// 		  .`,
// 		` ?|
// 		  a^
// 		  .`,
// 	)
// 	screenPrintTest(
// 		"print char",
// 		"a",
// 		`^. .`,
// 		` a^.`,
// 	)
// 	screenPrintTest(
// 		"print long string",
// 		"xyz",
// 		`^a b
// 		  . .`,
// 		` x y|
// 		  z^.`,
// 	)
// 	screenPrintTest(
// 		"print long string with roll",
// 		"pqrst",
// 		`^a b
// 		  c d`,
// 		` r s|
// 		  t^.`,
// 	)
// 	screenPrintTest(
// 		"print wide char",
// 		"👦",
// 		`^a  b  c`,
// 		` 👦   ^c`,
// 	)
// 	screenPrintTest(
// 		"print at eol, cursor past end",
// 		"a",
// 		`. .^. `,
// 		`. . a^`,
// 	)
// 	screenPrintTest(
// 		"print wide char near eol",
// 		"👦",
// 		` a ^b   c
// 		  .  .   .`,
// 		` a  👦   ^
// 		  .  .  .`,
// 	)
// 	screenPrintTest(
// 		"print char wider than screen",
// 		"👦",
// 		`^a
// 		  b
// 		  c`,
// 		` ?^
// 		  b
// 		  c`,
// 	)
// 	screenPrintTest(
// 		"print char at the end",
// 		"x",
// 		`a
// 		 b
// 		 c^`,
// 		`b
// 		 c|
// 		 x^`,
// 	)
// 	screenPrintTest(
// 		"print char wider than screen at the end",
// 		"👦",
// 		` a
// 		  b
// 		  c^`,
// 		` b
// 		  c|
// 		  ?^`,
// 	)
// }

// func TestScreenShift(t *testing.T) {
// 	screenShiftTest := func(testName string,
// 		y, n, delta int,
// 		visualInitial, visualExpected string) {

// 		// log.Infof("testing Shift - %v", testName)
// 		scr := makeScreenFromVisualForTest(visualInitial, testName)
// 		expected := makeScreenFromVisualForTest(visualExpected, testName)

// 		scr.MoveLines(y, n, delta)

// 		diff := ScreenDiff(expected, scr)
// 		if len(diff) > 0 {
// 			t.Errorf("failure:\n%s\n%s", testName, diff)
// 		}
// 	}

// 	screenShiftTest(
// 		"shift by 0",
// 		1, 3, 0,
// 		`^a b c
// 		  d e f
// 		  g h i`,
// 		`^a b c
// 		  d e f
// 		  g h i`,
// 	)
// 	screenShiftTest(
// 		"shift 1 by 1",
// 		1, 1, 1,
// 		`^a b c
// 		  d e f
// 		  g h i`,
// 		`^a b c
// 		  . . .
// 		  d e f`,
// 	)
// 	screenShiftTest(
// 		"shift 2 by 2",
// 		0, 2, 2,
// 		`^a b c
// 		  d e f
// 		  g h i`,
// 		`^. . .
// 		  . . .
// 		  a b c`,
// 	)
// 	screenShiftTest(
// 		"shift 2 by 1",
// 		0, 2, 1,
// 		`^a b c
// 		  d e f
// 		  g h i`,
// 		`^. . .
// 		  a b c
// 		  d e f`,
// 	)
// 	screenShiftTest(
// 		"shift 2 by -1",
// 		1, 2, -1,
// 		`^a b c
// 		  d e f
// 		  g h i`,
// 		`^d e f
// 		  g h i
// 		  . . .`,
// 	)
// 	screenShiftTest(
// 		"shift 1 by -2",
// 		2, 1, -2,
// 		`^a b c
// 		  d e f
// 		  g h i`,
// 		`^g h i
// 		  d e f
// 		  . . .`,
// 	)
// 	screenShiftTest(
// 		"shift 2 by -2",
// 		2, 2, -2,
// 		`^a b c
// 		  d e f
// 		  g h i
// 		  j k l`,
// 		`^g h i
// 		  j k l
// 		  . . .
// 		  . . .`,
// 	)
// 	// FIXME remove unnecessary tests
// 	screenShiftTest(
// 		"shift by -3",
// 		0, 3, -3,
// 		`^a b c
// 		  d e f
// 		  g h i
// 		  j k l`,
// 		`^. . .
// 		  . . .
// 		  . . .
// 		  j k l`,
// 	)
// 	screenShiftTest(
// 		"shift 2 by high positive",
// 		1, 2, 4,
// 		`^a b
// 		  c d
// 		  e f
// 		  g h`,
// 		`^a b
// 		  . .
// 		  . .
// 		  g h`,
// 	)
// 	screenShiftTest(
// 		"shift too many from bottom",
// 		1, 2, 1,
// 		`^a b
// 		  c d`,
// 		`^a b
// 		  . .`,
// 	)
// 	screenShiftTest(
// 		"shift by high negative",
// 		1, 2, -4,
// 		`^a b
//    		  c d
//    		  e f
//    		  g h`,
// 		`^a b
//    		  . .
//    		  . .
//    		  g h`,
// 	)
// }

// func TestScreenClearOnLine(t *testing.T) {
// 	testClearing := func(testName string,
// 		x, y int,
// 		n int,
// 		visualInitial, visualExpected string) {

// 		scr := makeScreenFromVisualForTest(visualInitial, testName)
// 		expected := makeScreenFromVisualForTest(visualExpected, testName)

// 		scr.ClearOnLine(x, y, n)

// 		diff := ScreenDiff(expected, scr)
// 		if len(diff) > 0 {
// 			t.Errorf("failure:\n%s\n%s", testName, diff)
// 		}
// 	}

// 	testClearing(
// 		"clear in middle",
// 		1, 1,
// 		2,
// 		` a b c d
// 		 ^d e f g`,
// 		` a b c d
// 		 ^d . . g`,
// 	)
// }

// func TestScreenClearLines(t *testing.T) {
// 	testDeletion := func(testName string,
// 		y, n int,
// 		visualInitial, visualExpected string) {

// 		scr := makeScreenFromVisualForTest(visualInitial, testName)
// 		expected := makeScreenFromVisualForTest(visualExpected, testName)

// 		scr.ClearLines(y, n)

// 		diff := ScreenDiff(expected, scr)
// 		if len(diff) > 0 {
// 			t.Errorf("failure:\n%s\n%s", testName, diff)
// 		}
// 	}

// 	testDeletion(
// 		"clear 2 lines",
// 		1, 2,
// 		` a b
// 	     ^c d
// 	      e f
// 	      g h`,
// 		` a b
// 	     ^. .
// 	      . .
// 	      g h`,
// 	)
// 	testDeletion(
// 		"clear too many lines",
// 		1, 9,
// 		` a b
// 	     ^c d
// 	      e f`,
// 		` a b
// 	     ^. .
// 	      . .`,
// 	)
// }

func TestScreenResize(t *testing.T) {
	testResize := func(testName string,
		w, h int,
		visualInitial, visualExpected string) {

		log.Infof("#### TestScreenResize %v", testName)
		scr := makeScreenFromVisualForTest(visualInitial, testName)
		expected := makeScreenFromVisualForTest(visualExpected, testName)

		scr.Resize(w, h)

		diff := ScreenDiff(expected, scr)
		if len(diff) > 0 {
			t.Errorf("failure:\n%s\n%s", testName, diff)
		}
	}

	testResize(
		"same size",
		3, 3,
		`a b c|
		 d^e f
		 g h i`,
		`a b c|
		 d^e f
		 g h i`,
	)
	testResize(
		"wider",
		4, 3,
		`^a b c
		  d e f|
		  g h .`,
		`^a b c .
	      d e f g|
		  h . . .`,
	)
	testResize(
		"narrower 3",
		3, 3,
		`^a b . .
	      c d e f|
		  g . . .`,
		`^a b .
	      c d e|
		  f g .`,
	)
	testResize(
		"fewer lines",
		3, 2,
		`^a b c
		  d e f
		  g h i`,
		`^d e f
		  g h i`,
	)
	testResize(
		"more lines",
		3, 4,
		`^a b c
		  d e f
		  . . .`,
		`^a b c
		  d e f
		  . . .
		  . . .`,
	)
	testResize(
		"narrower",
		3, 3,
		` a b . .
	      c d e^f|
		  g . . .|
		  . . . .`,
		` a b .
	      c d e|
		 ^f g .`,
	)
	testResize(
		"narrower 2",
		3, 5,
		` a b c d|
	      e f g h|
		  i .^. j|
		  k . . .`,
		` a b c|
		  d e f|
		  g h i|
		  .^. j|
		  k . . `,
	)
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func makeScreenFromVisualForTest(visual, testName string) Screen {
	screen, err := MakeScreenFromVisual(visual)
	if err != nil {
		log.Fatalf("problem with test '%v':\n%v", testName, err)
	}
	return screen
}

// func TestMakeScreenFromVisual(t *testing.T) {
// 	scr, err := MakeScreenFromVisual(`a   _
//  									  .   b|
//    									  👦   ^`)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	w, h := scr.Width(), scr.Height()
// 	if w != 2 || h != 3 {
// 		t.Error("unexpected screen size: ", w, h)
// 		return
// 	}
// 	expects := []struct {
// 		x, y int
// 		r    rune
// 	}{
// 		{0, 0, 'a'},
// 		{1, 0, ' '},
// 		{0, 1, 0},
// 		{1, 1, 'b'},
// 		{0, 2, '👦'},
// 		{1, 2, 0},
// 	}
// 	for _, expected := range expects {
// 		cell := scr.CellAt(expected.x, expected.y)
// 		if cell.Ch != expected.r {
// 			t.Errorf("rune in position %v, %v is %q instead of expected %q",
// 				expected.x, expected.y, cell.Ch, expected.r)
// 		}
// 	}
// 	for i, expected := range []bool{false, true, false} {
// 		actual := scr.lines[i].wrapped
// 		if scr.lines[i].wrapped != expected {
// 			t.Errorf("line wrap line %v: expected %v, was %v",
// 				i, expected, actual)
// 		}

// 	}
// 	actualCx, actualCy := scr.Cursor()
// 	if actualCx != 2 || actualCy != 2 {
// 		t.Errorf("cursor is  %v, %v instead of expected 2, 2", actualCx, actualCy)
// 	}
// }
