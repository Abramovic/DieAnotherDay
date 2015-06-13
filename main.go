package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
	"os/signal"
	"syscall"
	"log"
	"github.com/sevlyar/go-daemon"
)

var (
	daemonize = flag.String("d", "false", "process to run")
	seconds = flag.Int("s", 1, "seconds between retries")
	terminate = make(chan os.Signal, 1)
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	if *seconds < 1 {
		*seconds = 1
	}
	os.Args = flag.Args()

	if len(os.Args) == 0 {
		fmt.Println("[DieAnotherDay] Not enough commands")
		return
	}

	signal.Notify(terminate,
	  syscall.SIGINT,
	  syscall.SIGTERM,
	  syscall.SIGQUIT)

	cntxt := &daemon.Context{
		PidFileName: "dieanotherday.pid",
		PidFilePerm: 0644,
		LogFileName: "dieanotherday.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{""},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		fmt.Println(err)
	}
	if d != nil && *daemonize == "true" {
		return
	}
	defer cntxt.Release()
	runProcess()
}

func runProcess() {
	process, err := os.StartProcess(os.Args[0], os.Args[1:], &os.ProcAttr{})
	if err != nil {
		log.Println("[DieAnotherDay] Failure:", err.Error())
		return
	}

	state, err := process.Wait()
	if err != nil {
		log.Println("[DieAnotherDay] Error:", err.Error())
	}

	go func() {
		for {
			select {
			case <-terminate:
				err := process.Kill()
				if err != nil {
					log.Println("[DieAnotherDay] Error:", err.Error())
				}
				log.Println("[DieAnotherDay] Until next time, Mr Bond.")
				os.Exit(1)
			}
		}
		<-terminate
	}()

	var valid bool
	for !valid {
		valid = state.Exited()
		time.Sleep(time.Second * time.Duration(*seconds))
	}
	if state != nil {
		log.Println("[DieAnotherDay] Restarting process:", state.String())
	}
	runProcess()
}
