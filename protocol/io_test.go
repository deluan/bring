package protocol

import (
	"io/ioutil"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InstructionIO", func() {
	var server, client net.Conn
	var io *InstructionIO

	BeforeEach(func() {
		server, client = net.Pipe()
		io = NewInstructionIO(client)
	})
	AfterEach(func() {
		client.Close()
		server.Close()
	})

	It("formats and sends the instruction", func() {
		ins := NewInstruction("hello", "ग्वाकोमोल")
		go func() {
			io.Write(ins)
			io.Close()
		}()

		Expect(ioutil.ReadAll(server)).To(Equal([]byte("5.hello,9.ग्वाकोमोल;")))
	})

	It("reads and parses received instruction", func() {
		msg1 := "5.hello,9.ग्वाकोमोल;"
		msg2 := "5.empty,0.;"
		go func() {
			server.Write([]byte(msg1))
			server.Write([]byte(msg2))
			server.Close()
		}()
		ins1, err1 := io.Read()
		Expect(err1).To(BeNil())
		Expect(ins1.String()).To(Equal(msg1))

		ins2, err2 := io.Read()
		Expect(err2).To(BeNil())
		Expect(ins2.String()).To(Equal(msg2))
	})

})
