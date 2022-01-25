package errors

import (
	"errors"
	"fmt"
)

var (
	errNotFound     = errors.New("not found")
	errNotDestroyed = errors.New("not destroyed")
	errMismatched   = errors.New("mismatched")
)

func ResourceNotFoundError(key string) error {
	return fmt.Errorf("resource - %s %w", key, errNotFound)
}

func ResourceNotDestroyedError(key string) error {
	return fmt.Errorf("resource - %s %w", key, errNotDestroyed)
}

func ResourceTypeMismatched(expected, found string) error {
	return fmt.Errorf("resource type %w : expected - %s : found - %s", errMismatched, expected, found)
}
