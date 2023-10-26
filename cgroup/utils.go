package cgroup

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
	"syscall"
)

func readFile(p string) ([]byte, error) {
	data, err := os.ReadFile(p)
	for err != nil && errors.Is(err, syscall.EINTR) {
		data, err = os.ReadFile(p)
	}
	return data, err
}

func remove(name string) error {
	if name != "" {
		return os.Remove(name)
	}
	return nil
}

func writeFile(p string, content []byte, perm fs.FileMode) error {
	err := os.WriteFile(p, content, filePerm)
	for err != nil && errors.Is(err, syscall.EINTR) {
		err = os.WriteFile(p, content, filePerm)
	}
	return err
}

// StringToUint string => uint64
func (c *Cgroup) StringToUint(content string) (uint64, error) {
	u, err := strconv.ParseUint(content, 10, 64)
	if err != nil {
		return 0, nil
	}
	return u, nil
}
