package bring

import (
	"github.com/deluan/bring/protocol"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {
	server := &fakeServer{
		replies: map[string]string{
			"select":  "4.args,8.hostname,4.port,8.password;",
			"connect": "5.ready,21.$unique-connection-id;",
		},
	}
	addr := server.start()
	s, _ := newSession(addr, "rdp", map[string]string{
		"hostname": "host1",
		"port":     "port1",
		"password": "password123",
	}, &DefaultLogger{Quiet: true})

	It("should handshake with server", func() {
		Eventually(s.State).Should(Equal(SessionActive))

		err := s.Send(protocol.NewInstruction(disconnectOpcode))
		Expect(err).To(BeNil())

		Expect(server.opcodesReceived).To(Equal([]string{"select", "size", "audio", "video", "image", "connect"}))
		Expect(server.messagesReceived[0]).To(Equal("6.select,3.rdp;"))
		Expect(server.messagesReceived[len(server.messagesReceived)-1]).To(Equal("7.connect,5.host1,5.port1,11.password123;"))
		Expect(s.Id).To(Equal("$unique-connection-id"))
	})
})
