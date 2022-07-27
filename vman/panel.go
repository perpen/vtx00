package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gdamore/tcell"
	"github.com/perpen/vtx00/vterm"
	log "github.com/sirupsen/logrus"
)

var nextPanelID int

type panel struct {
	id          int
	phys        *physical
	term        *vterm.Term
	cmd         *exec.Cmd
	pty         os.File
	top         pair // Top-left corner position
	size        pair
	cur         pair
	borderStyle tcell.Style
	z           int
	termOffset  pair
	startChan   chan bool
	meta        string
}

func (c *container) newPanel(term *vterm.Term, ptmx *os.File, cmd *exec.Cmd) *panel {
	p := panel{}
	p.id = nextPanelID
	nextPanelID++
	p.phys = c.phys
	p.term = term
	p.cmd = cmd
	p.pty = *ptmx
	p.startChan = make(chan bool)

	go func() {
		<-p.startChan
		//defer p.pty.Close() //FIXME
		buf := make([]byte, 4096) // FIXME, for linux, where did I get this from?
		for {
			_, err := io.CopyBuffer(p.term, &p.pty, buf)
			if err != nil {
				switch err.(type) {
				case *os.PathError:
					log.Debugln("container.newPanel: got os.PathError, assuming process died")
					p.pty.Write([]byte("PROCESS EXITED\n")) //FIXME
					c.deathChan <- &p
					return
				default:
					log.Errorf("error reading from pty: %T: %v", err, err)
				}
			}
		}
	}()

	return &p
}

//func (p *panel) Write(p []byte) (n int, err error) {
//	return p.pty.Write(p)
//}

func (p *panel) processDamageThroughRect(dmg vterm.Damage, focused bool, r rect) {
	log.Debugf("processDamageThroughRect panel %v: dmg=%v r=%v", p.id, dmg, r)
	scr := dmg.Term.Screen

	if dmg.W * dmg.H > 0 {
		xOff := p.top.x + p.termOffset.x
		yOff := p.top.y + p.termOffset.y
		for x := dmg.X; x < dmg.X + dmg.W; x++ {
			for y := dmg.Y; y < dmg.Y + dmg.H; y++ {
				xAbs, yAbs := xOff+x, yOff+y
				if r.contains(xAbs, yAbs) {
					cell := scr.CellAt(x, y)
					if cell == nil {
						log.Error("invalid cell coords: ", x, ", ", y, " for panel ", p.id)
						continue
					}
					ch := cell.Ch
					if ch == 0 {
						ch = ' '
					}
					style := tcellStyle(&cell.Style)
					p.phys.tcs.SetContent(xAbs, yAbs, ch, nil, style)
				}
			}
		}
	}

    cx, cy := scr.Cursor()
	p.cur = pair{cx, cy}
	if focused {
		p.showCursor()
	}
}

func (p *panel) processDamage(dmg vterm.Damage, focused bool) {
	// log.Debugln("processDamage: ", dmg)
	scr := dmg.Term.Screen
	xOff := p.top.x + p.termOffset.x
	yOff := p.top.y + p.termOffset.y
	for x := dmg.X; x <= dmg.X + dmg.W; x++ {
		for y := dmg.Y; y <= dmg.Y + dmg.H; y++ {
			cell := scr.CellAt(x, y)
			if cell == nil {
				log.Error("invalid cell coords: ", x, ", ", y)
				continue
			}
			ch := cell.Ch
			if ch == 0 {
				ch = ' '
			}
			style := tcellStyle(&cell.Style)
			p.phys.tcs.SetContent(xOff+x, yOff+y, ch, nil, style)
		}
	}

    cx, cy := scr.Cursor()
	p.cur = pair{cx, cy}
	if focused {
		p.showCursor()
	}
}

func (p *panel) showCursor() {
	xOff := p.top.x + p.termOffset.x
	yOff := p.top.y + p.termOffset.y
	p.phys.tcs.ShowCursor(xOff+p.cur.x, yOff+p.cur.y)
	// p.phys.tcs.Show()
}

func (p *panel) resizeTerm(w, h int) {
	p.term.Resize(w, h)
	go func() {
		p.startChan <- true
	}()
	ws := pty.Winsize{Cols: uint16(w), Rows: uint16(h)}
	pty.Setsize(&p.pty, &ws)
}

func (p *panel) refresh() {
	p.term.Refresh()
}

func (p *panel) rect() rect {
	return rect{p.top.x, p.top.y, p.size.x, p.size.y}
}
