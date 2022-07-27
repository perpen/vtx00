package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	log "github.com/sirupsen/logrus"
)

type input struct {
	views.BoxLayout
	text               *views.Text
	value              string
	cursor             int
	style, cursorStyle tcell.Style
}

func newInput(label string) *input {
	input := input{}
	input.style = tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorYellow)
	input.cursorStyle = input.style.Reverse(true)

	labelView := views.NewText()
	labelView.SetText(label)
	labelView.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).
		Background(tcell.ColorRed))

	textView := views.NewText()
	textView.SetText(" ")
	textView.SetStyleAt(0, input.cursorStyle)
	textView.SetStyle(input.style)

	input.text = textView

	input.SetOrientation(views.Horizontal)
	input.AddWidget(labelView, .0001)
	input.AddWidget(textView, 1)
	log.Debugf("newInput: %v", &input)
	return &input
}

func (input *input) HandleEvent(ev tcell.Event) bool {
	// log.Debug("input.HandleEvent: ", ev)
	switch ev := ev.(type) {
	case *tcell.EventKey:
		runes := []rune(input.value)
		cur := input.cursor

		switch ev.Key() {
		case tcell.KeyCtrlA:
			log.Info("C-a")
			cur = 0
		case tcell.KeyCtrlE:
			log.Info("C-e")
			cur = len(runes)
		case tcell.KeyCtrlB, tcell.KeyLeft:
			log.Info("C-b")
			if cur > 0 {
				cur--
			}
		case tcell.KeyCtrlF, tcell.KeyRight:
			log.Info("C-f")
			if cur < len(runes) {
				cur++
			}
		case tcell.KeyDelete:
			log.Info("delete")
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			log.Info("backspace")
			if cur > 0 {
				runes = append(runes[:cur-1], runes[cur+0:]...)
				cur--
			}
		case tcell.KeyRune:
			runes = append(runes, 0)
			copy(runes[cur+1:], runes[cur:])
			runes[cur] = ev.Rune()
			cur++
		}
		input.value = string(runes)

		input.text.SetText(input.value + " ")
		input.text.SetStyleAt(cur+0, input.cursorStyle)
		input.cursor = cur
		return true
	}
	return false
}
