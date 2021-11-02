package logserver

import (
	"fmt"
	"log"
	"net"
)

type LogServer interface {

	// Logs returns channel to retreive logs. Channel will be closed when server is stopped.
	Logs() <-chan string

	// Run must be called before reading logs from channel.
	// When the received chan is closed, server starts its shutdown procedure.
	Run(<-chan struct{})
}

func New(bindIP string, port int, filters ...filter) LogServer {
	return &logserver{
		bindIP:   bindIP,
		bindPort: port,
		logCh:    make(chan string),
		filters:  filters,
	}
}

type logserver struct {
	bindIP   string
	bindPort int
	listener *net.UDPConn
	started  bool
	logCh    chan string
	filters  []filter
}

func (s *logserver) Logs() <-chan string {
	return s.logCh
}

// Run must to be invoked to start receiving logs from RouterOS
func (s *logserver) Run(stopCh <-chan struct{}) {
	if s.started {
		return
	}
	s.started = true
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(s.bindIP),
		Port: s.bindPort,
	})
	if err != nil {
		log.Fatal(err)
	}
	s.listener = listener
	go func() {
		data := make([]byte, 1024)
	nextLine:
		for s.started {
			n, _, err := s.listener.ReadFromUDP(data)
			if err != nil {
				fmt.Printf("error during read: %s", err)
			}
			line := string(data[:n])
			for _, f := range s.filters {
				if !f(line) {
					continue nextLine
				}
			}
			select {
			case s.logCh <- line:
			default:
				log.Printf("Dropped log: %s", data[:n])
			}
		}
		close(s.logCh)
	}()

	go func() {
		<-stopCh
		s.shutdown()
	}()

}

func (s *logserver) shutdown() {
	s.started = false
	s.listener.Close()
}
