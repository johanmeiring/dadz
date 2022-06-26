package main

import (
	"bytes"
	"dadz/model"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a *App

// MockDataProvider represents the principle of being able to test
// HTTP handlers independently of any actual data source.
type MockDataProvider struct{}

func (m *MockDataProvider) RetrieveRandomJoke() (model.Joke, error) {
	user := model.User{
		ID:     1,
		Name:   "The Joker",
		ApiKey: "abc123",
	}
	return model.Joke{
		ID:        5,
		Intro:     "This is",
		Punchline: "Funny",
		UserID:    aws.Int(1),
		User:      &user,
	}, nil
}

func (m *MockDataProvider) RetrieveJokes(page, limit int) ([]model.Joke, error) {
	user := model.User{
		ID:     1,
		Name:   "The Joker",
		ApiKey: "abc123",
	}
	return []model.Joke{
		{
			ID:        5,
			Intro:     "This is",
			Punchline: "Funny",
			UserID:    aws.Int(1),
			User:      &user,
		},
		{
			ID:        6,
			Intro:     "This, too",
			Punchline: "is Funny",
			UserID:    aws.Int(1),
			User:      &user,
		},
	}, nil
}

func (m *MockDataProvider) SaveNewJokes(jokes []model.Joke, userID int) ([]model.Joke, error) {
	user := model.User{
		ID:     1,
		Name:   "The Joker",
		ApiKey: "abc123",
	}
	return []model.Joke{
		{
			ID:        5,
			Intro:     "This is",
			Punchline: "Funny",
			UserID:    aws.Int(1),
			User:      &user,
		},
		{
			ID:        6,
			Intro:     "This, too",
			Punchline: "is Funny",
			UserID:    aws.Int(1),
			User:      &user,
		},
	}, nil
}

func (m *MockDataProvider) GetUserByApiKey(key string) (model.User, error) {
	return model.User{
		ID:     1,
		Name:   "The Joker",
		ApiKey: "abc123",
	}, nil
}

func TestMain(m *testing.M) {
	a = NewApp(&MockDataProvider{})
	code := m.Run()
	os.Exit(code)
}

func TestGetRandomJoke(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/random-joke", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check that the response body is a "joke".
	var joke model.Joke
	err := json.NewDecoder(rr.Body).Decode(&joke)
	assert.Nil(t, err)
}

func TestGetJokes(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/jokes", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var jokes []model.Joke
	err := json.NewDecoder(rr.Body).Decode(&jokes)
	assert.Nil(t, err)
	assert.Len(t, jokes, 2)
}

func TestPostJokesNoToken(t *testing.T) {
	jokes := []model.Joke{
		{
			ID:        5,
			Intro:     "This is",
			Punchline: "Funny",
		},
		{
			ID:        6,
			Intro:     "This, too",
			Punchline: "is Funny",
		},
	}
	b, _ := json.Marshal(jokes)
	req, _ := http.NewRequest(http.MethodPost, "/jokes", bytes.NewReader(b))
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var errs map[string]string
	json.NewDecoder(rr.Body).Decode(&errs)
	assert.NotEmpty(t, errs)
	assert.NotEmpty(t, errs["error"])
	assert.Equal(t, errs["error"], "authorization token is required")
}

func TestPostJokesSuccess(t *testing.T) {
	jokes := []model.Joke{
		{
			ID:        5,
			Intro:     "This is",
			Punchline: "Funny",
		},
		{
			ID:        6,
			Intro:     "This, too",
			Punchline: "is Funny",
		},
	}
	b, _ := json.Marshal(jokes)
	req, _ := http.NewRequest(http.MethodPost, "/jokes", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer abc123")
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	err := json.NewDecoder(rr.Body).Decode(&jokes)
	assert.Nil(t, err)
	// Check that jokes are actually returned as part of the response.
	assert.Len(t, jokes, 2)
}
