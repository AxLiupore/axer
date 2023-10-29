package main

import (
	"axer/cmd/main/container"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func logContainer(containerName string) {
	logFileLocation := fmt.Sprintf(container.InfoLocFormat, containerName) + container.LogFile
	file, err := os.Open(logFileLocation)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		logrus.Errorf("Log container open file %s error %v", logFileLocation, err)
		return
	}
	content, err := os.ReadFile(logFileLocation)
	if err != nil {
		logrus.Errorf("Log container read file %s error %v", logFileLocation, err)
		return
	}
	// Output file contents to standard output
	_, err = fmt.Fprint(os.Stdout, string(content))
	if err != nil {
		logrus.Errorf("Log container Fprint  error %v", err)
		return
	}
}
