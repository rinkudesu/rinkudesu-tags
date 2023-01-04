package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValid_NameEmpty_False(t *testing.T) {
	tag := Tag{}

	assert.False(t, tag.IsValid())
}

func TestIsValid_FieldsFilled_True(t *testing.T) {
	tag := Tag{Name: "test"}

	assert.True(t, tag.IsValid())
}
