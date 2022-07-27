package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/perpen/vtx00/vparser"
	"github.com/perpen/vtx00/vterm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

type spy struct {
	parser *vparser.Parser
}

func newSpy() *spy {
	spy := new(spy)

	// We need a function pointer, not a method
	handler := func(spec *vparser.ControlSpec, params []int) vterm.Rect {
		spy.Handle(spec, params)
		return vterm.Rect{}
	}

	controls := vterm.AllControls
	for _, control := range controls {
		control.UserData = &handler
	}

	bindings := vparser.NewBindings(controls)
	spy.parser = vparser.NewParser(bindings, spy)
	return spy
}

func (spy *spy) Write(data []byte) (int, error) {
	spy.parser.Parse(data)
	return len(data), nil
}

func (spy spy) Print(ch rune) {
	if showPrint {
		log.Infof("print   %c", ch)
	}
}

func (spy spy) Handle(spec *vparser.ControlSpec, params []int) {
	if len(params) > 0 {
		log.Infoln("handle ", spec.Name, params)
	} else {
		log.Infoln("handle ", spec.Name)
	}
}

func (spy spy) Osc(b byte, s string) {
	log.Infoln("osc:", b, ",", s)
}

func (spy spy) Unknown(triggerString string, params []int) {
	log.Errorln("unknown", triggerString, params)
}

func (spy spy) Error(msg string) {
	log.Errorln("error:", msg)
}

var showPrint = true

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--no-print" {
		showPrint = false
	}
	logFile := "spy.log"
	fmt.Printf("logging to %v\n", logFile)
	vparser.InitLogging(logFile, log.InfoLevel)
	spy := newSpy()

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	c := exec.Command("/bin/zsh")
	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	// when the terminal is resized we receive a SIGWINCH
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			//fmt.Fprintf(os.Stderr, "resize\n")
			if err := pty.InheritSize(os.Stdin, f); err != nil {
				fmt.Fprintf(os.Stderr, "resize error: %s\n", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH

	mw := io.MultiWriter(os.Stdout, spy)

	//go io.Copy(os.Stderr, f)
	//go io.Copy(os.Stdout, f)
	go io.Copy(mw, f)
	go io.Copy(f, os.Stdin)

	c.Wait()
	log.Println("Exiting...")
}
