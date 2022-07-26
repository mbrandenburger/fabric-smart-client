package common

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/onsi/gomega/gexec"
)

type Runner struct {
	Command       *exec.Cmd // process to be executed
	Name          string    // prefixes all output lines
	AnsiColorCode string    // colors the output
	exitCode      int
	lock          *sync.Mutex
}

func (r *Runner) Run(sigChan <-chan os.Signal, ready chan<- struct{}) error {
	r.lock = &sync.Mutex{}

	commandOut, commandErr := r.Command.Stdout, r.Command.Stdout
	if commandOut != nil {
		r.Command.Stdout = gexec.NewPrefixedWriter(
			fmt.Sprintf("\x1b[32m[o]\x1b[%s[%s]\x1b[0m ", r.AnsiColorCode, r.Name),
			commandOut,
		)
	}

	if commandErr != nil {
		r.Command.Stderr = gexec.NewPrefixedWriter(
			fmt.Sprintf("\x1b[91m[e]\x1b[%s[%s]\x1b[0m ", r.AnsiColorCode, r.Name),
			commandErr,
		)
	}

	if err := r.Command.Start(); err != nil {
		return err
	}

	exited := make(chan struct{})
	go r.monitorForExit(exited)

	logger.Infof("[%s] sleep before signal ready", r.Name)
	// TODO make this a proper check if readys
	// use GET /healthz
	time.Sleep(20 * time.Second)
	close(ready)

	for {
		select {
		case signal := <-sigChan:
			r.Command.Process.Signal(signal)

		case <-exited:
			logger.Infof("[%s] exited!!!", r.Name)

			if r.exitCode == 0 {
				return nil
			}

			return fmt.Errorf("exit status %d", r.exitCode)
		}
	}
}

func (r *Runner) monitorForExit(exited chan<- struct{}) {
	err := r.Command.Wait()
	r.lock.Lock()
	status := r.Command.ProcessState.Sys().(syscall.WaitStatus)
	if status.Signaled() {
		r.exitCode = 128 + int(status.Signal())
	} else {
		exitStatus := status.ExitStatus()
		if exitStatus == -1 && err != nil {
			r.exitCode = 254
		}
		r.exitCode = exitStatus
	}
	r.lock.Unlock()

	close(exited)
}

type CmdHelper struct {
	Out      io.Writer
	ErrOut   io.Writer
	ExitCode int
}
