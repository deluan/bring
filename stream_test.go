package bring

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStreams(t *testing.T) {
	Convey("Given a Streams map", t, func() {
		ss := newStreams()

		Convey("When I get an non-existent stream", func() {
			s := ss.get(1)

			Convey("It returns an newly initialized stream", func() {
				So(s.buffer, ShouldNotBeNil)
				So(ss[1], ShouldEqual, s)
			})
		})

		Convey("Given a new empty stream", func() {
			s := ss.get(2)

			Convey("When I call append on it", func() {
				_ = ss.append(2, "test data")

				Convey("It adds the data to the stream", func() {
					So(s.buffer.String(), ShouldEqual, "test data")
				})
			})

			Convey("When I have an endFunc assigned to it", func() {
				called := false
				s.onEnd = func(sp *stream) {
					called = true
					So(sp, ShouldEqual, s)
				}

				Convey("It executes the endFunc when I call end", func() {
					ss.end(2)
					So(called, ShouldBeTrue)
				})
			})

			Convey("When I call delete on it", func() {
				beforeSize := len(ss)
				ss.delete(2)

				Convey("It removes it from the streams map", func() {
					So(ss[2], ShouldBeNil)
					So(len(ss), ShouldEqual, beforeSize-1)
				})
			})
		})
	})
}
