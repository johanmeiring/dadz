package dbmodel_test

import (
	"dadz/model"
	"dadz/model/dbmodel"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoke_ToDadz(t *testing.T) {
	j := dbmodel.Joke{
		ID:        12,
		Intro:     "This is",
		Punchline: "Funny",
		UserID:    aws.Int(5),
		User: &dbmodel.User{
			ID:     5,
			Name:   "The Joker",
			ApiKey: "abc123",
		},
	}

	dj := j.ToDadz()

	assert.IsType(t, model.Joke{}, dj)
	assert.NotNil(t, dj.User)
	assert.Equal(t, "This is", dj.Intro)
	assert.Equal(t, "Funny", dj.Punchline)
	assert.Equal(t, 12, dj.ID)
}

func TestUser_ToDadz(t *testing.T) {
	dbUser := dbmodel.User{
		ID:     5,
		Name:   "The Joker",
		ApiKey: "abc123",
	}

	user := dbUser.ToDadz()

	assert.IsType(t, model.User{}, user)
	assert.Equal(t, "The Joker", user.Name)
	assert.Equal(t, "abc123", user.ApiKey)
}

func TestNewJokeFromDadz(t *testing.T) {
	joke := model.Joke{
		ID:        12,
		Intro:     "This is",
		Punchline: "Funny",
		UserID:    aws.Int(5),
	}

	dbJoke := dbmodel.NewJokeFromDadz(joke, 5)

	assert.IsType(t, dbmodel.Joke{}, dbJoke)
	assert.NotNil(t, dbJoke.UserID)
	assert.Equal(t, "This is", dbJoke.Intro)
	assert.Equal(t, "Funny", dbJoke.Punchline)
	assert.Equal(t, 12, dbJoke.ID)
}
