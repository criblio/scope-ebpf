package teardown

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func TeardownProc(timeout time.Duration) int {
	// Create a channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGUSR1)

	// Create a channel to implement the timeout
	timeoutChan := time.After(timeout)

	// Wait for either a signal or a timeout
	select {
	case <-signalChan:
		// fmt.Println("\nReceived signal:", sig.String())
		// fmt.Println("\nExiting")
		return 0
	case <-timeoutChan:
		// fmt.Printf("\nTimeout %v reached. Exiting...", timeout)
		return 1
	}

}
