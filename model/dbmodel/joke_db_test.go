package dbmodel_test

import (
	"dadz/model"
	"dadz/model/dbmodel"
	"database/sql"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getTestDB() *dbmodel.JokeDB {
	db, err := dbmodel.NewJokeDB("postgres", "blah", "localhost", "5434", "postgres")
	if err != nil {
		panic(err)
	}

	return db
}

func TestJokeDB_RetrieveRandomJoke(t *testing.T) {
	db := getTestDB()
	joke, err := db.RetrieveRandomJoke()

	assert.Nil(t, err)
	assert.IsType(t, model.Joke{}, joke)

	// This check will fail sometimes, obviously. But it serves its purpose well regardless.
	joke2, _ := db.RetrieveRandomJoke()
	assert.NotEqual(t, joke.ID, joke2.ID)
}

func TestJokeDB_RetrieveJokes(t *testing.T) {
	db := getTestDB()
	jokes, err := db.RetrieveJokes(1, 10)
	assert.Nil(t, err)
	assert.Len(t, jokes, 10)

	jokes2, err := db.RetrieveJokes(2, 10)
	assert.Nil(t, err)
	assert.Len(t, jokes2, 10)
	// Ensure that none of the results in the second dataset are present in the first dataset.
	for _, j2 := range jokes2 {
		for _, j1 := range jokes {
			assert.NotEqual(t, j1.ID, j2.ID)
		}
	}
}

func TestJokeDB_SaveNewJokes(t *testing.T) {
	db := getTestDB()

	var err error

	newJokes := []model.Joke{
		{
			Intro:     "This is",
			Punchline: "Funny",
			UserID:    aws.Int(1),
		},
		{
			Intro:     "This, too",
			Punchline: "is Funny",
			UserID:    aws.Int(1),
		},
	}

	// Double check that ID was not set.
	assert.Zero(t, newJokes[0].ID)
	newJokes, err = db.SaveNewJokes(newJokes, 1)
	assert.Nil(t, err)
	assert.Len(t, newJokes, 2)
	// Check that IDs have now been set by Postgres.
	assert.NotZero(t, newJokes[0].ID)
}

func TestJokeDB_GetUserByApiKeySuccess(t *testing.T) {
	db := getTestDB()
	user, err := db.GetUserByApiKey("9e44bed1ab8845508297079d4a0117b4")

	assert.Nil(t, err)
	assert.IsType(t, model.User{}, user)
}

func TestJokeDB_GetUserByApiKeyFailure(t *testing.T) {
	db := getTestDB()
	_, err := db.GetUserByApiKey("abcabc123123")

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}
