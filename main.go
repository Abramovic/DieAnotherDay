package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"
)

var (
	timeout = flag.Int("t", 0, "timeout before successful exit")
	seconds = flag.Int("s", 1, "seconds between retries")
	process *os.Process
	err     error
)

func startProcess() {
	log.Println("[DieAnotherDay] Starting:", os.Args[0], os.Args[1:])
	process, err = os.StartProcess(os.Args[0], os.Args[1:], &os.ProcAttr{Env: os.Environ()})
	if err != nil || process == nil {
		log.Println("[DieAnotherDay] Error:", err)
		return
	}
	log.Println("[DieAnotherDay] Waiting on Process:", process.Pid)
}

func killProcess() {
	if process != nil {
		err := process.Kill()
		log.Println("[DieAnotherDay] Killed:", err)
	}
}

func restartProcess() {
	killProcess()
	log.Println("[DieAnotherDay] Attempting Restart")
	runProcess()
}

func checkTimeout() {
	timeout := *timeout
	if timeout == 0 {
		timeout = 30
	}
	now := time.Now()
	var valid bool
	for !valid {
		duration := int((time.Since(now) * time.Second))
		if duration >= (timeout) {
			valid = true
		}
		time.Sleep(time.Second * 30)
	}
	log.Println("[DieAnotherDay] Time Out:", process.Pid)
	restartProcess()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	log.Println("[DieAnotherDay] Good evening, Mr Bond.")

	if *seconds < 1 {
		*seconds = 1
	}
	os.Args = flag.Args()

	if len(os.Args) == 0 {
		fmt.Println("[DieAnotherDay] The world is not enough, and neither are your commands. ")
		return
	}

	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile("dieanotherday.pid", []byte(pid), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write pid file: %v\n", err)
	}
	runProcess()
}

func childFinished() {
	log.Println("[DieAnotherDay] Until next time, Mr Bond.")
	os.Exit(1)
}

func runProcess() {
	// If we see a signal from the OS try to restart the process
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Kill, os.Interrupt)
	go func() {
		<-terminate
		restartProcess()
		return
	}()

	startProcess()
	go checkTimeout()
	state, err := process.Wait()
	if err != nil {
		log.Println("[DieAnotherDay] Failure:", err)
	}
	log.Println("[DieAnotherDay] Running:", state.Pid())

	var valid bool
	for !valid {
		valid = state.Exited()
		if valid == false && state.Success() == true {
			childFinished()
		}
		time.Sleep(time.Second * time.Duration(*seconds))
	}
	if state != nil {
		log.Println("[DieAnotherDay] Restarting:", state.String())
	}
}
