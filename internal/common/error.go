package common

type OutOfBoundsError struct{}

func (e OutOfBoundsError) Error() string {
	return "index out of bounds"
}

func ErrIfOutOfBounds(index uint8) error {
	if index < 8 {
		return nil
	}

	return OutOfBoundsError{}
}
