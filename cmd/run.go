package main

import (
	cgroup2 "axer/cmd/main/cgroup"
	"axer/cmd/main/container"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func Run(tty bool, cmdArray []string, limit *cgroup2.Limit, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}
	// create a builder
	builder, err := cgroup2.NewBuilder().WithCPU().WithCPUSet().WithMemory().FilterByEnv()
	if err != nil {
		return
	}
	cg, err := builder.Build("axer")
	defer func(cg *cgroup2.Cgroup) {
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
	_ = parent.Wait()
	pwd, err := os.Getwd()
	container.DeleteWorkSpace(pwd, volume)
}

// set init command
func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	logrus.Infof("command all is %s", command)
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}
