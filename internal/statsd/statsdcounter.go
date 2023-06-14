package statsd

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
)

// StatsDCounter describes the StatsD metrics counter
type StatsDCounter struct {
	Name    string
	Valobj  string
	Counter int
	Unit    string
	Dest    string
}

// Returns string value of statsD metrics
func (s *StatsDCounter) Add() {
	s.Counter += 1
}

// Returns string value of statsD metrics
func (s *StatsDCounter) string() string {
	return fmt.Sprintf("%s:%d=|c|#%s\n", s.Name, s.Counter, s.Valobj)
}

// Send statsD metrics
func (s *StatsDCounter) Send() error {
	conn, err := net.Dial("tcp", s.Dest)
	if err != nil {
		log.Error().Msgf("Failed to connect to %s, %v", s.Dest, err)
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(s.string()))
	if err != nil {
		log.Error().Msgf("Failed to send message %s, %v", s.Dest, err)
		return err
	}

	return nil
}
