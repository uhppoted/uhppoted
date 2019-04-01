package messages

type Message interface {
	Name() string
	Code() byte
}
