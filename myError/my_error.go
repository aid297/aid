package myError

type (
	IMyError interface {
		New(msg string) IMyError
		Wrap(err error) IMyError
		Panic() IMyError
		Error() string
		Is(target error) bool
	}
	MyError struct{ Msg string }
)
