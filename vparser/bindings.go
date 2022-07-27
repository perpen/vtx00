package vparser

import (
	"fmt"
)

type Bindings struct {
	specsByTrigger map[uint32]*ControlSpec
	specsByName    map[string]*ControlSpec
}

func NewBindings(specs map[string]ControlSpec) *Bindings {
	bindings := new(Bindings)
	bindings.specsByTrigger = make(map[uint32]*ControlSpec)
	bindings.specsByName = make(map[string]*ControlSpec)
	for _, spec := range specs {
		bindings.registerControl(spec)
	}
	return bindings
}

// Called from ctor to add the control to the lookup maps.
// Errors if the control name is already in use.
func (bindings *Bindings) registerControl(spec ControlSpec) error {
	name := spec.Name
	_, found := bindings.specsByName[name]
	if found {
		return fmt.Errorf("RegisterControl: already have %s", name)
	}

	// Basic validation
	if spec.ParamsNumber != -1 && len(spec.ParamsDefaults) > spec.ParamsNumber {
		return fmt.Errorf("RegisterControl: inconsistent number of default params for %s", name)
	}

	bindings.specsByName[name] = &spec

	for _, trigger := range spec.Triggers {
		// Check sequence not in use
		hash := trigger.hash()
		prevTrigger, found := bindings.specsByTrigger[hash]
		if found {
			return fmt.Errorf("addTrigger: cannot add %v for %v, already used for control %v",
				trigger.pretty(), name, prevTrigger.Name)
		}

		bindings.specsByTrigger[hash] = &spec
	}
	return nil
}

func (bindings *Bindings) LookupSpec(name string) (*ControlSpec, bool) {
	spec, found := bindings.specsByName[name]
	return spec, found
}

func (bindings *Bindings) specForTrigger(trigger Trigger) (*ControlSpec, bool) {
	hash := trigger.hash()
	spec, found := bindings.specsByTrigger[hash]
	if !found {
		return nil, false
	}
	return spec, true
}
