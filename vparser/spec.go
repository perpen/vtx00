package vparser

import (
	"fmt"
)

type ControlSpec struct {
	Name           string
	ParamsNumber   int
	ParamsDefaults []int
	Triggers       []Trigger
	Selections     map[int]string
	UserData       interface{}
}

type Trigger struct {
	Set      ControlSet
	Sequence []byte
}

type ControlSet uint8

// Type of input sequence.
const (
	SetC01 ControlSet = iota
	SetESC
	SetCSI
)

var ControlSetNames = [...]string{
	"C01",
	"ESC",
	"CSI",
}

// MergeDefaults returns the params, augmented with default values indicated by the spec.
func (spec *ControlSpec) MergeDefaults(params []int) []int {
	numDefaults := len(spec.ParamsDefaults)
	numProvided := len(params)
	numParams := numProvided
	if numProvided < numDefaults {
		numParams = numDefaults
	}
	fullParams := make([]int, numParams)
	for i := 0; i < numProvided; i++ {
		fullParams[i] = params[i]
	}
	if numParams != numProvided {
		for i := numProvided; i < numParams; i++ {
			fullParams[i] = spec.ParamsDefaults[i]
		}
	}
	return fullParams
}

func (t Trigger) hash() uint32 {
	var hash uint32
	hash |= uint32(t.Set) << (uint(0) * 8)
	for i, b := range t.Sequence {
		hash |= uint32(b) << (uint(i+1) * 8)
	}
	return hash
}

// Returns a human-readable representation of the control sequence.
func (t Trigger) pretty() string {
	setHuman := ControlSetNames[t.Set]
	switch len(t.Sequence) {
	case 1:
		return fmt.Sprintf("%s %c", setHuman, t.Sequence[0])
	case 2:
		return fmt.Sprintf("%s %c %c", setHuman, t.Sequence[0], t.Sequence[1])
	case 3:
		return fmt.Sprintf("%s %c %c %c", setHuman, t.Sequence[0], t.Sequence[1], t.Sequence[2])
	default:
		return fmt.Sprintf("BUG %v %v", t.Set, t.Sequence)
	}
}
