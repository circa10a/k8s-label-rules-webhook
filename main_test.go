package main

import "testing"

func TestFlags(t *testing.T) {
	t.Parallel()
	// Ensure no issue comes from flags()
	flags()
}
