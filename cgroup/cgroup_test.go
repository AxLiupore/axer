package cgroup

import (
	"testing"
)

func BenchmarkCgroup(b *testing.B) {
	builder, err := NewBuilder().WithCPU().WithCPUSet().WithMemory().FilterByEnv()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cg, err := builder.Build("test")
		if err != nil {
			b.Fatal(err)
		}
		if err := cg.SetCPUSet([]byte("0")); err != nil {
			b.Fatal(err)
		}
		if err := cg.SetMemoryLimit([]byte("4096")); err != nil {
			b.Fatal(err)
		}
		if err := cg.SetCPU("1"); err != nil {
			b.Fatal(err)
		}
		err = cg.Destroy()
		if err != nil {
			return
		}
	}
}
