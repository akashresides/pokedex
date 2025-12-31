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
                		},        }
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
