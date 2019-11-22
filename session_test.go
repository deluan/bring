package bring

import (
	"time"

	"github.com/deluan/bring/protocol"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {
	var (
		server *fakeServer
		s      *session
	)
	BeforeEach(func() {
		server = &fakeServer{
			replies: map[string]string{
				"select":  "4.args,8.hostname,4.port,8.password;",
				"connect": "5.ready,21.$unique-connection-id;",
			},
		}
		addr := server.start()
		s, _ = newSession(addr, "rdp", map[string]string{
			"hostname": "host1",
			"port":     "port1",
			"password": "password123",
		}, &DefaultLogger{Quiet: true})

		Eventually(func() SessionState {
			return s.State
		}, 3*time.Second, 100*time.Millisecond).Should(Equal(SessionActive))
	})

	It("should handshake with server", func() {
		err := s.Send(protocol.NewInstruction(disconnectOpcode))
		Expect(err).To(BeNil())

		Expect(server.opcodesReceived).To(Equal([]string{"select", "size", "audio", "video", "image", "connect"}))
		Expect(server.messagesReceived[0]).To(Equal("6.select,3.rdp;"))
		Expect(server.messagesReceived[len(server.messagesReceived)-1]).To(Equal("7.connect,5.host1,5.port1,11.password123;"))
		Expect(s.Id).To(Equal("$unique-connection-id"))
	})
})
