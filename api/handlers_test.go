package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAlphanumeric(t *testing.T) {
	var tests = []struct {
		s             string
		expectedValid bool
	}{
		{s: "foo-bar-213", expectedValid: true},
		{s: "foo_bar_213", expectedValid: true},
		{s: "foo-!@#$bar-213", expectedValid: false},
	}

	for _, tc := range tests {
		validation, err := validateAlphanumericString(tc.s)

		if tc.expectedValid {
			assert.NoError(t, err)
			assert.True(t, validation.valid)
		} else {
			assert.NoError(t, err)
			assert.False(t, validation.valid)
			assert.NotEmpty(t, validation.message)
		}
	}
}

func TestValidateLength(t *testing.T) {
	var tests = []struct {
		s             string
		min           int
		max           int
		expectedValid bool
	}{
		{s: "foo", min: 2, max: 3, expectedValid: true},
		{s: "foo", min: 4, max: 8, expectedValid: false},
		{s: "foo", min: 1, max: 2, expectedValid: false},
	}

	for _, tc := range tests {
		validation := validateLength(tc.s, tc.min, tc.max)

		if tc.expectedValid {
			assert.True(t, validation.valid)
		} else {
			assert.False(t, validation.valid)
			assert.NotEmpty(t, validation.message)
		}
	}
}
