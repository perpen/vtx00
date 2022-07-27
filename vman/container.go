package main

import (
	"os/exec"

	"github.com/creack/pty"
	"github.com/perpen/vtx00/vman/data"
	"github.com/perpen/vtx00/vterm"
	log "github.com/sirupsen/logrus"
)

type pair struct {
	x, y int
}

type container struct {
	phys                *physical
	panels              map[int]*panel
	panelByTerm         map[*vterm.Term]*panel
	focus               *panel
	borderStrategy      borderStrategy
	popupBorderStrategy borderStrategy
	deathChan           chan *panel
	apiChan             chan data.ApiRequest
	damageChan          chan vterm.Damage
}

func newContainer(phys *physical, dmgBufSize int) container {
	c := container{}
	c.phys = phys
	c.panels = make(map[int]*panel, 0)
	c.panelByTerm = make(map[*vterm.Term]*panel)
	c.deathChan = make(chan *panel)
	c.apiChan = make(chan data.ApiRequest)
	c.damageChan = make(chan vterm.Damage, dmgBufSize)

	c.popupBorderStrategy = bordersAllAround

	// c.borderStrategy = bordersAllAround
	c.borderStrategy = bordersInBetween
	// c.borderStrategy = bordersNone
	// c.borderStrategy = bordersWithTitle
	// c.borderStrategy = bordersTitleOnly

	go c.server()
	return c
}

// FIXME - or have the container method to create the panel?
func (c *container) add(p *panel) {
	c.panels[p.id] = p
	c.panelByTerm[p.term] = p
	// c.focus = p
}

func (c *container) del(panelID int) {
	delete(c.panels, panelID)
}

// newCommandPanel .
func (c *container) newCommandPanel(cmdLine string) *panel {
	log.Debugf("newCommandPanel(%v)\n", cmdLine)
	cmd := exec.Command(cmdLine)

	// termType := "xterm"
	// cmd.Env = append(os.Environ(), fmt.Sprintf("TERM=%s", termType))

	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal(err)
	}
	t := vterm.NewTerm(ptmx, c.damageChan)
	p := c.newPanel(t, ptmx, cmd)
	log.Infoln("newCommandPanel:", p, p.cmd.Process.Pid)

	p.top = pair{3, 3}
	p.size = pair{15, 15}

	c.add(p)
	return p
}

func (c *container) processDamage(dmg vterm.Damage) {
	// FIXME subtract all rectangles from panels above this one
	p := c.panelByTerm[dmg.Term]

	blockingRects := []rect{}
	for _, p2 := range c.panels {
		if p2.z > p.z {
			blockingRects = append(blockingRects, p2.rect())
		}
	}
	if len(blockingRects) == 0 {
		blockingRects = []rect{rect{}}
	}

	focused := p == c.focus
	r := rect{
		p.top.x + p.termOffset.x + dmg.X,
		p.top.y + p.termOffset.y + dmg.Y,
		dmg.W,
		dmg.H,
	}

	visibleRects := []rect{r}
	if len(blockingRects) > 0 {
		visibleRects = r.minusMany(blockingRects)
	}

	if ! r.isEmpty() {
		log.Infof("dmg panel %v: r=%v visible=%v", p.id, r, visibleRects)
		// log.Infof("dmg above:   %v", blockingRects)
	}
	for _, visibleRect := range visibleRects {
		p.processDamageThroughRect(dmg, focused, visibleRect)
	}

	c.phys.tcs.Show()
}
