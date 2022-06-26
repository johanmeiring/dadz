package dbmodel

import (
	"dadz/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/uptrace/bun"
)

// Joke represents a `joke` record from the `jokes` table.
type Joke struct {
	bun.BaseModel `bun:"jokes"`
	ID            int `bun:",pk,autoincrement"`
	Intro         string
	Punchline     string
	UserID        *int  `bun:",nullzero"`
	User          *User `bun:"rel:belongs-to,join:user_id=id"`
}

// ToDadz converts a bun-specific Joke value to model.Joke for use by the application.
func (j Joke) ToDadz() model.Joke {
	joke := model.Joke{
		ID:        j.ID,
		Intro:     j.Intro,
		Punchline: j.Punchline,
		UserID:    j.UserID,
	}

	if j.User != nil {
		userVal := *j.User
		dadzUser := userVal.ToDadz()
		joke.User = &dadzUser
	}

	return joke
}

// NewJokeFromDadz creates a Joke value from an incoming model.Joke value.
func NewJokeFromDadz(j model.Joke, userID int) Joke {
	return Joke{
		ID:        j.ID,
		UserID:    aws.Int(userID),
		Intro:     j.Intro,
		Punchline: j.Punchline,
	}
}

// User represents a `user` record from the `users` table.
type User struct {
	bun.BaseModel `bun:"users"`
	ID            int `bun:",pk,autoincrement"`
	Name          string
	ApiKey        string
}

// ToDadz converts a bun-specific User value to model.User for use by the application.
func (u User) ToDadz() model.User {
	return model.User{
		ID:     u.ID,
		Name:   u.Name,
		ApiKey: u.ApiKey,
	}
}
