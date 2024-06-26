package openapi

import (
	"fmt"
	"slices"

	"github.com/trwk76/openapi/spec"
)

func (a *API) Tag(name string, desc string) Tag {
	fnd := slices.ContainsFunc(a.s.Tags, func(i spec.Tag) bool { return i.Name == name })
	if fnd {
		panic(fmt.Errorf("tag '%s' is already declared", name))
	}

	a.s.Tags = append(a.s.Tags, spec.Tag{Name: name, Description: desc})
	return Tag{name: name}
}

type (
	Tag struct {
		name string
	}
)
