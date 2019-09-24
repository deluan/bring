package gocamole

import (
	"bufio"
	"io"
)

// InstructionIO ...
type InstructionIO struct {
	conn   io.ReadWriteCloser
	input  *bufio.Reader
	output *bufio.Writer
}

// NewInstructionIO ...
func NewInstructionIO(sock io.ReadWriteCloser) *InstructionIO {
	conn := sock
	return &InstructionIO{
		conn:   conn,
		input:  bufio.NewReaderSize(conn, MaxInstructionLength),
		output: bufio.NewWriter(conn),
	}
}

// Close closes the InstructionIO
func (io *InstructionIO) Close() error {
	return io.conn.Close()
}

// ReadRaw reads raw data from io input
func (io *InstructionIO) ReadRaw() ([]byte, error) {
	return io.input.ReadBytes(byte(';'))
}

// Read reads and parses the instruction from io input
func (io *InstructionIO) Read() (*Instruction, error) {
	raw, err := io.ReadRaw()
	if err != nil {
		return nil, err
	}
	return ParseInstruction(raw)
}

// WriteRaw writes raw buffer into io output
func (io *InstructionIO) WriteRaw(buf []byte) (n int, err error) {
	n, err = io.output.Write(buf)
	if err != nil {
		return
	}
	err = io.output.Flush()
	return
}

// Write writes and decodes an instruction to io output
func (io *InstructionIO) Write(ins *Instruction) (int, error) {
	return io.WriteRaw([]byte(ins.String()))
}
