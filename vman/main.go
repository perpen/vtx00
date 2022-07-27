package main

import (
	"flag"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/perpen/vtx00/vparser"
	"github.com/perpen/vtx00/vterm"
	log "github.com/sirupsen/logrus"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func batchDamages(batch []vterm.Damage) vterm.Damage {
	bigDmg := batch[0]
	x1, y1, x2, y2 := bigDmg.X, bigDmg.Y, bigDmg.X+bigDmg.W, bigDmg.Y+bigDmg.H
	count := 1
	if len(batch) > 1 {
		for i := 1; i < len(batch); i++ {
			dmg := batch[i]
			if dmg.Term != nil {
				count++
				x1 = min(x1, dmg.X)
				y1 = min(y1, dmg.Y)
				x2 = max(x2, dmg.X+dmg.W)
				y2 = max(y2, dmg.Y+dmg.H)
			}
		}
	}
	bigDmg.X = x1
	bigDmg.Y = y1
	bigDmg.W = x2 - x1
	bigDmg.H = y2 - y1
	log.Info("batchDamages: batched ", count, "/", len(batch))
	return bigDmg
}

var cpuprofile = flag.String("cpuprofile", "cpu.prof", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "mem.prof", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		log.Info("created CPU profile")
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	vparser.InitLogging("vman.log", log.DebugLevel)
	// vparser.InitLogging("vman.log", log.WarnLevel)

	phys := newPhysical()
	phys.start()
	defer phys.stop()

	dmgBufSize := 100
	cont := newContainer(phys, dmgBufSize)

	go func() {
		for {
			cmd := exec.Command("../vcon/vcon", apiSocketPath, notifSocketPath)
			// Ensure the child is killed when we exit
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Pdeathsig: syscall.SIGTERM,
			}
			err := cmd.Start()
			if err != nil {
				log.Errorln("cannot start controller:", err)
			}

			err = cmd.Wait()
			if err == nil {
				log.Println("controller restart")
			} else {
				log.Errorf("controller error: %v\n", err)
				time.Sleep(2 * time.Second)
			}
		}
	}()
	// time.Sleep(1 * time.Second)

	// p1 := cont.newCommandPanel("/bin/bash")
	// p1.top = pair{0, 0}
	// p1.size = pair{80, 15}
	// cont.add(p1)

	//p2 := cont.newCommandPanel("/bin/bash")
	//p2.top = pair{0, 16}
	//p2.size = pair{80, 15}
	//cont.add(p2)

	if false {
		// FIXME Move this to wherever Damage is defined
		go func() {
			for {
				dmg := <-cont.damageChan
				dmgs := make([]vterm.Damage, 0)
				dmgs = append(dmgs, dmg)
				lim := len(cont.damageChan) - 1
				// lim := 5
				for i := 0; i < lim; i++ {
					select {
					case dmg = <-cont.damageChan:
						dmgs = append(dmgs, dmg)
					default:
						break
					}
				}
				bigDmg := batchDamages(dmgs)
				// log.Info("bigDmg.Term: ", bigDmg.Term)
				p := cont.panelByTerm[bigDmg.Term]
				focused := p == cont.focus
				p.processDamage(bigDmg, focused)
				phys.tcs.Show()
			}
		}()
	}

	cmdMode := false
Here:
	for {
		// Event loop
		select {
		case damage := <-cont.damageChan:
			cont.processDamage(damage)

		case p := <-cont.deathChan:
			log.Debugln("deathChan:", p.id)
			cont.del(p.id)
			cont.notify("death", p.id)

		case data := <-phys.keyboardChan:
			// FIXME - don't check only first byte, check whole buffer
			// Although we don't want to handle a pasted C-t?
			//log.Debugf("keyboardChan gave %q %v %v", ch[0], ch[0], ch)
			ch := data[0]
			switch {
			case ch == 20: // ctrl-t
				cmdMode = true
			case cmdMode:
				cmdMode = false
				switch ch {
				case 'q':
					// panic("quitting")
					break Here
				default:
					cont.notify("key", ch)
				}
			default:
				cont.focus.pty.Write(data)
			}

		case call := <-cont.apiChan:
			cmd, out := call.Cmd, call.Out
			switch {
			case cmd.CreatePanelCmd != nil:
				cont.doCreatePanel(*cmd.CreatePanelCmd, out)
			case cmd.LayoutCmd != nil:
				cont.doLayout(*cmd.LayoutCmd, out)
			case cmd.KillCmd != nil:
				cont.doKill(*cmd.KillCmd, out)
			case cmd.GetStateCmd != nil:
				cont.doGetState(*cmd.GetStateCmd, out)
			default:
				log.Errorln("unknown cmd")
				out <- "unknown"
			}

		case sz := <-phys.resizeChan:
			log.Debugf("resizeChan gave %v", sz)
			cont.notify("resize", []int{sz.x, sz.y})
		}
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
