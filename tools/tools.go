//go:build tools

package tools

import (
	// Pin displaywidth to a version compatible with Go 1.26.
	// v0.6.2 (ent's minimum) fails to compile on Go 1.26.
	_ "github.com/clipperhouse/displaywidth"
	// Pin entc so go mod tidy keeps ent code-gen dependencies.
	_ "entgo.io/ent/entc"
)
