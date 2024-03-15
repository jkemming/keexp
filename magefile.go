//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Builds the binary.
func Build() error {
	return sh.Run("go", "build", "jkemming.com/keexp")
}
