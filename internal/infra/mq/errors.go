package mq

type SkippableError struct {
	err error
}

func NewSkippableError(err error) *SkippableError {
	return &SkippableError{err: err}
}

func (e *SkippableError) Error() string {
	return e.err.Error()
}
