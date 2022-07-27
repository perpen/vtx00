package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	runewidth "github.com/mattn/go-runewidth"
	log "github.com/sirupsen/logrus"
)

func runeWidth(r rune) int { return runewidth.RuneWidth(r) }

func stringWidth(s string) int { return runewidth.StringWidth(s) }

const (
	// These runes must have width of 1
	ellipsis          = '…'
	// FIXME allow no rune at all
	leftEnd, rightEnd = '┤', '├'
)

type barComp struct {
	offset int
	s      string
	style  tcell.Style
}

type borderStrategy interface {
	drawPanels([]*panel, tcell.Style, *physical)
}

var bordersInBetween = bordersInBetweenMaybeTitleStrategy{false}
var bordersWithTitle = bordersInBetweenMaybeTitleStrategy{true}
var bordersAllAround = bordersAllAroundStrategy{}
var bordersTitleOnly = bordersTitleOnlyStrategy{}
var bordersNone = bordersNoneStrategy{}

//////////////////////////////////////////////////////////////

type bordersAllAroundStrategy struct{}

func (b bordersAllAroundStrategy) drawPanels(
	panels []*panel,
	borderStyle tcell.Style,
	phys *physical,
) {

	for _, p := range panels {
		b.drawPanel(p)
	}
}

func (b bordersAllAroundStrategy) drawPanel(p *panel) {
	log.Infof("bordersAllAroundStrategy.drawPanel: %v", p)
	p.termOffset = pair{1, 1}
	p.resizeTerm(max(p.size.x-2, 1), max(p.size.y-2, 1))

	phys := p.phys
	style := p.borderStyle

	xOff := p.top.x
	yOff := p.top.y
	log.Debugln("top:", p.top, "size:", p.size)

	for x := 0; x < p.size.x; x++ {
		phys.tcs.SetContent(xOff+x, yOff, tcell.RuneHLine, nil, style)
		phys.tcs.SetContent(xOff+x, yOff+p.size.y-1, tcell.RuneHLine, nil, style)
	}
	for y := 0; y < p.size.y; y++ {
		phys.tcs.SetContent(xOff, yOff+y, tcell.RuneVLine, nil, style)
		phys.tcs.SetContent(xOff+p.size.x-1, yOff+y, tcell.RuneVLine, nil, style)
	}

	// Corners
	phys.tcs.SetContent(xOff, yOff, tcell.RuneULCorner, nil, style)
	phys.tcs.SetContent(xOff+p.size.x-1, yOff, tcell.RuneURCorner, nil, style)
	phys.tcs.SetContent(xOff, yOff+p.size.y-1, tcell.RuneLLCorner, nil, style)
	phys.tcs.SetContent(xOff+p.size.x-1, yOff+p.size.y-1, tcell.RuneLRCorner, nil, style)

	drawBar(p, 1)
}

//////////////////////////////////////////////////////////////

type bordersNoneStrategy struct{}

func (b bordersNoneStrategy) drawPanels(
	panels []*panel,
	borderStyle tcell.Style,
	phys *physical,
) {

	for _, p := range panels {
		p.termOffset = pair{0, 0}
		p.resizeTerm(p.size.x, p.size.y)
	}
}

//////////////////////////////////////////////////////////////

type bordersTitleOnlyStrategy struct{}

func (b bordersTitleOnlyStrategy) drawPanels(
	panels []*panel,
	borderStyle tcell.Style,
	phys *physical,
) {

	for _, p := range panels {
		b.drawPanel(p)
	}
}

func (b bordersTitleOnlyStrategy) drawPanel(p *panel) {
	phys := p.phys

	p.termOffset = pair{0, 1}
	p.resizeTerm(p.size.x, p.size.y-1)

	barStyle := p.borderStyle

	xOff := p.top.x
	yOff := p.top.y
	log.Debugln("top:", p.top, "size:", p.size)

	for x := 0; x < p.size.x-0; x++ {
		phys.tcs.SetContent(xOff+x, yOff, tcell.RuneHLine, nil, barStyle)
	}

	drawBar(p, 0)
}

//////////////////////////////////////////////////////////////

type bordersInBetweenMaybeTitleStrategy struct {
	showTitle bool
}

