package utils

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	fakeUserId := primitive.NewObjectID()

	resultString := GenerateToken(fakeUserId)
	assert.NotEmpty(t, resultString)
}

func TestValidateToken(t *testing.T) {
	fakeToken := "fakeToken"

	_, _, err := ValidateToken(fakeToken)

	assert.Error(t, err)
}
