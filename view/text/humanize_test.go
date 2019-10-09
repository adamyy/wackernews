package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanize(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "Normally I'd agree, but the fact that the AirPower mat is printed onto the back of the new AirPods case box does seem to be a pretty good indicator that this was a very sudden decision.\nThat is unless they were trying to make it <i>look</i> sudden, but I don't think we have good enough reason to be that cynical about it.",
			output: "Normally I'd agree, but the fact that the AirPower mat is printed onto the back of the new AirPods case box does seem to be a pretty good indicator that this was a very sudden decision.\nThat is unless they were trying to make it [3mlook [0m sudden, but I don't think we have good enough reason to be that cynical about it.",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.output, Humanize(test.input))
	}
}
