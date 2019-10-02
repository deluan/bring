package protocol

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

const maxInstructionLength = 8192

// Possible errors returned by the ParseInstruction function
var (
	ErrInstructionMissDot   = errors.New("instruction without dot")
	ErrInstructionMissComma = errors.New("instruction without comma")
	ErrInstructionBadDigit  = errors.New("instruction with bad digit")
	ErrInstructionBadRune   = errors.New("instruction with bad rune")
)

type Instruction struct {
	Opcode string
	Args   []string
}

func NewInstruction(opcode string, args ...string) *Instruction {
	return &Instruction{
		Opcode: opcode,
		Args:   args,
	}
}

// ParseInstruction parses an instruction: 1.a,2.bc,3.def,10.abcdefghij;
// Copied from Occamy project: https://github.com/changkun/occamy/blob/master/protocol/instruction.go#L38
func ParseInstruction(raw []byte) (ins *Instruction, err error) {
	var (
		cursor   int
		elements []string
	)

	bytes := len(raw)
	for cursor < bytes {
		// 1. parse digit
		lengthEnd := -1
		for i := cursor; i < bytes; i++ {
			if raw[i]^'.' == 0 {
				lengthEnd = i
				break
			}
		}
		if lengthEnd == -1 { // cannot find '.'
			return nil, ErrInstructionMissDot
		}
		length, err := strconv.Atoi(string(raw[cursor:lengthEnd]))
		if err != nil {
			return nil, ErrInstructionBadDigit
		}

		// 2. parse rune
		cursor = lengthEnd + 1
		element := new(strings.Builder)
		for i := 1; i <= length; i++ {
			r, n := utf8.DecodeRune(raw[cursor:])
			if r == utf8.RuneError {
				return nil, ErrInstructionBadRune
			}
			cursor += n
			element.WriteRune(r)
		}
		elements = append(elements, element.String())

		// 3. done
		if cursor == bytes-1 {
			break
		}

		// 4. parse next
		if raw[cursor]^',' != 0 {
			return nil, ErrInstructionMissComma
		}

		cursor++
	}

	//noinspection ALL
	return NewInstruction(elements[0], elements[1:]...), nil
}

func (i *Instruction) String() string {
	var b = strings.Builder{}

	b.WriteString(strconv.Itoa(len(i.Opcode)))
	b.WriteString(".")
	b.WriteString(i.Opcode)

	for _, a := range i.Args {
		b.WriteString(",")
		b.WriteString(strconv.FormatInt(int64(utf8.RuneCountInString(a)), 10))
		b.WriteString(".")
		b.WriteString(a)
	}
	b.WriteString(";")

	return b.String()
}
