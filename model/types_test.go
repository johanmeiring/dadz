package model_test

import (
	"dadz/model"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"testing"
)

// This test ensures that a user's API key isn't accidentally leaked in an HTTP response.
func TestUserJSONOutput(t *testing.T) {
	user := model.User{
		ID:     5,
		Name:   "The Joker",
		ApiKey: "ABC123",
	}
	b, err := json.Marshal(user)
	assert.Nil(t, err)
	assert.Equal(t, `{"id":5,"name":"The Joker"}`, string(b))

	joke := model.Joke{
		ID:        12,
		Intro:     "This is",
		Punchline: "Funny",
		UserID:    aws.Int(5),
		User:      &user,
	}
	b, err = json.Marshal(joke)
	assert.Nil(t, err)
	assert.NotContains(t, "key", string(b))
}
