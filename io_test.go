package bring

import (
	"io/ioutil"
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInstructionIO(t *testing.T) {
	Convey("Describe InstructionIO", t, func() {
		server, client := net.Pipe()
		defer client.Close()
		defer server.Close()
		io := NewInstructionIO(client)

		Convey("When we send an instruction", func() {
			ins := NewInstruction("hello", "ग्वाकोमोल")
			go func() {
				io.Write(ins)
				io.Close()
			}()

			Convey("The server receives the formatted instruction", func() {
				buf, err := ioutil.ReadAll(server)
				server.Close()
				So(err, ShouldBeNil)
				So(string(buf), ShouldEqual, "5.hello,9.ग्वाकोमोल;")
			})
		})

		Convey("When the server sends instructions", func() {
			msg1 := "5.hello,9.ग्वाकोमोल;"
			msg2 := "5.empty,0.;"
			go func() {
				server.Write([]byte(msg1))
				server.Write([]byte(msg2))
				server.Close()
			}()

			Convey("Read() returns the parsed instructions", func() {
				ins1, err1 := io.Read()
				So(err1, ShouldBeNil)
				So(ins1.String(), ShouldEqual, msg1)

				ins2, err2 := io.Read()
				So(err2, ShouldBeNil)
				So(ins2.String(), ShouldEqual, msg2)
			})
		})
	})
}
