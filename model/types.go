package model

// Provider defines an interface for interacting with a model source. It can be implemented by a database, caching engine
// file reader, HTTP API consumer, etc.
type Provider interface {
	SaveNewJokes(jokes []Joke, userID int) ([]Joke, error)
	RetrieveJokes(page, limit int) ([]Joke, error)
	RetrieveRandomJoke() (Joke, error)
	GetUserByApiKey(key string) (User, error)
}

// Joke represents a `joke` record from any model source.
type Joke struct {
	ID        int    `json:"id"`
	Intro     string `json:"intro"`
	Punchline string `json:"punchline"`
	UserID    *int   `json:"-"`
	User      *User  `json:"user,omitempty"`
}

// User represents a `user` record from any model source.
type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	ApiKey string `json:"-"`
}
