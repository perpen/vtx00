package vterm

import (
	"testing"
)

func TestControls(t *testing.T) {
	for _, test := range allControlTests {
		// log.Println("testing", name)
		test(t)
	}
}
