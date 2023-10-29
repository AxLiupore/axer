package main

import (
	"axer/cmd/main/cgroup"
	"axer/cmd/main/container"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func Run(tty bool, cmdArray []string, limit *cgroup.Limit, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}
	// create a builder
	builder, err := cgroup.NewBuilder().WithCPU().WithCPUSet().WithMemory().FilterByEnv()
	if err != nil {
		return
	}
	cg, err := builder.Build("axer")
	defer func(cg *cgroup.Cgroup) {
		_ = cg.Destroy()
	}(&cg)
	// set cgroup resource limit
	err = cg.SetLimit(limit)
	if err != nil {
		return
	}
	// add pid into this process
	err = cg.AddProc(parent.Process.Pid)
	if err != nil {
		return
	}
	// send command args to NewParentProcess
	sendInitCommand(cmdArray, writePipe)

	rootPath, _ := os.Getwd()
	txt, _ := os.ReadFile(filepath.Join(rootPath, "util", "banner.txt"))
	fmt.Println(string(txt))

	// if tty parent process block
	if tty {
		_ = parent.Wait()
	}

	//pwd, err := os.Getwd()
	//container.DeleteWorkSpace(pwd, volume)
}

// set init command
func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}
