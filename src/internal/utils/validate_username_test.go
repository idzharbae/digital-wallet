package utils_test

import (
	"testing"

	"github.com/idzharbae/digital-wallet/src/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	t.Run("Must return false if it contains space", func(t *testing.T) {
		result := utils.ValidateUserName("asdf esdj")
		assert.False(t, result)
	})

	t.Run("Must return false if it contains special characters", func(t *testing.T) {
		result := utils.ValidateUserName("asdf-esdj")
		assert.False(t, result)

		result = utils.ValidateUserName("asdf!esdj")
		assert.False(t, result)
	})

	t.Run("Must return false if it starts with number", func(t *testing.T) {
		result := utils.ValidateUserName("123asdasd")
		assert.False(t, result)

		result = utils.ValidateUserName("4asdf")
		assert.False(t, result)
	})

	t.Run("Must return true if it starts with alphabet and is alphanumeric", func(t *testing.T) {
		result := utils.ValidateUserName("testingA123")
		assert.True(t, result)

		result = utils.ValidateUserName("ABCDefgh9129")
		assert.True(t, result)
	})
}
