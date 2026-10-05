package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

func prompt(msg string) string {
	fmt.Print(msg)
	input, err := readInputLine(os.Stdin)
	if err != nil {
		exitError("Read input: %v", err)
	}
	return strings.TrimSpace(input)
}

func promptOptional(msg string) string { return prompt(msg) }

func promptPassword(msg string) string {
	fmt.Fprint(os.Stderr, msg)
	if term.IsTerminal(int(os.Stdin.Fd())) {
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			exitError("Read password: %v", err)
		}
		return string(password)
	}
	password, err := readInputLine(os.Stdin)
	if err != nil {
		exitError("Read password: %v", err)
	}
	return password
}

func readInputLine(reader io.Reader) (string, error) {
	var input strings.Builder
	var character [1]byte
	for {
		if _, err := io.ReadFull(reader, character[:]); err != nil {
			if errors.Is(err, io.EOF) {
				return input.String(), nil
			}
			return "", fmt.Errorf("read input line: %w", err)
		}
		if character[0] == '\n' {
			return input.String(), nil
		}
		if character[0] != '\r' {
			input.WriteByte(character[0])
		}
	}
}
