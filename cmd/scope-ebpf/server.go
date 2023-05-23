package main

import (
	"github.com/criblio/scope-ebpf/internal/ebpf/oom"
	"github.com/criblio/scope-ebpf/internal/ebpf/sigdel"
	"github.com/criblio/scope-ebpf/internal/prom"
	"github.com/rs/zerolog/log"
)

// scope-ebpf server will start the loader in server mode
func server(scopeServer string) {
	log.Debug().Msgf("Scope-ebpf will send data to: %s", scopeServer)
	// Serve OOM
	oomEventChan := make(chan string, 25)
	go oom.Serve(oomEventChan)
	signalEventChan := make(chan string, 50)
	go sigdel.Serve(signalEventChan)

	oomElement := prom.PromMetricCounter{
		Name:    "oom_kill",
		Counter: 0,
		Unit:    "process",
		Dest:    scopeServer,
	}

	signalElement := prom.PromMetricCounter{
		Name:    "signal_fault",
		Counter: 0,
		Unit:    "signal",
		Dest:    scopeServer,
	}

	for {
		select {
		case oomEvent := <-oomEventChan:
			log.Debug().Msgf("Out of memory Event: %s", oomEvent)
			oomElement.Valobj = oomEvent
			oomElement.Add()
			oomElement.Send()
		case sigEvent := <-signalEventChan:
			log.Debug().Msgf("Signal Event: %s", sigEvent)
			signalElement.Valobj = sigEvent
			signalElement.Add()
			signalElement.Send()
		}
	}
}
