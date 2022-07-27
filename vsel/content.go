package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type content struct {
	parser     *vparser.Parser
	parserRcvr parserRcvr
	view       views.View
	rawString  string
	// FIXME wasteful?
	// FIXME I think a style uses 64 bits, then what's the point of using a pointer?
	rawStyles []*tcell.Style
	lastStyle tcell.Style
}

func newContent(path string) *content {
	cnt := new(content)
	cnt.parserRcvr = parserRcvr{cnt}
	var allControls = map[string]vparser.ControlSpec{
		"SGR": {
			Name:           "SGR",
			ParamsNumber:   -1,
			ParamsDefaults: []int{0},
			Triggers: []vparser.Trigger{
				// CSI m
				vparser.Trigger{
					Set:      vparser.SetCSI,
					Sequence: []byte{109},
				},
			},
			UserData: doSGR,
		},
		"LF": {
			Name:           "LF",
			ParamsNumber:   0,
			ParamsDefaults: []int{},
			Triggers: []vparser.Trigger{
				// Ctrl 0o12
				vparser.Trigger{
					Set:      vparser.SetC01,
					Sequence: []byte{10},
				},
			},
			UserData: doLF,
		},
	}
	bindings := vparser.NewBindings(allControls)
	cnt.parser = vparser.NewParser(bindings, cnt.parserRcvr)

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	// log.Infof("File contents: %s", bytes)

	_, err = cnt.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return cnt
}

func (cnt *content) Write(p []byte) (n int, err error) {
	// log.Debug("content.Write:", p)
	cnt.parser.Parse(p)
	return len(p), nil
}

func (cnt *content) print(ch rune) {
	log.Infof("term.print: %q", ch)
	cnt.rawString += string(ch)
	dup := cnt.lastStyle
	cnt.rawStyles = append(cnt.rawStyles, &dup)
	// log.Infof("print: rawStyles[%v]: %v", len(t.rawStyles)-1, t.lastStyle)
}

func (cnt *content) Draw() {
	w, _ := cnt.view.Size()
	x, y := 0, 0
	cnt.view.Clear()
	for i, ch := range cnt.rawString {
		style := cnt.rawStyles[i]
		// log.Infof("draw: rawStyles[%v]: %v", i, *style)
		cnt.view.SetContent(x, y, ch, nil, *style)
		if x == w-1 || ch == '\n' {
			x = 0
			y++
		} else {
			x++
		}
	}
}

func (cnt *content) Resize() {
	log.Debugln("content.Resize")
}

func (cnt *content) HandleEvent(ev tcell.Event) bool {
	// log.Debug("content.HandleEvent: ", ev)
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			app.Quit()
			return true
		case tcell.KeyRune:
			switch ev.Rune() {
			case '/':
				log.Info("search")
				return true
			}
			return false
		}
	}
	return false
}

func (cnt *content) SetView(view views.View) {
	cnt.view = view
}

func (cnt *content) Size() (int, int) {
	return cnt.view.Size()
}

func (cnt *content) Watch(handler tcell.EventHandler) {
}

func (cnt *content) Unwatch(handler tcell.EventHandler) {
}
