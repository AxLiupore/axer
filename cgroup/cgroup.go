package cgroup

import (
	"path"
	"strconv"
)

type Cgroup struct {
	path string
}

// WriteFile writes cgroup file and handles potential EINTR error while writes to the slow device (cgroup)
func (c *Cgroup) WriteFile(name string, content []byte) error {
	p := path.Join(c.path, name)
	return writeFile(p, content, filePerm)
}

// ReadFile reads cgroup file and handles potential EINTR error while read to the slow device (cgroup)
func (c *Cgroup) ReadFile(name string) ([]byte, error) {
	p := path.Join(c.path, name)
	return readFile(p)
}

func (c *Cgroup) Destroy() error {
	return remove(c.path)
}

// AddProc add a process into this cgroup
func (c *Cgroup) AddProc(pid int) error {
	return c.WriteFile(cgroupProcs, []byte(strconv.FormatUint(uint64(pid), 10)))
}

// SetMemoryLimit set memory limit
func (c *Cgroup) SetMemoryLimit(content string) error {
	memory, err := c.StringToUint(content)
	if err != nil {
		return err
	}
	memory = memory * 1024 * 1024
	return c.WriteFile("memory.max", []byte(strconv.FormatUint(memory, 10)))
}

func (c *Cgroup) SetCPUSet(content string) error {
	return c.WriteFile("cpuset.cpus", []byte(content))
}

func (c *Cgroup) SetCPU(input string) error {
	quota, err := c.StringToUint(input)
	if err != nil {
		return err
	}
	quota = quota * 1000
	content := strconv.FormatUint(quota, 10) + " " + "1000000"
	return c.WriteFile("cpu.max", []byte(content))
}

// SetLimit set all resource limit
func (c *Cgroup) SetLimit(limit *Limit) error {
	if limit.CPU != "" {
		err := c.SetCPU(limit.CPU)
		if err != nil {
			return err
		}
	}
	if limit.CPUSet != "" {
		err := c.SetCPUSet(limit.CPUSet)
		if err != nil {
			return err
		}
	}
	if limit.Memory != "" {
		err := c.SetMemoryLimit(limit.Memory)
		if err != nil {
			return err
		}
	}
	return nil
}
