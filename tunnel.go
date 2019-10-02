package bring

import (
	"net"
	"sync"
	"time"
)

type TunnelState int

// All possible tunnel states.
const (
	TunnelClosed TunnelState = iota
	TunnelOpen
)

type Tunnel interface {
	// Connect to the tunnel with the given optional data. This data is
	// typically used for authentication. The format of data accepted is
	// up to the tunnel implementation.
	Connect(data string) error

	// Disconnect from the tunnel.
	Disconnect()

	// Send the given message through the tunnel to the service on the other
	// side. All messages are guaranteed to be received in the order sent.
	SendInstruction(ins ...*Instruction) error

	// Receive a instruction from the service on the other side. This is a
	// blocking call
	ReceiveInstruction() (*Instruction, error)
}

const connectionTimeout = 5 * time.Second

// Simple tunnel implementation over TCP sockets. This allows to connect directly to guacd,
// without the need of any middleware
type InetSocketTunnel struct {
	address    string
	socket     net.Conn
	state      TunnelState
	io         *InstructionIO
	writeMutex sync.Mutex
}

func NewInetSocketTunnel(address string) (*InetSocketTunnel, error) {
	t := &InetSocketTunnel{address: address}

	return t, nil
}

func (t *InetSocketTunnel) SendInstruction(ins ...*Instruction) error {
	if t.state != TunnelOpen {
		return ErrNotConnected
	}

	if len(ins) == 0 {
		return nil
	}

	t.writeMutex.Lock()
	defer t.writeMutex.Unlock()

	var err error
	for _, in := range ins {
		_, err := t.io.Write(in)
		if err != nil {
			break
		}
	}
	return err
}

// TODO Implement timeout
func (t *InetSocketTunnel) ReceiveInstruction() (*Instruction, error) {
	if t.state != TunnelOpen {
		return nil, ErrNotConnected
	}

	return t.io.Read()
}

func (t *InetSocketTunnel) Connect(data string) error {
	sock, err := net.DialTimeout("tcp", t.address, connectionTimeout)
	if err != nil {
		return err
	}

	t.socket = sock
	t.state = TunnelOpen
	t.io = NewInstructionIO(t.socket)
	return nil
}

func (t *InetSocketTunnel) Disconnect() {
	t.closeTunnel()
}

func (t *InetSocketTunnel) closeTunnel() {
	if t.state == TunnelClosed {
		return
	}

	t.state = TunnelClosed
	_ = t.io.Close()
}