func (b bordersInBetweenMaybeTitleStrategy) drawPanels(
	panels []*panel,
	borderStyle tcell.Style,
	phys *physical,
) {

	corners := map[int]map[int]string{}
	add := func(x, y int, corner string) {
		if corners[x] == nil {
			corners[x] = map[int]string{}
		}
		corners[x][y] += corner
	}

	// draw the edges and titles for all panels, and collect the corners
	for _, p := range panels {
		termOffset := pair{}
		termSize := p.size

		if p.top.x > 0 && (p.top.y > 0 || b.showTitle) {
			add(p.top.x, p.top.y, "es")
		}
		if p.top.y > 0 || b.showTitle {
			add(p.top.x+p.size.x-0, p.top.y, "sw")
			termOffset.y++
			termSize.y--
		}
		if p.top.x > 0 {
			add(p.top.x, p.top.y+p.size.y-0, "en")
			termOffset.x++
			termSize.x--
		}
		add(p.top.x+p.size.x-0, p.top.y+p.size.y-0, "nw")

		barStyle := p.borderStyle
		xOff := p.top.x
		yOff := p.top.y
		log.Debugln("top:", p.top, "size:", p.size)

		if p.top.y > 0 || b.showTitle {
			// top edge
			for x := 0; x < p.size.x-0; x++ {
				phys.tcs.SetContent(xOff+x, yOff, tcell.RuneHLine, nil, barStyle)
			}
		}
		if p.top.x > 0 {
			// left edge
			for y := 0; y < p.size.y-0; y++ {
				phys.tcs.SetContent(xOff, yOff+y, tcell.RuneVLine, nil, barStyle)
			}
		}

		p.termOffset = termOffset
		p.resizeTerm(termSize.x, termSize.y)
		if b.showTitle {
			drawBar(p, p.termOffset.x)
		}
	}

	// then draw the corners
	rules := map[string]rune{
		// FIXME - use bits instead of strings, then we can just OR them
		"ensw": '┼',
		"enw":  '┴',
		"esw":  '┬',
		"ens":  '├',
		"nsw":  '┤',
		"es":   '┌',
		"sw":   '┐',
		"nw":   '┘',
		"ne":   '└',
	}

	normalise := func(s string) string {
		tokens := map[rune]bool{}
		for _, token := range s {
			tokens[token] = true
		}
		result := ""
		for _, token := range "ensw" {
			if tokens[token] {
				result += string(token)
			}
		}
		return result
	}

	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	xOff := 0
	yOff := 0
	for x, yMap := range corners {
		for y, cornerRep := range yMap {
			norm := rules[normalise(cornerRep)]
			if norm != 0 {
				phys.tcs.SetContent(xOff+x, yOff+y, norm, nil, style)
			}
		}
	}
}

//////////////////////////////////////////////////////////////

func drawBar(p *panel, offset int) {
	styles := []tcell.Style{
		tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorTeal),
		tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen),
		tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue),
	}

	components := []string{
		"✅",
		fmt.Sprintf("%vx%v", p.size.x, p.size.y),
	}

	barComps := formatPanelBar("bash", styles[0], components,
		styles[1:],
		p.size.x - 2*offset,
	)

	phys := p.phys
	xOff := p.top.x + offset
	yOff := p.top.y

	for i := len(barComps) - 1; i >= 0; i-- {
		comp := barComps[i]
		log.Debugf("comp: %v", comp)
		reader := strings.NewReader(comp.s)
		x := comp.offset - 1
		phys.tcs.SetContent(xOff+x, yOff, leftEnd, nil, p.borderStyle)
		x++
		for {
			r, n, _ := reader.ReadRune()
			if n == 0 {
				break
			}
			phys.tcs.SetContent(xOff+x, yOff, r, nil, comp.style)
			x += runeWidth(r)
		}
		phys.tcs.SetContent(xOff+x, yOff, rightEnd, nil, p.borderStyle)
	}
}

func trimBarComponents(ellidableComponent string, fixedComponents []string, width int) (string, []string) {
	// "┤the title├────┤paused├┤12├"
	leftEndWidth, rightEndWidth := runeWidth(leftEnd), runeWidth(rightEnd)

	// Collect the fixed-size components we have space for
	free := width
	shown := []string{}
	for _, comp := range fixedComponents {
		cw := leftEndWidth + stringWidth(comp) + rightEndWidth
		if cw > free {
			break
		}
		free -= cw
		shown = append(shown, comp)
	}

	// Try fit the ellidable component
	ellisionFree := free - leftEndWidth - rightEndWidth
	ellided := ""
	if ellisionFree > 0 {
		ellided = ellide(ellidableComponent, ellisionFree)
	}

	return ellided, shown
}

// formatPanelBar returns a string which hopefully contains the title
// aligned to the left, and the other components pushed to the right, in
// reverse order.
// If there is not enough space we try:
// - Elliding the title
// - Then drop components that do not fit, starting from the left
func formatPanelBar(title string, titleStyle tcell.Style, components []string, compStyles []tcell.Style, width int) []barComp {
	ellided, shown := trimBarComponents(title, components, width)
	comps := make([]barComp, 0)
	offset := width

	for i := 0; i < len(shown); i++ {
		offset -= stringWidth(shown[i]) + i + 1
		comps = append(comps, barComp{offset, shown[i], compStyles[i]})
	}

	if len(ellided) > 0 {
		comps = append(comps, barComp{1, ellided, titleStyle})
	}

	return comps
}

// returns a buffer containing a possibly truncated version
// of the param, so that it occupies visually at most the give width.
func ellide(s string, maxWidth int) string {

	impossibleError := func(msg string, err error) {
		if err != nil {
			log.Errorf("bug: getting error: %v: %v", msg, err)
			panic(err)
		}
	}
	sw := stringWidth(s)
	ew := 1
	reader := strings.NewReader(s)
	buf := bytes.NewBuffer(make([]byte, 0, maxWidth))

	if sw > maxWidth {
		// How many cells do we need truncate
		truncw := maxWidth - ew
		sw = 0
		// Add truncated string
		for {
			r, _, err := reader.ReadRune()
			impossibleError("vman.truncate/ReadRune", err)
			if err != nil {
				log.Errorln("bug: getting ReadRune error", err)
			}
			rw := runewidth.RuneWidth(r)
			if sw+rw > truncw {
				break
			}
			_, err = buf.WriteRune(r)
			impossibleError("vman.truncate/WriteRune", err)
			sw += rw
		}
		// Add ellipsis
		_, err := buf.WriteRune(ellipsis)
		impossibleError("vman.truncate/WriteRune ellipsis", err)
	} else {
		_, err := buf.WriteString(s)
		impossibleError("vman.truncate/WriteString", err)
	}

	return buf.String()
}
