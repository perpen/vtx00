package main

import (
	"github.com/gdamore/tcell"
	"github.com/perpen/vtx00/vscreen"
	log "github.com/sirupsen/logrus"
)

var defStyle tcell.Style

var vtermToTcellColors = map[string]tcell.Color{
	"black":   tcell.ColorBlack,
	"red":     tcell.ColorRed,
	"green":   tcell.ColorGreen,
	"yellow":  tcell.ColorYellow,
	"blue":    tcell.ColorBlue,
	"magenta": tcell.ColorPurple,
	"cyan":    tcell.ColorTeal,
	"white":   tcell.ColorWhite,
}

func tcellStyle(vs *vscreen.Style) tcell.Style {
	//log.Debugln("tcellStyle")
	fg := tcellColor(vs.Fg, tcell.ColorWhite)
	bg := tcellColor(vs.Bg, tcell.ColorBlack)

	ts := tcell.StyleDefault.
		Background(bg).
		Foreground(fg).
		Reverse(vs.Reverse).
		Bold(vs.Bold).
		Underline(vs.Underlined)

	return ts
}

func tcellColor(vname string, def tcell.Color) tcell.Color {
	color, found := vtermToTcellColors[vname]
	if !found {
		//FIXME
		//log.Warningln("vman.tcellColor: unknown color:", color)
		color = def
	}
	return color
}

func makeTcellStyle(fgAPI, bgAPI string) tcell.Style {
	log.Debugln("makeTcellStyle:", fgAPI, bgAPI)
	fg := tcellColor(fgAPI, tcell.ColorWhite)
	bg := tcellColor(bgAPI, tcell.ColorBlack)
	ts := tcell.StyleDefault.
		Foreground(fg).
		Background(bg)
	return ts
}
