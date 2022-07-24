package Models

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValid_NameEmpty_False(t *testing.T) {
	id, _ := uuid.NewV4()
	tag := Tag{UserId: id}

	assert.False(t, tag.IsValid())
}

func TestIsValid_UserIdEmpty_False(t *testing.T) {
	tag := Tag{Name: "test"}

	assert.False(t, tag.IsValid())
}

func TestIsValid_FieldsFilled_True(t *testing.T) {
	id, _ := uuid.NewV4()
	tag := Tag{UserId: id, Name: "test"}

	assert.True(t, tag.IsValid())
}
