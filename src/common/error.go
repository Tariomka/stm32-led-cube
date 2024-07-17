package common

import "errors"

func ErrIfOutOfBounds(num uint8, name string) error {
	if num < 8 {
		return nil
	}

	// println(name, " is out of range[0-7]. Received layer:", num)
	return errors.New(name + "out of bounds")
}
