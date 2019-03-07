package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadingText(t *testing.T) {
	gen := LoadingText()
	assert.Equal(t, "loading \\", gen())
	assert.Equal(t, "loading |", gen())
	assert.Equal(t, "loading /", gen())
	assert.Equal(t, "loading -", gen())
	assert.Equal(t, "loading \\", gen())
}
