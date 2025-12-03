package client

type MessageOption struct {
	text  string
	valid bool
}

func WithMessage(text string) MessageOption {
	return MessageOption{
		text:  text,
		valid: true,
	}
}

func KeepExistingMessage() MessageOption {
	return MessageOption{
		text:  "",
		valid: true,
	}
}
