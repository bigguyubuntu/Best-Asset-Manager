package api

import "testing"

func TestExtractSuffixId(t *testing.T) {
	input := []string{
		"www.google.com/userId",
		"banana.ca/apples/mangos/12",
		"/",
		"sup",
		"//",
		"//////",
	}
	expected := []string{
		"userId",
		"12",
		"",
		"",
		"",
		"",
	}
	for i := 0; i < len(input); i++ {
		in := input[i]
		e := expected[i]
		out := ExtractSuffixId(in)
		if e != out {
			t.Errorf("Expcted result to be %s, but it was %s", e, out)
		}
	}
}
