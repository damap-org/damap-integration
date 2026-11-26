package domain

func (e Error) Error() string {
	return e.Message
}
