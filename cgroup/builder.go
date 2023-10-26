package cgroup

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// Builder builds cgroup directories
// available: cpuacct, memory, pids
type Builder struct {
	CPU    bool
	CPUSet bool
	Memory bool
}

// NewBuilder create a new Builder to controller a cgroup
func NewBuilder() *Builder {
	return &Builder{}
}

// WithCPU includes cpu cgroup
func (b *Builder) WithCPU() *Builder {
	b.CPU = true
	return b
}

// WithCPUSet includes cpuset cgroup
func (b *Builder) WithCPUSet() *Builder {
	b.CPUSet = true
	return b
}

// WithMemory includes memory cgroup
func (b *Builder) WithMemory() *Builder {
	b.Memory = true
	return b
}

// FilterByEnv reads /proc/cgroups and filter out non-exists ones
func (b *Builder) FilterByEnv() (*Builder, error) {
	m, err := b.getAvailableController()
	if err != nil {
		return b, err
	}
	b.CPU = b.CPU && m["cpu"]
	b.CPUSet = b.CPUSet && m["cpuset"]
	b.Memory = b.Memory && m["memory"]
	return b, nil
}

// getAvailableController reads /sys/fs/cgroup/cgroup.controllers to get all controller
func (b *Builder) getAvailableController() (map[string]bool, error) {
	c, err := readFile(path.Join(basePath, cgroupControllers))
	if err != nil {
		return nil, err
	}
	m := make(map[string]bool)
	f := strings.Fields(string(c))
	for _, v := range f {
		m[v] = true
	}
	return m, nil
}

func (b *Builder) loopControllerNames(f func(name string)) {
	for _, t := range []struct {
		name    string
		enabled bool
	}{
		{"cpu", b.CPU},
		{"cpuset", b.CPUSet},
		{"memory", b.Memory},
	} {
		if t.enabled {
			f(t.name)
		}
	}
}

func (b *Builder) controllerNames() []string {
	s := make([]string, 0, 5)
	b.loopControllerNames(func(name string) {
		s = append(s, name)
	})
	return s
}

// String prints the build properties
func (b *Builder) String() string {
	s := b.controllerNames()
	return fmt.Sprintf("cgroup builder: [%s]", strings.Join(s, ", "))
}

// Build create a new cgroup
func (b *Builder) Build(name string) (cg Cgroup, err error) {
	p := path.Join(basePath, name)
	defer func() {
		if err != nil {
			_ = remove(p)
		}
	}()
	// mkdir
	if err := os.Mkdir(p, dirPerm); err != nil {
		return cg, err
	}
	return Cgroup{p}, nil
}
