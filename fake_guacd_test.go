package bring

import (
	"net"
	"testing"
	"time"

	"github.com/deluan/bring/protocol"
)

const disconnectOpcode = "disconnect"

type fakeServer struct {
	replies          map[string]string
	messagesReceived []string
	opcodesReceived  []string
}

func (s *fakeServer) start() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				panic(err)
			}
			s.handleRequest(conn)
		}
	}()
	return ln.Addr().String()
}

func (s *fakeServer) handleRequest(conn net.Conn) {
	defer conn.Close()
	io := protocol.NewInstructionIO(conn)
	for {
		recv, err := io.Read()
		if err != nil {
			return
		}
		opcode := recv.Opcode

		_, err = io.WriteRaw([]byte(s.replies[opcode]))
		if err != nil {
			panic(err)
		}
		if opcode == disconnectOpcode {
			break
		}
		s.messagesReceived = append(s.messagesReceived, recv.String())
		s.opcodesReceived = append(s.opcodesReceived, opcode)
	}
	_, err := io.Write(protocol.NewInstruction(disconnectOpcode))
	if err != nil {
		panic(err)
	}
}

func waitForHandshake(t *testing.T, s *session) {
	// Wait the end of the Handshake for 2 seconds
	for i := 0; i < 20; i++ {
		if s.State == SessionActive {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("TimeOut waiting for handshake. Session= %+v", s)
}

func disconnectFromFakeServer(t *testing.T, s *session) {
	err := s.Send(protocol.NewInstruction(disconnectOpcode))
	if err != nil {
		t.Fatalf("Error trying to disconnect from fake server: %s", err)
	}
}
