package prom

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
)

// PromMetricCounter describes the Prometheus metrics counter
type PromMetricCounter struct {
	Name    string
	Valobj  string
	Counter int
	Unit    string
	Dest    string
}

// Returns string value of prometheus metrics
func (p *PromMetricCounter) Add() {
	p.Counter += 1
}

// Returns string value of prometheus metrics
func (p *PromMetricCounter) string() string {
	return fmt.Sprintf("# TYPE %s counter\n%s{%s,unit=\"%s\"} %d\n", p.Name, p.Name, p.Valobj, p.Unit, p.Counter)
}

// Send prometheus metrics
func (p *PromMetricCounter) Send() error {
	conn, err := net.Dial("tcp", p.Dest)
	if err != nil {
		log.Error().Msgf("Failed to connect to %s, %v", p.Dest, err)
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(p.string()))
	if err != nil {
		log.Error().Msgf("Failed to send message %s, %v", p.Dest, err)
		return err
	}

	return nil
}
