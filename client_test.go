package bring

import (
	"image"
	"math"
	"strconv"

	"github.com/deluan/bring/protocol"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var s *session
	var c *Client
	var t *mockTunnel

	BeforeEach(func() {
		t = &mockTunnel{}
		l := &DefaultLogger{Quiet: true}
		s = &session{
			In:       make(chan *protocol.Instruction, 100),
			State:    SessionActive,
			done:     make(chan bool),
			logger:   l,
			tunnel:   t,
			protocol: "vnc",
		}
		c = &Client{
			session: s,
			display: newDisplay(l),
			streams: newStreams(),
			logger:  l,
		}
	})

	Context("Active session", func() {
		It("exposes the session state", func() {
			s.State = SessionHandshake
			Expect(c.State()).To(Equal(SessionHandshake))
			s.State = SessionClosed
			Expect(c.State()).To(Equal(SessionClosed))
		})

		It("sends the mouse position to the remote server", func() {
			err := c.SendMouse(image.Pt(10, 20))
			Expect(err).To(BeNil())

			Expect(t.sent[0].Opcode).To(Equal("mouse"))
			Expect(t.sent[0].Args).To(Equal([]string{"10", "20", "0"}))
		})

		It("sends the position to the remote server", func() {
			err := c.SendMouse(image.Pt(10, 20), MouseLeft, MouseDown)
			Expect(err).To(BeNil())

			Expect(t.sent[0].Opcode).To(Equal("mouse"))
			Expect(t.sent[0].Args[2]).To(Equal(strconv.Itoa(1 + 16)))
		})

		It("sends a single keyscan to the remote server", func() {
			err := c.SendKey(KeyBackspace, false)
			keyBackspace := keySyms[KeyBackspace]

			Expect(err).To(BeNil())
			Expect(t.sent).To(HaveLen(len(keyBackspace)))
			Expect(t.sent[0].Opcode).To(Equal("key"))
			Expect(t.sent[0].Args).To(Equal([]string{strconv.Itoa(keyBackspace[0]), "0"}))
		})

		It("sends a key with multiple keyscans to the remote server", func() {
			err := c.SendKey(KeyRightShift, true)
			KeyRightShift := keySyms[KeyRightShift]

			Expect(err).To(BeNil())
			Expect(t.sent).To(HaveLen(len(KeyRightShift)))
			Expect(t.sent[0].Opcode).To(Equal("key"))
			Expect(t.sent[0].Args).To(Equal([]string{strconv.Itoa(KeyRightShift[0]), "1"}))
			Expect(t.sent[1].Opcode).To(Equal("key"))
			Expect(t.sent[1].Args).To(Equal([]string{strconv.Itoa(KeyRightShift[1]), "1"}))
		})

		It("returns an ErrInvalidKeyCode when receiving an invalid keycode", func() {
			err := c.SendKey(KeyCode(math.MaxInt32), true)

			Expect(err).To(Equal(ErrInvalidKeyCode))
			Expect(t.sent).To(BeEmpty())
		})

		It("sends a text as individual keystrokes", func() {
			err := c.SendText("bring")
			Expect(err).To(BeNil())
			Expect(t.sent).To(HaveLen(10))
			for i, c := range "bring" {
				Expect(t.sent[i*2]).To(Equal(protocol.NewInstruction("key", toAscii(c), "1")))
				Expect(t.sent[i*2+1]).To(Equal(protocol.NewInstruction("key", toAscii(c), "0")))
			}
		})
	})

	Context("Session is disconnected", func() {
		BeforeEach(func() {
			s.State = SessionClosed
		})

		It("does not send anything", func() {
			err := c.SendKey(KeyEnter, true)
			Expect(err).To(Equal(ErrNotConnected))

			err = c.SendText("abc")
			Expect(err).To(Equal(ErrNotConnected))

			err = c.SendMouse(image.Pt(0, 0), MouseRight)
			Expect(err).To(Equal(ErrNotConnected))
		})
	})
})

func toAscii(c int32) string {
	return strconv.Itoa(int(c))
}

type mockTunnel struct {
	protocol.Tunnel
	sent []*protocol.Instruction
}

func (mt *mockTunnel) SendInstruction(ins ...*protocol.Instruction) error {
	mt.sent = append(mt.sent, ins...)
	return nil
}
