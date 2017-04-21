package message

// Message tells the out plugin how to respond to the event
type Message struct {
	OK      bool
	Content string
	Err     error
}

// NewMessage returns message
func NewMessage(ok bool, content string, err error) *Message {
	return &Message{
		OK:      ok,
		Content: content,
		Err:     err,
	}
}
