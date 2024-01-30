package main // Replace with the actual package name

import (
	"testing"
)

func TestExtractDirName(t *testing.T) {

	tests := []struct {
		description string
		want        string
		input       string
	}{
		{description: "It should fail", want: "vericred_wi_icare_", input: "vericred_wi_icare_20240101.json.zip"},
		{description: "It should pass", want: "vericred_wi_icare", input: "vericred_wi_icare_20240101.json.zip"},
		{description: "It should fail", want: "vericred_wi_icare_", input: "vericred_wi_icare_20240101.json"},
		{description: "It should pass", want: "vericred_wi_icare", input: "vericred_wi_icare_20240101.json.zip"},
		{description: "It should fail", want: "vericred_wi_icare_", input: "vericred_wi_icare_20240101_summary_report.html"},
		{description: "It should pass", want: "vericred_wi_icare", input: "vericred_wi_icare_20240101_summary_report.html"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got := extractDirName(test.input)

			if test.description == "It should fail" {
				if got == test.want {
					t.Errorf("Expected %q, got %q", test.want, got)
				}
			} else {
				if got != test.want {
					t.Errorf("Expected %q, got %q", test.want, got)
				}
			}
		})

	}
}
