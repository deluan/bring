package bring

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInstruction(t *testing.T) {
	Convey("Describe String()", t, func() {
		Convey("When I create a new instruction without args", func() {
			i := NewInstruction("nop")
			Convey("Then I get a properly formatted instruction", func() {
				So(i.String(), ShouldEqual, "3.nop;")
			})
		})
		Convey("When I create a new instruction with UTF args", func() {
			i := NewInstruction("hello", "世界", "yes")
			Convey("Then I get a properly formatted instruction", func() {
				So(i.String(), ShouldEqual, "5.hello,2.世界,3.yes;")
			})
		})
		Convey("When I create a new instruction with an empty arg", func() {
			i := NewInstruction("empty", "")
			Convey("Then I get a properly formatted instruction", func() {
				So(i.String(), ShouldEqual, "5.empty,0.;")
			})
		})
	})

	Convey("Describe ParseInstruction()", t, func() {
		Convey("Given a raw string with an instruction without args", func() {
			ri := "3.nop;"
			Convey("It parses it to a Instruction struct", func() {
				i, _ := ParseInstruction([]byte(ri))
				So(i.opcode, ShouldEqual, "nop")
				So(i.args, ShouldHaveLength, 0)
			})
		})
		Convey("Given a raw string with an instruction with UTF args", func() {
			ri := "5.hello,2.世界,3.yes;"
			Convey("It parses it to a Instruction struct", func() {
				i, _ := ParseInstruction([]byte(ri))
				So(i.opcode, ShouldEqual, "hello")
				So(i.args, ShouldResemble, []string{"世界", "yes"})
			})
		})
		Convey("Given a raw string with an instruction with an empty arg", func() {
			ri := "4.test,0.;"
			Convey("It parses it to a Instruction struct", func() {
				i, _ := ParseInstruction([]byte(ri))
				So(i.opcode, ShouldEqual, "test")
				So(i.args, ShouldResemble, []string{""})
			})
		})
	})
}

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
