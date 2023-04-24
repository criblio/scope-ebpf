package main

import (
	"fmt"
	"os"
	"time"

	"github.com/criblio/scope-ebpf/internal/ebpf/oom"
	"github.com/criblio/scope-ebpf/internal/teardown"
)

const timeout = 60 * time.Second

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("This binary must be run with sudo for elevated privileges.")
		return
	}
	fmt.Println("Loader started, PID: ", os.Getpid())

	o, err := oom.Setup()
	if err != nil {
		fmt.Println("oom.Setup failed")
		return
	}
	defer o.Teardown()
	os.Exit(teardown.TeardownProc(timeout))
}
