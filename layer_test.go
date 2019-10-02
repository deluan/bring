package bring

import (
	"image/draw"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLayers(t *testing.T) {
	Convey("Given a Layers map with a default layer", t, func() {
		layers := newLayers()
		layers.getDefault().Resize(1024, 768)

		Convey("getDefault returns the layer 0", func() {
			So(layers.getDefault(), ShouldEqual, layers[0])
		})

		Convey("When I request a new visible layer", func() {
			l := layers.get(30)
			Convey("It is created as a non-autosizeable layer", func() {
				So(l.autosize, ShouldBeFalse)
			})

			Convey("It returns a layer with the same dimensions as the default", func() {
				So(layers[30], ShouldNotBeNil)
				So(l.width, ShouldEqual, layers.getDefault().width)
				So(l.height, ShouldEqual, layers.getDefault().height)
				So(l.visible, ShouldBeTrue)
			})
		})

		Convey("When I request a new invisible layer (aka buffer)", func() {
			l := layers.get(-3)
			Convey("It returns a layer with dimensions 0x0", func() {
				So(layers[-3], ShouldNotBeNil)
				So(l.width, ShouldEqual, 0)
				So(l.height, ShouldEqual, 0)
				So(l.visible, ShouldBeFalse)
			})
		})

		Convey("When I call delete on an existing layer", func() {
			layers.get(40)
			beforeSize := len(layers)
			layers.delete(40)
			Convey("It removes it from the layers map", func() {
				So(layers[40], ShouldBeNil)
				So(len(layers), ShouldEqual, beforeSize-1)
			})
		})

		Convey("When I call delete on the default layer", func() {
			l := layers.getDefault()
			layers.delete(0)
			Convey("It does not removes", func() {
				So(layers[0], ShouldEqual, l)
			})
		})

		Convey("When I have a empty buffer", func() {
			dst := newBuffer()
			src := newBuffer()
			src.Resize(100, 100)
			Convey("When I copy another buffer to it", func() {
				dst.Copy(src, 0, 0, 100, 100, 0, 0, draw.Src)

				Convey("It grows the buffer to fit the source", func() {
					So(dst.width, ShouldEqual, 100)
					So(dst.height, ShouldEqual, 100)
				})
			})

			Convey("When I try to copy a larger rectangle from the source buffer", func() {
				dst.Copy(src, 0, 0, 200, 200, 0, 0, draw.Src)

				Convey("It clips the source canvas", func() {
					So(dst.width, ShouldEqual, 100)
					So(dst.height, ShouldEqual, 100)
				})
			})

			Convey("When I try to copy a rectangle from outside of the src canvas", func() {
				dst.Copy(src, 120, 120, 10, 10, 0, 0, draw.Src)

				Convey("It does not copy anything", func() {
					So(dst.width, ShouldEqual, 0)
					So(dst.height, ShouldEqual, 0)
				})
			})
		})

	})

	Convey("Given an empty buffer", t, func() {
		l := newBuffer()

		Convey("It is created as a autosizeable layer", func() {
			So(l.autosize, ShouldBeTrue)
		})

		Convey("It is created with a size of 0x0", func() {
			So(l.width, ShouldEqual, 0)
			So(l.height, ShouldEqual, 0)
		})

		Convey("When I call fitRect with a rectangle", func() {
			l.fitRect(10, 10, 20, 30)

			Convey("It grows the layer to fit the rectangle", func() {
				So(l.width, ShouldEqual, 30)
				So(l.height, ShouldEqual, 40)
			})

		})
	})

	Convey("Given a big buffer", t, func() {
		l := newBuffer()
		l.Resize(100, 100)

		Convey("When I call fitRect with a small rectangle", func() {
			l.fitRect(10, 10, 20, 30)

			Convey("It does not affect the size of the buffer", func() {
				So(l.width, ShouldEqual, 100)
				So(l.height, ShouldEqual, 100)
			})

		})
	})

	Convey("Given a buffer", t, func() {
		l := newBuffer()
		l.Resize(100, 10)

		Convey("When I call fitRect with a rectangle that overflows the buffer", func() {
			l.fitRect(10, 10, 20, 30)

			Convey("It resizes the buffer", func() {
				So(l.width, ShouldEqual, 100)
				So(l.height, ShouldEqual, 40)
			})

		})
	})

}
