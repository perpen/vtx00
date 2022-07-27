package main

import (
	"fmt"
	"testing"

	"github.com/gdamore/tcell"
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestEllide(t *testing.T) {
	vparser.InitLogging("", log.DebugLevel)
	tests := []struct {
		s        string
		max      int // visual width
		expected string
	}{
		{"ab", 2, "ab"},
		{"八", 2, "八"},
		{"八c", 2, "…"},
		{"八cd", 3, "八…"},
		{"abc", 2, "a…"},
		{"abcd", 3, "ab…"},
		{"ab", 1, "…"},
		{"", 2, ""},
	}

	for i, test := range tests {
		actual := ellide(test.s, test.max)
		assert.Equal(t, test.expected, actual, fmt.Sprintf("test %v right: %v", i, test))
	}
}

func TestBarComponentsTrimming(t *testing.T) {
	vparser.InitLogging("", log.DebugLevel)
	tests := []struct {
		// Input params
		ellidable string
		fixed     []string
		w         int
		// Expected values
		ellided string
		shown   []string
	}{
		// "|a…||b||c|"
		{"aaa", []string{"b", "c"}, 10, "a…", []string{"b", "c"}},
		// "|…||b||c|"
		{"八a", []string{"b", "c"}, 10, "…", []string{"b", "c"}},
		// "|…||b||c|"
		{"aa", []string{"b", "c"}, 9, "…", []string{"b", "c"}},
		// "|a||b||c|"
		{"a", []string{"b", "c"}, 9, "a", []string{"b", "c"}},
		// "  |b||c|"
		{"a", []string{"b", "c"}, 8, "", []string{"b", "c"}},
		// "|b||c|"
		{"a", []string{"b", "c"}, 6, "", []string{"b", "c"}},
		// "  |b|"
		{"a", []string{"b", "c"}, 5, "", []string{"b"}},
		// "|b|"
		{"a", []string{"b", "c"}, 3, "", []string{"b"}},
		// "  "
		{"a", []string{"b", "c"}, 2, "", []string{}},
		// ""
		{"a", []string{"b", "c"}, 0, "", []string{}},
		// "|a.||bbb|" - some fixed-size components cannot fit, but ellided does
		{"aaa", []string{"bbb", "ccc"}, 9, "a…", []string{"bbb"}},
		// "|aaa|" - no fixed-size components fit, but ellided does
		{"aaa", []string{"bbbb", "ccc"}, 5, "aaa", []string{}},
	}

	for i, test := range tests {
		ellided, shown := trimBarComponents(test.ellidable, test.fixed, test.w)
		assert.Equal(t, test.ellided, ellided, fmt.Sprintf("test %v ellided: %v", i, test))
		assert.Equal(t, test.shown, shown, fmt.Sprintf("test %v: %v fixed", i, test))
	}
}

func TestFormatPanelBar(t *testing.T) {
	vparser.InitLogging("", log.DebugLevel)

	s := []tcell.Style{
		tcell.StyleDefault.Foreground(tcell.ColorBlack),
		tcell.StyleDefault.Foreground(tcell.ColorRed),
		tcell.StyleDefault.Foreground(tcell.ColorGreen),
		tcell.StyleDefault.Foreground(tcell.ColorBlue),
	}

	tests := []struct {
		// Input params
		ellidable string
		fixed     []string
		w         int
		// Expected values
		exp []barComp
	}{
		// "┤title├─┤tag├┤123├"
		//  012345678901234567
		{"title", []string{"123", "tag"}, 18, []barComp{{14, "123", s[1]}, {9, "tag", s[2]}, {1, "title", s[0]}}},
		// "┤八tle├─┤tag├┤123├"
		//  01234567890123456
		{"八tle", []string{"123", "tag"}, 18, []barComp{{14, "123", s[1]}, {9, "tag", s[2]}, {1, "八tle", s[0]}}},
		// "┤八tle├─┤八g├┤123├"
		//  01234567890123456
		{"八tle", []string{"123", "八g"}, 18, []barComp{{14, "123", s[1]}, {9, "八g", s[2]}, {1, "八tle", s[0]}}},
		// "┤title├┤tag├┤123├"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 17, []barComp{{13, "123", s[1]}, {8, "tag", s[2]}, {1, "title", s[0]}}},
		// "┤tit…├┤tag├┤123├"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 16, []barComp{{12, "123", s[1]}, {7, "tag", s[2]}, {1, "tit…", s[0]}}},
		// "┤t…├┤tag├┤123├"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 14, []barComp{{10, "123", s[1]}, {5, "tag", s[2]}, {1, "t…", s[0]}}},
		// "┤…├┤tag├┤123├"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 13, []barComp{{9, "123", s[1]}, {4, "tag", s[2]}, {1, "…", s[0]}}},
		// "──┤tag├┤123├"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 12, []barComp{{8, "123", s[1]}, {3, "tag", s[2]}}},
		// "┤tag├┤123├"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 10, []barComp{{6, "123", s[1]}, {1, "tag", s[2]}}},
		// "┤t…├┤123├"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 9, []barComp{{5, "123", s[1]}, {1, "t…", s[0]}}},
		// "┤…├┤123├"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 8, []barComp{{4, "123", s[1]}, {1, "…", s[0]}}},
		// "──┤123├"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 7, []barComp{{3, "123", s[1]}}},
		// "┤123├"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 5, []barComp{{1, "123", s[1]}}},
		// "┤t…├"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 4, []barComp{{1, "t…", s[0]}}},
		// "┤…├"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 3, []barComp{{1, "…", s[0]}}},
		// "──"
		//  01234567890123456
		{"title", []string{"123", "tag"}, 2, []barComp{}},
	}

	for i, test := range tests {
		//log.Infoln(i, test)
		comps := formatPanelBar(test.ellidable, s[0], test.fixed, s[1:], test.w)
		assert.Equal(t, test.exp, comps, fmt.Sprintf("test %v: %v", i, test))
	}
}
