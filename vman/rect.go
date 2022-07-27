package main

import (
// log "github.com/sirupsen/logrus"
)

type rect struct{ x, y, w, h int }

func (r rect) isEmpty() bool {
	return r.w*r.h == 0
}

func (r rect) contains(x, y int) bool {
	return r.x <= x && x < r.x+r.w && r.y <= y && y < r.y+r.h
}

func (a rect) minus(b rect) []rect {
	abs := func(n int) int {
		if n < 0 {
			return -n
		}
		return n
	}
	ordered := func(low, n, high int) bool {
		return low < n && n < high
	}

	inter := a.intersect(b)
	if inter.isEmpty() {
		return []rect{a}
	}

	// FIXME too many allocs
	ay2 := a.y + a.h
	bx2 := b.x + b.w
	by2 := b.y + b.h

	remainders := []rect{}

	// log.Infof("a: %v  b: %v", a, b)
	if ordered(a.y, b.y, ay2) {
		// log.Info("top band")
		remainders = append(remainders, rect{a.x, a.y, a.w, b.y - a.y})
	}
	if a.contains(b.x, b.y) {
		// log.Info("left side")
		remainders = append(remainders, rect{a.x, b.y, abs(b.x - a.x), min(abs(b.x-a.x), b.h)})
	}
	if a.contains(bx2, b.y) {
		// log.Info("right side")
		remainders = append(remainders, rect{bx2, b.y, abs(a.x - b.x), min(abs(a.x-b.x), b.h)})
	}
	if ordered(a.y, by2, ay2) {
		// log.Info("bottom band")
		remainders = append(remainders, rect{a.x, by2, a.w, min(ay2-by2, b.h)})
	}
	// log.Infof("remainders: %v", remainders)
	return remainders
}

func (a rect) intersect(b rect) rect {
	x := max(min(a.x, a.x+a.w), min(b.x, b.x+b.w))
	y := max(min(a.y, a.y+a.w), min(b.y, b.y+b.h))
	x2 := min(max(a.x, a.x+a.h), max(b.x, b.x+b.w))
	y2 := min(max(a.y, a.y+a.h), max(b.y, b.y+b.h))
	if x < x2 && y < y2 {
		return rect{x, y, x2 - x, y2 - y}
	}
	return rect{}
}

func (a rect) minusMany(bs []rect) []rect {
	as := []rect{a}
	return manyRectsMinusMany(as, bs)
}

// FIXME needed?
func manyRectsMinusMany(as, bs []rect) []rect {
	remainders := []rect{}
	for _, a := range as {
		for _, b := range bs {
			cs := a.minus(b)
			for _, c := range cs {
				if !c.isEmpty() {
					remainders = append(remainders, c)
				}
			}
		}
	}
	return remainders
}
