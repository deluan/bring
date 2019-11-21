package bring

import (
	"errors"
	"strings"
	"time"

	"github.com/deluan/bring/protocol"
)

type SessionState int

const (
	SessionClosed SessionState = iota
	SessionHandshake
	SessionActive
)

var ErrNotConnected = errors.New("not connected")

const pingFrequency = 5 * time.Second

// Session is used to create and keep a connection with a guacd server,
// and it is responsible for the initial handshake and to send and receive instructions.
// Instructions received are put in the In channel. Instructions are sent using the Send() function
type session struct {
	In    chan *protocol.Instruction
	State SessionState
	Id    string

	tunnel   protocol.Tunnel
	logger   Logger
	done     chan bool
	config   map[string]string
	protocol string
}

// newSession creates a new connection with the guacd server, using the configuration provided
func newSession(addr string, remoteProtocol string, config map[string]string, logger Logger) (*session, error) {
	t, err := protocol.NewInetSocketTunnel(addr)
	if err != nil {
		return nil, err
	}

	err = t.Connect("")
	if err != nil {
		logger.Errorf("Error connecting to '%s': %s", addr, err)
		return nil, err
	}

	s := &session{
		In:       make(chan *protocol.Instruction, 100),
		State:    SessionClosed,
		done:     make(chan bool),
		logger:   logger,
		tunnel:   t,
		config:   config,
		protocol: remoteProtocol,
	}

	s.logger.Infof("Initiating %s session with %s", strings.ToUpper(remoteProtocol), addr)
	err = s.Send(protocol.NewInstruction("select", remoteProtocol))
	if err != nil {
		s.logger.Errorf("Failed sending 'select': %s", err)
		return nil, err
	}

	s.State = SessionHandshake
	s.startReader()

	return s, nil
}

// Terminate the current session, disconnecting from the server
func (s *session) Terminate() {
	if s.State == SessionClosed {
		return
	}
	close(s.done)
	s.State = SessionClosed
	_ = s.tunnel.SendInstruction(protocol.NewInstruction("disconnect"))
	s.tunnel.Disconnect()
}

// Send instructions to the server. Multiple instructions are sent in one single transaction
func (s *session) Send(ins ...*protocol.Instruction) error {
	for _, i := range ins {
		s.logger.Debugf("C> %s", i)
	}
	return s.tunnel.SendInstruction(ins...)
}

func (s *session) startKeepAlive() {
	go func() {
		ping := time.NewTicker(pingFrequency)
		defer ping.Stop()
		for {
			select {
			case <-ping.C:
				err := s.Send(protocol.NewInstruction("nop"))
				if err != nil {
					s.logger.Errorf("Failed ping the server: %s", err)
				}
			case <-s.done:
				return
			}
		}
	}()
}

func (s *session) startReader() {
	go func() {
		for {
			ins, err := s.tunnel.ReceiveInstruction()
			if err != nil {
				s.logger.Warnf("Disconnecting from server. Reason: " + err.Error())
				s.Terminate()
				break
			}
			if ins.Opcode == "blob" {
				s.logger.Tracef("S> %s", ins)
			} else {
				s.logger.Debugf("S> %s", ins)
			}
			if ins.Opcode == "nop" {
				continue
			}
			if ins.Opcode == "ready" {
				s.State = SessionActive
				s.Id = ins.Args[0]
				s.logger.Infof("Handshake successful. Got connection ID %s", s.Id)
				s.startKeepAlive()
				continue
			}
			if s.State == SessionHandshake {
				s.logger.Infof("Handshake started at %s", time.Now().Format(time.RFC3339))
				s.handShake(ins)
				continue
			}
			if s.State == SessionActive {
				s.In <- ins
				continue
			}
			s.logger.Warnf("Received out of order instruction: %s", ins)
		}
	}()
}

func (s *session) handShake(argsIns *protocol.Instruction) {
	options := []*protocol.Instruction{
		protocol.NewInstruction("size", "1024", "768", "96"),
		protocol.NewInstruction("audio", ""),
		protocol.NewInstruction("video", ""),
		protocol.NewInstruction("image", ""),
	}

	err := s.Send(options...)
	if err != nil {
		s.logger.Errorf("Failed handshake: %s", err)
		s.Terminate()
	}

	connectValues := make([]string, len(argsIns.Args))
	for i, argName := range argsIns.Args {
		connectValues[i] = s.config[argName]
	}

	err = s.Send(protocol.NewInstruction("connect", connectValues...))
	if err != nil {
		s.logger.Errorf("Failed handshake when sending 'connect': %s", err)
		s.Terminate()
	}
}
