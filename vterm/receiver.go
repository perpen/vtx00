package vterm

import (
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
)

type parserRcvr struct {
	term *Term
}

// Print is a vparser.Parser callback
func (p parserRcvr) Print(ch rune) {
	// log.Infof("term.print: %q", ch)

	style := p.term.settings.style
	x, y, w, h := p.term.print(string(ch), style)
	if p.term.fineDamage {
		p.term.pushDamage(Rect{x, y, w, h})
	}
}

// Handle is a vparser.Parser callback
func (p parserRcvr) Handle(spec *vparser.ControlSpec, params []int) {
	if spec.UserData == nil {
		log.Warnf("%v %v - not implemented", spec.Name, params)
		return
	}
	if len(params) > 0 && len(spec.Selections) > 0 {
		log.Debugf("%v %v -> %v", spec.Name, params, spec.Selections[params[0]])
	} else {
		log.Debugf("%v %v", spec.Name, params)
	}
	handler := spec.UserData.(func(spec *vparser.ControlSpec, params []int, term *Term) Rect)
	zone := handler(spec, params, p.term)
	if p.term.fineDamage {
		p.term.pushDamage(zone)
	}
}

// Osc is a vparser.Parser callback.
func (p parserRcvr) Osc(b byte, s string) {
	log.Infoln("term.osc:", b, ",", s)
}

// Unknown is a vparser.Parser callback.
func (p parserRcvr) Unknown(triggerString string, params []int) {
	log.Errorln("unknown", triggerString, params)
}

// Error is a vparser.Parser callback.
func (p parserRcvr) Error(msg string) {
	log.Infoln("term.error:", msg)
}

func unimplemented(spec *vparser.ControlSpec, params []int, term *Term) {
	log.Warnf("%v %v - not implemented", spec.Name, params)
}
