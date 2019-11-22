package protocol

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Instruction", func() {
	Describe("String()", func() {
		It("formats instructions without args", func() {
			i := NewInstruction("nop")
			Expect(i.String()).To(Equal("3.nop;"))
		})

		It("formats instructions with UTF args", func() {
			i := NewInstruction("hello", "世界", "yes")
			Expect(i.String()).To(Equal("5.hello,2.世界,3.yes;"))
		})

		It("formats instructions with an empty arg", func() {
			i := NewInstruction("empty", "")
			Expect(i.String()).To(Equal("5.empty,0.;"))
		})
	})

	Describe("ParseInstruction()", func() {
		It("parses instructions without args", func() {
			ri := "3.nop;"
			i, err := ParseInstruction([]byte(ri))
			Expect(err).To(BeNil())
			Expect(i.Opcode).To(Equal("nop"))
			Expect(i.Args).To(HaveLen(0))
		})

		It("parses instructions with UTF args", func() {
			ri := "5.hello,2.世界,3.yes;"
			i, err := ParseInstruction([]byte(ri))
			Expect(err).To(BeNil())
			Expect(i.Opcode).To(Equal("hello"))
			Expect(i.Args).To(Equal([]string{"世界", "yes"}))
		})

		It("parses instructions with an empty arg", func() {
			ri := "4.test,0.;"
			i, err := ParseInstruction([]byte(ri))
			Expect(err).To(BeNil())
			Expect(i.Opcode).To(Equal("test"))
			Expect(i.Args).To(Equal([]string{""}))
		})
	})

})

func BenchmarkInstruction_String(b *testing.B) {
	i := NewInstruction("test", "p1", "p2", "p3", "p4")

	b.ReportAllocs()
	for x := 0; x < b.N; x++ {
		_ = i.String()
	}
}

func BenchmarkInstruction_ParseInstruction(b *testing.B) {
	ri := "4.test,3.one,3.two,5.three,4.four;"

	b.ReportAllocs()
	for x := 0; x < b.N; x++ {
		_, _ = ParseInstruction([]byte(ri))
	}
}
