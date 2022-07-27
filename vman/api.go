package main

import (
	"encoding/json"
	"io"
	"net"
	"syscall"

	"github.com/perpen/vtx00/vman/data"
	log "github.com/sirupsen/logrus"
)

var apiSocketPath = "/tmp/vman-api.sock"
var notifSocketPath = "/tmp/vman-notif.sock"

func (c *container) server() {
	log.Debugln("vman.server")
	syscall.Unlink(apiSocketPath)
	// if err != nil {
	// 	log.Warningln(err)
	// }
	listener, err := net.Listen("unix", apiSocketPath)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("error listening on socket", err)
		}
		c.apiHandler(conn)
		log.Debugln("api loop")
	}
}

func (c *container) apiHandler(conn net.Conn) {
	log.Debugln("apiHandler")

	dec := json.NewDecoder(conn)
	var cmd data.AbstractCmd
	if err := dec.Decode(&cmd); err == io.EOF {
		log.Debugln("EOF")
		return
	} else if err != nil {
		log.Errorln(err)
		return
	}
	log.Debugln("decoded:", cmd)

	out := make(chan interface{})
	c.apiChan <- data.ApiRequest{Cmd: cmd, Out: out}
	log.Debugln("apiHandler: pushed call to apiChan")
	resp := <-out
	log.Debugf("apiHandler: returning to client: %v\n", resp)

	enc := json.NewEncoder(conn)
	if err := enc.Encode(resp); err != nil {
		log.Errorln("vman.apiHandler:", err)
		return
	}
}

// event types: key, process death, resize, readline
func (c *container) notify(evtType string, evtDetails interface{}) {
	conn, err := net.Dial("unix", notifSocketPath)
	if err != nil {
		log.Errorln("vman.notify:", err)
		return
	}
	enc := json.NewEncoder(conn)
	focus_id := -1
	if c.focus != nil {
		focus_id = c.focus.id
	}
	req := data.Event{
		Type:    evtType,
		Target:  focus_id,
		Details: evtDetails,
	}
	log.Debugln("notify: sending:", req)
	if err = enc.Encode(req); err != nil {
		log.Errorln("vman.notify:", err)
		return
	}
	// FIXME - don't close
	err = conn.Close()
	if err != nil {
		log.Errorln("vman.notify:", err)
		return
	}
}

func (c *container) panelSerial(p *panel) data.Panel {
	return data.Panel{
		ID:     p.id,
		Pos:    []int{p.top.x, p.top.y, p.size.x, p.size.y},
		Border: data.Border{},
		Meta:   p.meta,
	}
}

func (c *container) doGetState(req data.GetStateCmd, out chan interface{}) {
	log.Debugln("doGetState:", req)
	panelThings := make([]data.Panel, len(c.panels))
	for _, p := range c.panels {
		panelThing := c.panelSerial(p)
		panelThings = append(panelThings, panelThing)
	}
	focusId := -1
	if c.focus != nil {
		focusId = c.focus.id
	}
	w, h := c.phys.tcs.Size()
	out <- data.State{
		FocusID: focusId,
		Panels:  panelThings,
		Size:    []int{w, h},
	}
}

func (c *container) doCreatePanel(req data.CreatePanelCmd, out chan interface{}) {
	log.Debugln("doCreatePanel:", req)
	p := c.newCommandPanel(req.Argv[0])
	p.meta = req.Meta
	out <- c.panelSerial(p)
}

func (c *container) doLayout(req data.LayoutCmd, out chan interface{}) {
	log.Debugln("doLayout:", req)
	c.focus = c.panels[req.FocusID]
	panels := []*panel{}
	popups := []*panel{}
	for _, reqPanel := range req.Panels {
		panel := c.panels[reqPanel.ID]
        panel.z = reqPanel.Z
		panel.meta = reqPanel.Meta

		if reqPanel.Z > 0 {
			popups = append(popups, panel)
		} else {
			panels = append(panels, panel)
		}

		pos := reqPanel.Pos
		if len(pos) == 0 {
			log.Errorf("doLayout: panel %v missing position", panel.id)
			continue
		}
		panel.size = pair{pos[2], pos[3]}
		panel.top = pair{pos[0], pos[1]}

		border := reqPanel.Border
		reqBorderStyle := makeTcellStyle(border.Style.Fg, border.Style.Bg)
		if reqBorderStyle != panel.borderStyle {
			panel.borderStyle = reqBorderStyle
		}
	}
	c.borderStrategy.drawPanels(panels, defStyle, c.phys)
	c.popupBorderStrategy.drawPanels(popups, defStyle, c.phys)
	c.focus.showCursor()
	c.phys.tcs.Show()
	out <- true
}

func (c *container) doKill(req data.KillCmd, out chan interface{}) {
	log.Debugln("doKill:", req)
	out <- "hi"
}
