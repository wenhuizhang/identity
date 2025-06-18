// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package cmdutil

import (
	"fmt"
	"os"
)

const (
	// If the user just presses Enter, the input will be "" and err will be an "unexpected newline" error.
	errNewLine = "unexpected newline"
)

func ScanRequiredIfNotSet(msg string, in *string) error {
	if in != nil && *in != "" {
		return nil
	}

	return ScanRequired(msg, in)
}

func ScanRequired(msg string, in *string) error {
	fmt.Fprintf(os.Stdout, "%s: ", msg)

	_, err := fmt.Scanln(in)
	if err != nil {
		if err.Error() != errNewLine {
			return err
		}
	}

	if *in == "" {
		return fmt.Errorf("field cannot be empty")
	}

	return nil
}

func ScanOptionalIfNotSet(msg string, in *string) error {
	if in != nil && *in != "" {
		return nil
	}

	return ScanOptional(msg, in)
}

func ScanOptional(msg string, in *string) error {
	fmt.Fprintf(os.Stdout, "(Optional) %s: ", msg)

	_, err := fmt.Scanln(in)
	if err != nil {
		if err.Error() != errNewLine {
			return err
		}
	}

	return nil
}

func ScanWithDefaultIfNotSet(msg, defaultValue string, in *string) error {
	if in != nil && *in != "" {
		return nil
	}

	return ScanWithDefault(msg, defaultValue, in)
}

func ScanWithDefault(msg, defaultValue string, in *string) error {
	fmt.Fprintf(os.Stdout, "%s (default %s): ", msg, defaultValue)

	_, err := fmt.Scanln(in)
	if err != nil {
		if err.Error() != errNewLine {
			return err
		}
	}

	if *in == "" {
		*in = defaultValue
	}

	return nil
}
