package main

import (
	"fmt"
	"os"

	// "github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/perpen/vtx00/vparser"
	log "github.com/sirupsen/logrus"
)

type stack struct {
    views.BoxLayout
}

var app = &views.Application{}

func main() {
	vparser.InitLogging("vsel.log", log.DebugLevel)

	path := os.Args[1]
	cnt := newContent(path)
    // input := newInput("regex: ")

    stack := stack{}
	stack.SetOrientation(views.Vertical)
	stack.AddWidget(cnt, 1)
	// stack.AddWidget(input, 0.1)

	app.SetRootWidget(&stack)
	if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
}
