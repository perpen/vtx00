package main

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// FIXME delete this test
func TestRectMinus(t *testing.T) {
	tests := []struct {
		desc     string
		a, b     rect
		expected []rect
	}{
		{
			"no overlap",
			rect{0, 0, 2, 2}, rect{3, 3, 5, 5},
			[]rect{
				rect{0, 0, 2, 2},
			},
		},
		{
			"top-left corner in",
			rect{0, 0, 4, 4}, rect{2, 2, 6, 6},
			[]rect{
				rect{0, 0, 4, 2},
				rect{0, 2, 2, 2},
			},
		},
		{
			"top-right corner in",
			rect{3, 1, 4, 4}, rect{1, 3, 4, 4},
			[]rect{
				rect{3, 1, 4, 2},
				rect{5, 3, 2, 2},
			},
		},
		{
			"top in",
			rect{1, 1, 6, 4}, rect{3, 3, 2, 2},
			[]rect{
				rect{1, 1, 6, 2},
				rect{1, 3, 2, 2},
				rect{5, 3, 2, 2},
			},
		},
		{
			"all in",
			rect{1, 1, 6, 6}, rect{3, 3, 2, 2},
			[]rect{
				rect{1, 1, 6, 2},
				rect{1, 3, 2, 2},
				rect{5, 3, 2, 2},
				rect{1, 5, 6, 2},
			},
		},
		{
			"left in",
			rect{1, 1, 4, 6}, rect{3, 3, 4, 2},
			[]rect{
				rect{1, 1, 4, 2},
				rect{1, 3, 2, 2},
				rect{1, 5, 4, 2},
			},
		},
		{
			"right in",
			rect{3, 1, 4, 6}, rect{1, 3, 4, 2},
			[]rect{
				rect{3, 1, 4, 2},
				rect{5, 3, 2, 2},
				rect{3, 5, 4, 2},
			},
		},
	}

	for i, test := range tests {
		log.Infof("running test %v: %v", i, test.desc)
		actual := test.a.minus(test.b)
		assert.Equal(t, test.expected, actual,
			fmt.Sprintf("test %v: %v", i, test.desc))
	}
}
