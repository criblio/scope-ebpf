package sigdel

import "C"

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target $GOARCH -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf sigdel_bpf.c -- -I/usr/include/bpf -I.

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/rlimit"
)

type SigDelStruct struct {
	objs bpfObjects
	link link.Link
}

type sigFaultEvent struct {
	Pid     uint32
	NsPid   uint32
	Sig     uint32
	Errno   uint32
	Code    uint32
	Uid     uint32
	Gid     uint32
	Handler uint64
	Flags   uint64
	Comm    [16]byte
}

// Returns string value of metrics in prometheus format
func (sfe *sigFaultEvent) StringProm() string {
	return fmt.Sprintf("signal=\"%d\", pid=\"%d\", name=\"%s\", uid=\"%d\", gid=\"%d\"", sfe.Sig, sfe.Pid, bytes.Trim(sfe.Comm[:], "\x00"), sfe.Uid, sfe.Gid)
}

// Returns string value of metrics in statsD format
func (sfe *sigFaultEvent) StringStatsd() string {
	return fmt.Sprintf("signal=%d,pid=%d,name=%s,uid=%d,gid=%d", sfe.Sig, sfe.Pid, bytes.Trim(sfe.Comm[:], "\x00"), sfe.Uid, sfe.Gid)
}

// // Serve Signal fault events
func Serve(sigEventChan chan<- string) error {
	var err error
	fn := "signal_deliver"

	// Allow the current process to lock memory for eBPF resources.
	if err = rlimit.RemoveMemlock(); err != nil {
		return err
	}
	objs := bpfObjects{}

	// Load BPF code
	if err = loadBpfObjects(&objs, nil); err != nil {
		return err
	}
	defer objs.Close()

	// Attach BPF code
	lnk, err := link.Tracepoint("signal", fn, objs.SigDeliver, nil)
	if err != nil {
		return err
	}
	defer lnk.Close()

	rd, err := perf.NewReader(objs.Events, os.Getpagesize())
	if err != nil {
		return err
	}
	defer rd.Close()

	for {
		ev, err := rd.Read()
		if err != nil {
			return err
		}

		if ev.LostSamples != 0 {
			continue
		}

		b_arr := bytes.NewBuffer(ev.RawSample)

		var event sigFaultEvent
		if err := binary.Read(b_arr, binary.LittleEndian, &event); err != nil {
			continue
		}

		sigEventChan <- event.StringStatsd()
	}
}

// Setup Sigdel structure
func Setup() (*SigDelStruct, error) {
	var err error
	fn := "signal_deliver"

	// Allow the current process to lock memory for eBPF resources.
	if err = rlimit.RemoveMemlock(); err != nil {
		return nil, err
	}
	sigdel := new(SigDelStruct)

	// Load BPF code
	if err = loadBpfObjects(&sigdel.objs, nil); err != nil {
		return nil, err
	}

	// Attach BPF code
	sigdel.link, err = link.Tracepoint("signal", fn, sigdel.objs.SigDeliver, nil)
	if err != nil {
		sigdel.objs.Close()
		return nil, err
	}
	return sigdel, nil
}

// Teardown Sigdel structure
func (sigdel *SigDelStruct) Teardown() {
	sigdel.objs.Close()
	sigdel.link.Close()
}
