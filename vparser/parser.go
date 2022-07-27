package vparser

import (
	"fmt"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"
)

// If you increase this then you'll need to use a type larger than uint32 as key for the
// parser maps.
const maxIntermediates = 2

// Parser .
type Parser struct {
	state uint
	// We reserve 2 bytes for controlSet and final char
	sequence         [1 + maxIntermediates + 1]byte
	numIntermediates int
	// Enough for utf-8
	printable     [4]byte
	numPrintable  int
	oscString     [1024]byte
	oscStringLen  uint
	ignoreFlagged bool
	params        [16]int
	numParams     int
	bindings      *Bindings
	receiver      Receiver
}

// Receiver objects are notified of each control function identified by the Parser.
type Receiver interface {
	Print(ch rune)
	// On C0, C1, CSI, ESC input
	Handle(spec *ControlSpec, params []int)
	Osc(b byte, s string)
	// Valid input but no control function registered for it
	Unknown(triggerString string, params []int)
	// Invalid input
	Error(msg string)
}

func NewParser(bindings *Bindings, receiver Receiver) *Parser {
	parser := new(Parser)
	parser.state = stateGround
	parser.bindings = bindings
	parser.receiver = receiver
	return parser
}

// Parse parses the byte array and calls the Receiver methods.
func (parser *Parser) Parse(data []byte) {
	for _, ch := range data {
		change := stateTable[parser.state-1][ch]
		parser.doStateChange(change, ch)
	}
}

// Looks up the spec registered against the sequence and invokes the receiver.
func (parser *Parser) dispatch(set ControlSet, ch byte) {
	if set == SetESC && parser.numIntermediates > 0 {
		// To keep the receiver.Handle call simple, and the implementations
		// consistent with the CSI controls, we transform the intermediates
		// into params.
		// FIXME - rehaul all this

		// log.Printf("dispatch 1: ESC: %q %v - %v inter",
		// 		   ch, parser.params[:parser.numParams], parser.numIntermediates)
	    parser.params[0] = int(ch)
		for i := 1; i < parser.numIntermediates; i++ {
			parser.params[i] = int(parser.sequence[i-1])
		}
		parser.numParams = parser.numIntermediates
		ch = parser.sequence[0]
		parser.numIntermediates = 0
		// log.Printf("dispatch 2: ESC: %q %v - %v inter",
		// 		   ch, parser.params[:parser.numParams], parser.numIntermediates)
	}

	parser.sequence[parser.numIntermediates] = ch
	seq := parser.sequence[:parser.numIntermediates+1]

	trigger := Trigger{
		Set:      set,
		Sequence: seq,
	}

	spec, found := parser.bindings.specForTrigger(trigger)
	if found {
		params := spec.MergeDefaults(parser.params[:parser.numParams])
		parser.receiver.Handle(spec, params)
	} else {
		triggerString := trigger.pretty()
		parser.receiver.Unknown(triggerString, parser.params[:parser.numParams])
	}
}

// Makes the bytes (possibly utf-8) into a rune and sends it to the receiver.
func (parser *Parser) print(ch byte) {
	parser.printable[parser.numPrintable] = ch
	parser.numPrintable++
	bytes := parser.printable[:parser.numPrintable]
	r, _ := utf8.DecodeRune(bytes)
	if r == utf8.RuneError {
		parser.receiver.Error(fmt.Sprintf("error decoding rune from %v", bytes))
	}
	parser.receiver.Print(r)
}

// nolint:gocyclo
func (parser *Parser) doAction(action uint, ch byte) {
	receiver := parser.receiver

	//glog.Infoln("doAction: numPrintable:", parser.numPrintable)
	switch action {
	case actionUtf:
		parser.printable[parser.numPrintable] = ch
		parser.numPrintable++

	case actionPrint:
		parser.print(ch)
		fallthrough

	case actionClear:
		parser.numPrintable = 0
		parser.numIntermediates = 0
		parser.numParams = 0
		parser.oscStringLen = 0
		parser.ignoreFlagged = false

	case actionExecute:
		parser.dispatch(SetC01, ch)
	case actionCsiDispatch:
		parser.dispatch(SetCSI, ch)
	case actionEscDispatch:
		parser.dispatch(SetESC, ch)

	case actionOscStart:
		parser.oscStringLen = 0
	case actionOscPut:
		parser.oscString[parser.oscStringLen] = ch
		parser.oscStringLen++
	case actionOscEnd:
		// FIXME - From http://invisible-island.net/xterm/ctlseqs/ctlseqs.html:
		// - Parse params before the string
		receiver.Osc(1, "my title")

	case actionHook:
		fallthrough
	case actionPut:
		fallthrough
	case actionUnhook:
		log.Infoln("action not implemented:", action)

	case actionCollect:
		if parser.numIntermediates+1 > maxIntermediates {
			log.Infoln("OOPS too many intermediates")
			parser.ignoreFlagged = true
		} else {
			parser.sequence[parser.numIntermediates] = ch
			parser.numIntermediates++
		}

	case actionParam:
		/* process the param character */
		if ch == ';' || ch == '$' {
			parser.numParams++
			parser.params[parser.numParams-1] = 0
		} else {
			/* the character is a digit */
			if parser.numParams == 0 {
				parser.numParams = 1
				parser.params[0] = 0
			}
			currentParam := parser.numParams - 1
			parser.params[currentParam] *= 10
			parser.params[currentParam] += int(ch - '0')
		}

	case actionIgnore:
		/* do nothing */

	default:
		receiver.Error("parsing error")
	}
}

func (parser *Parser) doStateChange(change uint, ch byte) {
	/* A state change is an action and/or a new state to transition to. */
	newState := stateChangeState(change)
	action := stateChangeAction(change)

	//log.Infoln("doStateChange: action:", actionNames[action], "state:", stateNames[newState])

	if newState != 0 {
		/* Perform up to three actions:
		 *   1. the exit action of the old state
		 *   2. the action associated with the transition
		 *   3. the entry action of the new state
		 */
		exitAction := exitActions[parser.state-1]
		entryAction := entryActions[newState-1]

		if exitAction != 0 {
			//log.Infoln("doStateChange: exitAction:", actionNames[exitAction])
			parser.doAction(exitAction, 0)
		}
		if action != 0 {
			parser.doAction(action, ch)
		}
		if entryAction != 0 {
			//log.Infoln("doStateChange: entryAction:", actionNames[entryAction])
			parser.doAction(entryAction, 0)
		}

		parser.state = newState
	} else {
		parser.doAction(action, ch)
	}
}

func stateChangeAction(stateChange uint) uint {
	return stateChange & 0x0F
}

func stateChangeState(stateChange uint) uint {
	return stateChange >> 4
}
