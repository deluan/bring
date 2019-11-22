package bring

import (
	"image/draw"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Layers", func() {
	var layers layers
	BeforeEach(func() {
		layers = newLayers()
		layers.getDefault().Resize(1024, 768)
	})

	It("returns the default layer", func() {
		Expect(layers.getDefault()).To(Equal(layers[0]))
	})

	It("properly creates new layers", func() {
		l := layers.get(30)
		Expect(layers[30]).ToNot(BeNil(), "the new layer should be saved to the layers map")
		Expect(l.autosize).To(BeFalse(), "the new layer should be not be autosizeable")
		Expect(l.visible).To(BeTrue(), "the new layer should be visible")
		Expect(l.width).To(Equal(layers.getDefault().width), "the new layer should have the default width")
		Expect(l.height).To(Equal(layers.getDefault().height), "the new layer should have the default height")
	})

	It("properly creates new buffers (invisible layers)", func() {
		l := layers.get(-3)
		Expect(layers[-3]).ToNot(BeNil())
		Expect(l.autosize).To(BeTrue(), "the new layer should be not be autosizeable")
		Expect(l.width).To(BeZero(), "the new layer should have zero width")
		Expect(l.height).To(BeZero(), "the new layer should have zero height")
		Expect(l.visible).To(BeFalse(), "the new layer should be invisible")
	})

	It("deletes the layer from the layers map", func() {
		previousSize := len(layers)
		layers.get(40)
		Expect(layers).To(HaveLen(previousSize + 1))
		layers.delete(40)
		Expect(layers[40]).To(BeNil())
		Expect(layers).To(HaveLen(previousSize))
	})

	It("must not delete the default layer", func() {
		l := layers.getDefault()
		layers.delete(0)
		Expect(layers[0]).To(Equal(l))
	})
})

var _ = Describe("Layer", func() {
	var layers layers
	var src, dst *layer

	BeforeEach(func() {
		layers = newLayers()
		layers.getDefault().Resize(1024, 768)
		dst = newBuffer()
		src = newBuffer()
		src.Resize(100, 100)
	})

	Describe("fitRect", func() {
		It("grows the buffer to fit a rectangle", func() {
			dst.fitRect(10, 10, 20, 30)
			Expect(dst.width).To(Equal(30))
			Expect(dst.height).To(Equal(40))
		})

		It("does not change the size of the buffer if the rectangle is smaller", func() {
			dst.Resize(100, 100)
			dst.fitRect(10, 10, 20, 30)
			Expect(dst.width).To(Equal(100))
			Expect(dst.height).To(Equal(100))
		})

		It("When I call fitRect with a rectangle that overflows the buffer", func() {
			dst.Resize(100, 10)
			dst.fitRect(10, 10, 20, 30)
			Expect(dst.width).To(Equal(100))
			Expect(dst.height).To(Equal(40))
		})
	})

	Describe("Copy", func() {
		It("grows the buffer to fit a larger source", func() {
			dst.Copy(src, 0, 0, 100, 100, 0, 0, draw.Src)
			Expect(dst.width).To(Equal(100))
			Expect(dst.height).To(Equal(100))
		})

		It("clips the source canvas if copying a larger rectangle from source", func() {
			dst.Copy(src, 0, 0, 200, 200, 0, 0, draw.Src)
			Expect(dst.width).To(Equal(100))
			Expect(dst.height).To(Equal(100))
		})

		It("does not copy anything if rectangle is outside of the src canvas", func() {
			dst.Copy(src, 120, 120, 10, 10, 0, 0, draw.Src)
			Expect(dst.width).To(BeZero())
			Expect(dst.height).To(BeZero())
		})
	})
})
