package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  I am Zidane  ",
			expected: []string{"i", "am", "zidane"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "HELLO WORLD",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello\tworld\n",
			expected: []string{"hello", "world"},
		}}
	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("length of slices don't match")
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("expected: %s, got: %s", expectedWord, word)
				t.Fail()
			}

		}
	}
}

func TestCommandsMapContainsExpectedCommands(t *testing.T) {
	expectedCommands := []string{"help", "exit"}

	for _, cmdName := range expectedCommands {
		_, exists := commands[cmdName]
		if !exists {
			t.Errorf("expected command '%s' to be registered", cmdName)
			t.Fail()
		}
	}
}

func TestCommandStructureValidation(t *testing.T) {
	for key, cmd := range commands {
		if cmd.name != key {
			t.Errorf("command key '%s' does not match cmd.name '%s'", key, cmd.name)
			t.Fail()
		}

		if cmd.description == "" {
			t.Errorf("command '%s' has empty description", cmd.name)
			t.Fail()
		}

		if cmd.callback == nil {
			t.Errorf("command '%s' has nil callback", cmd.name)
			t.Fail()
		}
	}
}

func TestInvalidCommandLookup(t *testing.T) {
	invalidCommands := []string{"nonexistent", "invalid", "fake"}

	for _, cmdName := range invalidCommands {
		cmd, exists := commands[cmdName]
		if exists {
			t.Errorf("unexpectedly found command '%s' in commands map", cmdName)
			t.Fail()
		}
		if cmd.name != "" || cmd.description != "" || cmd.callback != nil {
			t.Errorf("command '%s' should be zero value, got: %+v", cmdName, cmd)
			t.Fail()
		}
	}
}

func TestCommandNamesAreUnique(t *testing.T) {
	seenNames := make(map[string]bool)

	for key, cmd := range commands {
		if seenNames[cmd.name] {
			t.Errorf("duplicate command name found: '%s'", cmd.name)
			t.Fail()
		}
		seenNames[cmd.name] = true

		if key != cmd.name {
			t.Errorf("map key '%s' does not match command name '%s'", key, cmd.name)
			t.Fail()
		}
	}
}

func TestCommandDescriptionsArePopulated(t *testing.T) {
	for _, cmd := range commands {
		if cmd.description == "" {
			t.Errorf("command '%s' has empty description", cmd.name)
			t.Fail()
		}

		trimmed := cmd.description
		if len(trimmed) == 0 {
			t.Errorf("command '%s' description contains only whitespace", cmd.name)
			t.Fail()
		}
	}
}
