package githubapi

import (
	"fmt"
)

func WrapError(message string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", message, err)
	}

	return nil
}
