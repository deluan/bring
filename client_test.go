package bring

import (
	"image"
	"math"
	"strconv"
	"testing"

	"github.com/deluan/bring/protocol"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {
	Convey("Given a Client with an active session", t, func() {
		t := &mockTunnel{}
		s := &Session{
			In:       make(chan *protocol.Instruction, 100),
			State:    SessionActive,
			done:     make(chan bool),
			logger:   &DefaultLogger{Quiet: true},
			tunnel:   t,
			protocol: "vnc",
		}
		c, _ := NewClient(s, &DefaultLogger{Quiet: true})

		Convey("It exposes the session state", func() {
			s.State = SessionHandshake
			So(c.State(), ShouldEqual, SessionHandshake)
			s.State = SessionClosed
			So(c.State(), ShouldEqual, SessionClosed)
		})

		Convey("When it receives a mouse position", func() {
			err := c.SendMouse(image.Pt(10, 20))

			Convey("It sends the position to the tunnel", func() {
				So(err, ShouldBeNil)
				So(t.sent[0].Opcode, ShouldEqual, "mouse")
				So(t.sent[0].Args, ShouldResemble, []string{"10", "20", "0"})
			})
		})

		Convey("When it receives mouse buttons", func() {
			err := c.SendMouse(image.Pt(10, 20), MouseLeft, MouseDown)

			Convey("It sends the position to the tunnel", func() {
				So(err, ShouldBeNil)
				So(t.sent[0].Opcode, ShouldEqual, "mouse")
				So(t.sent[0].Args[2], ShouldEqual, strconv.Itoa(1+16))
			})
		})

		Convey("When it receives a key with single keyscan", func() {
			err := c.SendKey(KeyBackspace, false)
			keyBackspace := keySyms[KeyBackspace]

			Convey("It sends the keycode", func() {
				So(err, ShouldBeNil)
				So(t.sent, ShouldHaveLength, 1)
				So(t.sent[0].Opcode, ShouldEqual, "key")
				So(t.sent[0].Args, ShouldResemble, []string{strconv.Itoa(keyBackspace[0]), "0"})
			})
		})

		Convey("When it receives a key with multiple keyscans", func() {
			err := c.SendKey(KeyRightShift, true)
			keyRightShift := keySyms[KeyRightShift]

			Convey("It sends the keycode", func() {
				So(err, ShouldBeNil)
				So(t.sent, ShouldHaveLength, len(keyRightShift))
				So(t.sent[0].Opcode, ShouldEqual, "key")
				So(t.sent[0].Args, ShouldResemble, []string{strconv.Itoa(keyRightShift[0]), "1"})
				So(t.sent[1].Opcode, ShouldEqual, "key")
				So(t.sent[1].Args, ShouldResemble, []string{strconv.Itoa(keyRightShift[1]), "1"})
			})
		})

		Convey("When it receives an invalid KeyCode", func() {
			err := c.SendKey(KeyCode(math.MaxInt32), true)

			Convey("It returns ErrInvalidKeyCode", func() {
				So(err, ShouldResemble, ErrInvalidKeyCode)
				So(t.sent, ShouldHaveLength, 0)
			})
		})

		Convey("When it receives a text to be sent", func() {
			err := c.SendText("bring")

			Convey("It sends all keycodes", func() {
				So(err, ShouldBeNil)
				So(t.sent, ShouldHaveLength, 10)
				So(t.sent[0], ShouldResemble, protocol.NewInstruction("key", toAscii("b"), "1"))
				So(t.sent[1], ShouldResemble, protocol.NewInstruction("key", toAscii("b"), "0"))
				So(t.sent[2], ShouldResemble, protocol.NewInstruction("key", toAscii("r"), "1"))
				So(t.sent[3], ShouldResemble, protocol.NewInstruction("key", toAscii("r"), "0"))
			})
		})

		Convey("When it is disconnected", func() {
			s.State = SessionClosed

			Convey("It does not send anything", func() {
				err := c.SendKey(KeyEnter, true)
				So(err, ShouldResemble, ErrNotConnected)

				err = c.SendText("abc")
				So(err, ShouldResemble, ErrNotConnected)

				err = c.SendMouse(image.Pt(0, 0), MouseRight)
				So(err, ShouldResemble, ErrNotConnected)

				So(t.sent, ShouldBeEmpty)
			})
		})
	})
}

func toAscii(c string) string {
	return strconv.Itoa(int(c[0]))
}

type mockTunnel struct {
	protocol.Tunnel
	sent []*protocol.Instruction
}

func (mt *mockTunnel) SendInstruction(ins ...*protocol.Instruction) error {
	mt.sent = append(mt.sent, ins...)
	return nil
}
