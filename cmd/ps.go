package main

import (
	"axer/cmd/main/container"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"text/tabwriter"
)

func PsContainers() {
	// Read all dirs in the directory where container information stored
	files, err := os.ReadDir(container.InfoLoc)
	if err != nil {
		logrus.Errorf("read dir %s error %v", container.InfoLoc, err)
		return
	}
	containers := make([]*container.Info, 0, len(files))
	for _, file := range files {
		tmpContainer, err := getContainerInfo(file)
		if err != nil {
			logrus.Errorf("get container info error %v", err)
			continue
		}
		containers = append(containers, tmpContainer)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	_, err = fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	if err != nil {
		logrus.Errorf("Fprint error %v", err)
	}
	for _, item := range containers {
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreatedTime)
		if err != nil {
			logrus.Errorf("Fprint error %v", err)
		}
	}
	if err = w.Flush(); err != nil {
		logrus.Errorf("Flush error %v", err)
	}
}

// get container info from file
func getContainerInfo(file os.DirEntry) (*container.Info, error) {
	// Concatenate the full path based on the file name
	containerName := file.Name()
	configFileDir := fmt.Sprintf(container.InfoLocFormat, containerName)
	configFileDir = configFileDir + container.ConfigName
	// Read config from container's config
	content, err := os.ReadFile(configFileDir)
	if err != nil {
		logrus.Errorf("read file %s error %v", configFileDir, err)
		return nil, err
	}
	info := new(container.Info)
	if err = json.Unmarshal(content, info); err != nil {
		logrus.Errorf("json unmarshal error %v", err)
		return nil, err
	}

	return info, nil
}
