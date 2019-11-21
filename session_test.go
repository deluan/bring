package bring

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSession(t *testing.T) {
	Convey("When creating a new Session", t, func() {
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
		}, &DefaultLogger{})

		waitForHandshake(t, s)
		disconnectFromFakeServer(t, s)

		Convey("It executes the handshake", func() {
			So(server.opcodesReceived, ShouldResemble, []string{"select", "size", "audio", "video", "image", "connect"})
			So(server.messagesReceived[0], ShouldEqual, "6.select,3.rdp;")
			So(server.messagesReceived[len(server.messagesReceived)-1], ShouldEqual, "7.connect,5.host1,5.port1,11.password123;")
			So(s.Id, ShouldEqual, "$unique-connection-id")
		})
	})
}
