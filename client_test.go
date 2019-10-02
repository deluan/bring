package bring

import (
	"image"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {
	Convey("Given a Client with an active session", t, func() {
		t := &mockTunnel{}
		s := &Session{
			In:       make(chan *Instruction, 100),
			State:    SessionActive,
			done:     make(chan bool),
			logger:   &DiscardLogger{},
			tunnel:   t,
			protocol: "vnc",
		}
		c, _ := NewClient(s, &DiscardLogger{})

		Convey("It exposes the session state", func() {
			s.State = SessionHandshake
			So(c.State(), ShouldEqual, SessionHandshake)
			s.State = SessionClosed
			So(c.State(), ShouldEqual, SessionClosed)
		})

		Convey("When it receives a mouse position", func() {
			c.SendMouse(image.Pt(10, 20))

			Convey("It sends the position to the tunnel", func() {
				So(t.sent[0].opcode, ShouldEqual, "mouse")
				So(t.sent[0].args, ShouldResemble, []string{"10", "20", "0"})
			})
		})

		Convey("When it receives mouse buttons", func() {
			c.SendMouse(image.Pt(10, 20), MouseLeft, MouseDown)

			Convey("It sends the position to the tunnel", func() {
				So(t.sent[0].opcode, ShouldEqual, "mouse")
				So(t.sent[0].args[2], ShouldEqual, strconv.Itoa(1+16))
			})
		})

		Convey("When it receives a key with single keyscan", func() {
			c.SendKey(KeyBackspace, false)

			Convey("It sends the keycode", func() {
				So(t.sent, ShouldHaveLength, 1)
				So(t.sent[0].opcode, ShouldEqual, "key")
				So(t.sent[0].args, ShouldResemble, []string{strconv.Itoa(KeyBackspace[0]), "0"})
			})
		})

		Convey("When it receives a key with multiple keyscans", func() {
			c.SendKey(KeyRightShift, true)

			Convey("It sends the keycode", func() {
				So(t.sent, ShouldHaveLength, len(KeyRightShift))
				So(t.sent[0].opcode, ShouldEqual, "key")
				So(t.sent[0].args, ShouldResemble, []string{strconv.Itoa(KeyRightShift[0]), "1"})
				So(t.sent[1].opcode, ShouldEqual, "key")
				So(t.sent[1].args, ShouldResemble, []string{strconv.Itoa(KeyRightShift[1]), "1"})
			})
		})

		Convey("When it receives a text to be sent", func() {
			c.SendText("bring")

			Convey("It sends all keycodes", func() {
				So(t.sent, ShouldHaveLength, 10)
				So(t.sent[0], ShouldResemble, NewInstruction("key", to_i("b"), "1"))
				So(t.sent[1], ShouldResemble, NewInstruction("key", to_i("b"), "0"))
				So(t.sent[2], ShouldResemble, NewInstruction("key", to_i("r"), "1"))
				So(t.sent[3], ShouldResemble, NewInstruction("key", to_i("r"), "0"))
			})
		})

		Convey("When it is disconnected", func() {
			s.State = SessionClosed

			Convey("It does not send anything", func() {
				c.SendKey(KeyEnter, true)
				c.SendText("abc")
				c.SendMouse(image.Pt(0, 0), MouseRight)

				So(t.sent, ShouldBeEmpty)
			})
		})
	})
}

func to_i(c string) string {
	return strconv.Itoa(int(c[0]))
}

type mockTunnel struct {
	Tunnel
	sent []*Instruction
}

func (mt *mockTunnel) SendInstruction(ins ...*Instruction) error {
	mt.sent = append(mt.sent, ins...)
	return nil
}
