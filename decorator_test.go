package gw

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAllPermDecorator(t *testing.T) {
	resource := "User"
	permList := NewPermAllDecorator(resource)
	assert.True(t, permList.Has("Creation"))
	assert.True(t, permList.Has("Deletion"))
	assert.True(t, permList.Has("Modification"))
	assert.True(t, permList.Has("ReadAll"))
	assert.True(t, permList.Has("ReadDetail"))
	assert.True(t, permList.Has("Administration"))
}

func TestNewCrudPermDecorator(t *testing.T) {
	resource := "User"
	permList := NewPermAllDecorator(resource)
	assert.True(t, permList.Has("Creation"))
	assert.True(t, permList.Has("Deletion"))
	assert.True(t, permList.Has("Modification"))
	assert.True(t, permList.Has("ReadAll"))
	assert.True(t, permList.Has("ReadDetail"))
}
