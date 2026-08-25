package adapter

import "fmt"

func NormalizeConversionError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("codec call failed: %v", err)
}

func IsTemporary(err error) bool { return err != nil }
