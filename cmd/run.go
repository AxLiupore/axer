package main

import (
	"axer/cmd/main/cgroup"
	"axer/cmd/main/constant"
	"axer/cmd/main/container"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Run container
func Run(tty bool, cmdArray []string, limit *cgroup.Limit, volume, containerName string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	containerName, err := recordContainerInfo(parent.Process.Pid, cmdArray, containerName)

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
		deleteContainerInfo(containerName)
	}

	pwd, err := os.Getwd()
	container.DeleteWorkSpace(pwd, volume)
}

// set init command
func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}

func recordContainerInfo(containerPid int, cmdArray []string, containerName string) (string, error) {
	// Randomly generate a 10-digit containerId
	id := randStringBytes(container.IDLength)
	// Set now time as container created time
	createTime := time.Now().Format("2006-01-02 15:04:05")
	// If no container name is specified, a randomly generated containerID is used
	if containerName == "" {
		containerName = id
	}

	command := strings.Join(cmdArray, "")
	containerInfo := &container.Info{
		Id:          id,
		Pid:         strconv.Itoa(containerPid),
		Command:     command,
		CreatedTime: createTime,
		Status:      container.RUNNING,
		Name:        containerName,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		logrus.Errorf("Record container info error %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	// Splice out the path to the container information file. If the directory does not exist, create it cascaded
	dirPath := fmt.Sprintf(container.InfoLocFormat, containerName)
	if err := os.MkdirAll(dirPath, constant.Perm0755); err != nil {
		logrus.Errorf("MkdirAll error %s error %v", dirPath, err)
		return "", err
	}

	// Write message into file
	fileName := filepath.Join(dirPath, container.ConfigName)
	file, err := os.Create(fileName)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		logrus.Errorf("Create file %s error %v", fileName, err)
		return "", err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		logrus.Errorf("File write string error %v", err)
		return "", err
	}

	return containerName, nil
}

// Generate random characters
func randStringBytes(n int) string {
	letterBytes := "1234567890"
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func deleteContainerInfo(containerName string) {
	dirPath := fmt.Sprintf(container.InfoLocFormat, containerName)
	if err := os.RemoveAll(dirPath); err != nil {
		logrus.Errorf("RemoveAll dir %s error %v", dirPath, err)
	}
}
