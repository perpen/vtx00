package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gdamore/tcell"
	log "github.com/sirupsen/logrus"
)

type physical struct {
	tcs          tcell.Screen
	keyboardChan chan []byte
	resizeChan   chan pair
	winchChan    chan os.Signal
}

func newPhysical() *physical {
	phys := physical{}
	phys.keyboardChan = make(chan []byte)
	phys.resizeChan = make(chan pair)
	return &phys
}

func (phys *physical) start() {
	// Tcell stuff
	tcs, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err := tcs.Init(); err != nil {
		log.Fatal(err)
	}
	//tcs.EnableMouse()
	tcs.Clear()
	tcs.Show()

	phys.tcs = tcs

	// Winch detection
	winchChan := make(chan os.Signal)
	go func() {
		for {
			<-winchChan
			// Wait a bit before reading the new screen size
			time.Sleep(50 * time.Millisecond)
			w, h := phys.tcs.Size()
			log.Infoln("winch signal", w, h)
			phys.resizeChan <- pair{w, h}
		}
	}()
	signal.Notify(winchChan, syscall.SIGWINCH)

	// Read from stdin
	go func() {
		buf := make([]byte, 4096) // FIXME why 4k?
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				log.Errorln(err)
			}
			phys.keyboardChan <- buf[:n]
		}
	}()
}

func (phys *physical) stop() {
	phys.tcs.Fini()
}

func (phys *physical) putString(x, y int, s string, style tcell.Style) {
	w, _ := phys.tcs.Size()
	reader := strings.NewReader(s)
	for ; x < w; x++ {
		r, _, _ := reader.ReadRune()
		phys.tcs.SetContent(x, y, r, nil, style)
	}
}
