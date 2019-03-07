package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordBreak(t *testing.T) {
	tests := []struct {
		str    string
		width  int
		center bool
		result []string
	}{
		{
			str:    "The Text Justification algorithm will ensure that the output from your program is both left and right justified when displayed in amono-spaced font such as Courier. This paragraph is an example of such justification. All the lines (except the last) of the output from a given run of your program should have the same length, and the last line is to be no longer than the other lines.",
			width:  100,
			center: false,
			result: []string{
				"The Text Justification algorithm will ensure that the output from your program is both left and",
				"right justified when displayed in amono-spaced font such as Courier. This paragraph is an example",
				"of such justification. All the lines (except the last) of the output from a given run of your",
				"program should have the same length, and the last line is to be no longer than the other lines.",
			},
		},
		{
			str:    "pneumonoultramicroscopicsilicovolcanoconiosis",
			width:  10,
			center: false,
			result: []string{
				"pneumonoultramicroscopicsilicovolcanoconiosis",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.result, Justify(test.str, test.width, test.center))
	}
}
