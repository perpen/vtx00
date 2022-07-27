package vscreen

type Pair struct {
	X, Y int
}

// FIXME make the nil value usable?
// Possible if the default colors for fg/bg are static?
// FIXME use bits
type Style struct {
	Bold       bool
	Reverse    bool
	Italics    bool
	Underlined bool
	Fg         string
	Bg         string
}

var defaultStyle = Style{
	Bold:       false,
	Reverse:    false,
	Italics:    false,
	Underlined: false,
	Fg:         "white",
	Bg:         "black",
}

// FIXME make the nil value usable
type Cell struct {
	Ch          rune
	unprintable rune
	Style       Style
	width       int
}

func (c *Cell) set(ch rune, width int, style Style) {
	c.Ch = ch
	c.unprintable = 0
	c.width = width
	c.Style = style
}

func (c *Cell) clear() {
	c.Ch = 0
	c.unprintable = 0
	c.width = 0
	c.Style = Style{}
}
