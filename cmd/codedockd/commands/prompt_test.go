package commands

import (
	"strings"
	"testing"
)

func TestSetupInputPreservesFullNamesAndNextPrompt(t *testing.T) {
	reader := strings.NewReader("Jane Doe\r\n password with spaces \n")
	name, err := readInputLine(reader)
	if err != nil {
		t.Fatal(err)
	}
	if name != "Jane Doe" {
		t.Fatal("multiword name was truncated")
	}
	password, err := readInputLine(reader)
	if err != nil {
		t.Fatal(err)
	}
	if password != " password with spaces " {
		t.Fatal("password input or the next prompt was altered")
	}
}

func TestSetupInputAcceptsFinalLineWithoutNewline(t *testing.T) {
	value, err := readInputLine(strings.NewReader("owner@example.com"))
	if err != nil {
		t.Fatal(err)
	}
	if value != "owner@example.com" {
		t.Fatal("final line was lost")
	}
}
