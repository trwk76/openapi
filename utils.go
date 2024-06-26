package openapi

import (
	"fmt"
	"strings"
)

func refKey(ref string) string {
	if idx := strings.LastIndexByte(ref, '/'); idx >= 0 {
		return ref[idx+1:]
	}

	return ref
}

func uniqueName[T any](c map[string]T, base string) string {
	if _, fnd := c[base]; !fnd {
		return base
	}

	i := 1
	name := fmt.Sprintf("%s%d", base, i)

	for _, fnd := c[name]; fnd; {
		i += 1
		name = fmt.Sprintf("%s%d", base, i)
	}

	return name
}
