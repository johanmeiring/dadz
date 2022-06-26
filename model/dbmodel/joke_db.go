package dbmodel

import (
	"context"
	"dadz/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/thoas/go-funk"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"math/rand"
	"time"
)

// JokeDB is a container for a connection to the database.
type JokeDB struct {
	db bun.IDB
}

// NewJokeDB returns a pointer to a value of JokeDB with the provided DB connection value passed in.
func NewJokeDB(username, password, host, port, database string) (*JokeDB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	return &JokeDB{db: db}, nil
}

// RetrieveRandomJoke fetches a random joke from the DB.
func (jdb *JokeDB) RetrieveRandomJoke() (model.Joke, error) {
	// Attempt decent randomness.
	rand.Seed(time.Now().UTC().UnixNano())

	// Determine the max ID of a joke record in the DB.
	var maxID int
	err := jdb.db.
		NewSelect().
		Model((*Joke)(nil)).
		ColumnExpr("max(id)").
		Limit(1).
		Scan(context.Background(), &maxID)
	if err != nil {
		return model.Joke{}, fmt.Errorf("couldn't determine number of jokes in the DB: %w", err)
	}

	var joke Joke
	found := false
	// We iterate, because there could potentially be gaps in IDs. Continue until we find a record.
	for !found {
		id := rand.Intn(maxID) + 1
		err = jdb.db.
			NewSelect().
			Model(&joke).
			Relation("User").
			Where("?TablePKs = ?", id).
			Scan(context.Background())

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			// Something slightly more serious has gone wrong.
			return model.Joke{}, fmt.Errorf("could not fetch a random joke from the DB: %w", err)
		}

		if err == nil {
			found = true
		}
	}

	return joke.ToDadz(), nil
}

// RetrieveJokes fetches a paged collection of jokes from the DB.
// Potential improvement: use cursor-based pagination instead of traditional limit+offset. This would provide
// consistency in the dataset in the event that a new joke is created while a consumer is busy reading paged model.
func (jdb *JokeDB) RetrieveJokes(page, limit int) ([]model.Joke, error) {
	var jokes []Joke
	err := jdb.db.
		NewSelect().
		Model(&jokes).
		Relation("User").
		Limit(limit).
		Offset((page - 1) * limit).
		OrderExpr("id").
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return mapDBJokesToDadzJokes(jokes), nil
}

// SaveNewJokes stores the provided list of new jokes in the DB, and returns the newly-hydrated values.
func (jdb *JokeDB) SaveNewJokes(jokes []model.Joke, userID int) ([]model.Joke, error) {
	dbJokes := funk.Map(jokes, func(j model.Joke) Joke {
		return NewJokeFromDadz(j, userID)
	}).([]Joke)

	_, err := jdb.db.
		NewInsert().
		Model(&dbJokes).
		Exec(context.Background())

	if err != nil {
		return nil, err
	}

	return mapDBJokesToDadzJokes(dbJokes), nil
}

func mapDBJokesToDadzJokes(dbJokes []Joke) []model.Joke {
	return funk.Map(dbJokes, func(j Joke) model.Joke {
		return j.ToDadz()
	}).([]model.Joke)
}

// GetUserByApiKey retrieves a user based on API key.
func (jdb *JokeDB) GetUserByApiKey(key string) (model.User, error) {
	var user User
	err := jdb.db.
		NewSelect().
		Model(&user).
		Where("api_key = ?", key).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		return model.User{}, err
	}

	return user.ToDadz(), nil
}
