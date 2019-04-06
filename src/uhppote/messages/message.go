package messages

type Field struct {
	Offset uint8
	Length uint8
	Value  interface {
	}
}

type Message struct {
	Code   byte
	Fields []Field
}
