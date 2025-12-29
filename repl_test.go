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
	};

	for _, c := range cases {

		actualSlice := cleanInput(c.input)

		if len(actualSlice) != len(c.expected) {
			t.Errorf("Different length: %v and %v", len(actualSlice), len(c.expected))
		}

		for i, actualWord := range actualSlice {
			expectedWord := c.expected[i]
			if actualWord != expectedWord {
				t.Errorf("Unexpected word: %v and %v", actualWord, expectedWord)
			}
		}
	}
}
