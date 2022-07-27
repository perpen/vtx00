package vterm

// FIXME avoid creating new slice of *Cell, use existing slice of Cell from screen
type Damage struct {
	Term       *Term
	X, Y, W, H int
}

type Rect struct {
	x, y, w, h int
}

var emptyZone = Rect{}
