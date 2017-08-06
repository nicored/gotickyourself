package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidAlis(t *testing.T) {
	testCases := []struct {
		Alias string
		Valid bool
	}{
		{"dev", true},
		{"dev-dms", true},
		{"  dev  ", true},
		{"d.d", true},
		{"1dev", false},
		{"1.2", false},
		{"today", false},
		{"yesterday", false},
		{"   yesterday   ", false},
		{"today-dms", true},
		{"3 days", false},
	}

	for _, tc := range testCases {
		valid, _ := isValidAlias(tc.Alias)
		assert.Equal(t, tc.Valid, valid, tc.Alias)
	}
}
