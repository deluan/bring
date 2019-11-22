package bring

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Streams", func() {
	var ss streams
	BeforeEach(func() {
		ss = newStreams()
	})

	It("creates a new stream when it does not exist", func() {
		s := ss.get(1)

		Expect(s.buffer).ToNot(BeNil())
		Expect(ss[1]).To(Equal(s))
	})

	Context("Given an new empty stream", func() {
		var s *stream
		BeforeEach(func() {
			s = ss.get(2)
		})

		It("appends data to the stream", func() {
			err := ss.append(2, "test data")

			Expect(err).To(BeNil())
			Expect(s.buffer.String()).To(Equal("test data"))
		})

		It("decodes bas64 images", func() {
			err := ss.append(2, "iVBORw0KGgoAAAANSUhEUgAAAAEAAAAPAgMAAABYcU1qAAAACVBMVEX8/Pzc3Nzr6+uSJe5dAAAAEUlEQVQImWNgAAIHhgYGrAAAEd4AwbcvDeEAAAAASUVORK5CYII=")
			Expect(err).To(BeNil())

			img, err := s.image()
			Expect(err).To(BeNil())
			Expect(img.Bounds()).To(Equal(image.Rect(0, 0, 1, 15)))
		})

		It("executes the endFunc when I call end", func() {
			called := false
			s.onEnd = func(sp *stream) {
				called = true
				Expect(sp).To(Equal(s))
			}
			ss.end(2)
			Expect(called).To(BeTrue())
		})

		It("removes it from the streams map", func() {
			previousSize := len(ss)
			ss.delete(2)

			Expect(ss[2]).To(BeNil())
			Expect(ss).To(HaveLen(previousSize - 1))
		})

	})
})
