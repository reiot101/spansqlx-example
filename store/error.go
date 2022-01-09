package store

const (
	ErrAccountNotFound    = Error("account not found")
	ErrEmailAlreadyExists = Error("account already exists")
	ErrTaskNotFound       = Error("task not found")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
