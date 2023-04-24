package oom

import (
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)
import "C"

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target $GOARCH -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf oom_bpf.c -- -I/usr/include/bpf -I.

type OomStruct struct {
	objs bpfObjects
	link link.Link
}

// Setup
func Setup() (*OomStruct, error) {
	var err error

	fn := "oom_kill_process"

	oom := new(OomStruct)
	// Allow the current process to lock memory for eBPF resources.
	if err = rlimit.RemoveMemlock(); err != nil {
		return nil, err
	}

	// Load BPF code
	if err = loadBpfObjects(&oom.objs, nil); err != nil {
		return nil, err
	}

	// Attach BPF code
	oom.link, err = link.Kprobe(fn, oom.objs.KprobeOomKillProcess, nil)
	if err != nil {
		oom.objs.Close()
		return nil, err
	}
	return oom, nil
}

func (oom *OomStruct) Teardown() {
	oom.objs.Close()
	oom.link.Close()
}
