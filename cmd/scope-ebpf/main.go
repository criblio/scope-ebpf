package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

func main() {
	var dest string
	var debug bool

	if os.Geteuid() != 0 {
		fmt.Println("This binary must be run with sudo for elevated privileges.")
		return
	}

	// fmt.Println(os.Args[0], "started, PID:", os.Getpid())
	flag.StringVar(&dest, "dest", "", "Destination point")
	flag.BoolVar(&debug, "debug", false, "Enable debug message")
	flag.Parse()

	// Setup Logging
	logLevel := zerolog.ErrorLevel
	if debug {
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	if dest != "" {
		server(dest)
	} else {
		loader()
	}
}
