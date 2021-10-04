package xname

type ValidationError struct {
	ValidationFailures []string
}

func (e *ValidationError) Error() string {

}