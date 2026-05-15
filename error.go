// Copyright ©2026 xiayoudi. All rights reserved.
// Author: xiayoudi
// Email: ur@xiaud.com

package ud

import (
	"errors"
	"fmt"
	"strings"
)

func Err(parts ...string) error {
	switch len(parts) {
	case 0:
		return errors.New("unknown error")
	case 1:
		return errors.New(parts[0])
	default:
		return errors.New(strings.Join(parts, ": "))
	}
}

func Wrap(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	context := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %w", context, err)
}
