package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// NewParentProcess create a new process and configuration to container init process
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		logrus.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWCGROUP | syscall.CLONE_NEWNET | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUTS,
	}
	// if flags have -it
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	// ParentProcess send readPipe to childProcess
	cmd.ExtraFiles = []*os.File{readPipe}
	pwd, err := os.Getwd()
	err = NewWorkSpace(pwd)
	if err != nil {
		return nil, nil
	}
	// this is image's path
	cmd.Dir = filepath.Join(pwd, "worker", "container")
	return cmd, writePipe
}

// NewPipe create an anonymous pipe
func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
