package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemovePrefixes(t *testing.T) {
	testCases := []struct {
		Name     string
		Given    string
		APIPath  string
		Expected Path
	}{
		{
			Name:     "lightwell-next status",
			Given:    "/api/lightwell-next/v1.0/status",
			APIPath:  "lightwell-next",
			Expected: Path{"status"},
		},
		{
			Name:     "lightwell-next ping",
			Given:    "/api/lightwell-next/v1/ping",
			APIPath:  "lightwell-next",
			Expected: Path{"ping"},
		},
		{
			Name:     "lightwell ping after switch",
			Given:    "/api/lightwell/v1/ping",
			APIPath:  "lightwell",
			Expected: Path{"ping"},
		},
		{
			Name:     "wrong api path",
			Given:    "/api/lightwell-next/v1/ping",
			APIPath:  "lightwell",
			Expected: Path{},
		},
		{
			Name:     "plain ping",
			Given:    "/ping",
			APIPath:  "lightwell-next",
			Expected: Path{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.Equal(t, testCase.Expected, NewPathWithString(testCase.Given).RemovePrefixes(testCase.APIPath))
		})
	}
}

func TestStartWithResources(t *testing.T) {
	assert.True(t, Path{"ping"}.StartWithResources([]string{"ping"}))
	assert.False(t, Path{"status"}.StartWithResources([]string{"ping"}))
	assert.False(t, Path{}.StartWithResources([]string{"ping"}))
}
