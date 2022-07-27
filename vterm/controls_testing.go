package vterm

// These things are required only for testing, but since I find it nice
// to define the unit tests for a control in the same file as the control
// itself, we compile this in.

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/perpen/vtx00/vparser"
	"github.com/perpen/vtx00/vscreen"
	log "github.com/sirupsen/logrus"
)

type testState struct {
	visualScreen string
	settings     TermSettings
}

func testImpl(specName, testName string, params []int,
	initialVisual, expectedVisual testState,
	expectedZone Rect,
	t *testing.T) {

	log.Infof("testing %v - %v", specName, testName)

	spec := AllControls[specName]
	handler := spec.UserData.(func(spec *vparser.ControlSpec, params []int, term *Term) Rect)

	mkState := func(visual testState) termState {
		state, err := makeStateFromVisual(visual)
		if err != nil {
			log.Fatalf("problem with test '%v %v':\n%v", specName, testName, err)
		}
		return state
	}

	initial := mkState(initialVisual)
	expected := mkState(expectedVisual)

	term := Term{}
	term.settings = initial.settings
	term.Screen = initial.Screen
	params = spec.MergeDefaults(params)
	actualZone := handler(&spec, params, &term)

	actual := termState{
		settings: term.settings,
		Screen:   term.Screen,
	}
	diff := statesDiff(expected, actual)
	if len(diff) > 0 {
		t.Errorf("failure:\n%s %s\n%s", specName, testName, diff)
	}
	if actualZone != expectedZone {
		t.Errorf("failure:\n%s %s\n%v != %v", specName, testName,
			expectedZone, actualZone)
	}
}

// Returns a description of the differences between the states
// Empty string if they are equal
func statesDiff(s1, s2 termState) string {
	scrDiff := vscreen.ScreenDiff(s1.Screen, s2.Screen)
	setDiff := settingsDiff(s1.settings, s2.settings)
	return scrDiff + setDiff
}

// Returns a string showing the differences between the settings
// Empty string if they are equal
func settingsDiff(stg1, stg2 TermSettings) string {
	var buf bytes.Buffer
	check := func(name string, val1, val2 interface{}) {
		if val1 != val2 {
			buf.WriteString(fmt.Sprintf("     %s:\t%v\t%v\n", name, val1, val2))
		}
	}

	check("title", stg1.title, stg2.title)
	check("bold", stg1.style.Bold, stg2.style.Bold)
	check("reverse", stg1.style.Reverse, stg2.style.Reverse)
	check("italics", stg1.style.Italics, stg2.style.Italics)
	check("underlined", stg1.style.Underlined, stg2.style.Underlined)
	check("fg", stg1.style.Fg, stg2.style.Fg)
	check("bg", stg1.style.Bg, stg2.style.Bg)
	check("savedCursor", stg1.savedCursor, stg2.savedCursor)
	check("regtop", stg1.regtop, stg2.regtop)
	check("regbot", stg1.regbot, stg2.regbot)
	check("DECCKM", stg1.DECCKM, stg2.DECCKM)

	settingsDiff := buf.String()
	if len(settingsDiff) > 0 {
		buf.WriteString("settings (expected, actual):\n")
		buf.WriteString(settingsDiff)
	}

	return buf.String()
}

func makeStateFromVisual(visual testState) (termState, error) {
	if len(visual.visualScreen) == 0 {
		// No screen was provided, we make a tiny one
		visual.visualScreen = ".^"
	}
	scr, err := vscreen.MakeScreenFromVisual(visual.visualScreen)
	if err != nil {
		return termState{}, err
	}
	settings := mergeSettingsWithDefault(visual.settings, scr)
	state := termState{
		settings: settings,
		Screen:   scr,
	}
	return state, nil
}

func mergeSettingsWithDefault(explicit TermSettings, scr vscreen.Screen) TermSettings {
	merged := explicit
	if explicit.title == "" {
		merged.title = defaultSettings.title
	}
	if !explicit.style.Bold {
		merged.style.Bold = defaultSettings.style.Bold
	}
	if !explicit.style.Reverse {
		merged.style.Reverse = defaultSettings.style.Reverse
	}
	if !explicit.style.Italics {
		merged.style.Italics = defaultSettings.style.Italics
	}
	if !explicit.style.Underlined {
		merged.style.Underlined = defaultSettings.style.Underlined
	}
	if explicit.style.Fg == "" {
		merged.style.Fg = defaultSettings.style.Fg
	}
	if explicit.style.Bg == "" {
		merged.style.Bg = defaultSettings.style.Bg
	}
	if explicit.regtop == 0 && explicit.regbot == 0 {
		// the spec says region is at least 2 lines high, thus we can use
		// these zero values to signify we want the region to be fullscreen
		merged.regtop = 0
		merged.regbot = scr.Height() - 1
	}
	return merged
}
