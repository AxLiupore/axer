package container

import (
	"axer/cmd/main/constant"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// NewParentProcess create a new process and configuration to container init process
func NewParentProcess(tty bool, volume, containerName string) (*exec.Cmd, *os.File) {
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
	} else {
		// For background running containers, the standard output is redirected to the log file to facilitate subsequent queries
		dirPath := fmt.Sprintf(InfoLocFormat, containerName)
		if err := os.MkdirAll(dirPath, constant.Perm0755); err != nil {
			logrus.Errorf("NewParentProcess mkdir %s error %v", dirPath, err)
			return nil, nil
		}
		stdLogFilePath := dirPath + LogFile
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			logrus.Errorf("NewParentProcess create file %s error %v", stdLogFilePath, err)
			return nil, nil
		}
		cmd.Stdout = stdLogFile
	}

	// ParentProcess send readPipe to childProcess
	cmd.ExtraFiles = []*os.File{readPipe}
	pwd, err := os.Getwd()
	err = NewWorkSpace(pwd, volume)
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
