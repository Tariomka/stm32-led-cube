package utils

import "github.com/Tariomka/led-common-lib/pkg/common"

func ErrIfOutOfBounds(index uint8) error {
	if index < 8 {
		return nil
	}

	return common.ErrOutOfBounds
}
