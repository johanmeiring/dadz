package main

import (
	"dadz/model"
	"dadz/model/dbmodel"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// App defines the core components and dependencies of our application.
type App struct {
	DataProvider model.Provider
	Router       *mux.Router
}

// NewApp returns a pointer to a value of App, after initialising its data provider and route handlers.
func NewApp(dataProvider model.Provider) *App {
	a := &App{
		DataProvider: dataProvider,
	}

	a.Router = mux.NewRouter()

	a.Router.HandleFunc("/jokes", a.getJokes).Methods(http.MethodGet)
	a.Router.HandleFunc("/random-joke", a.getRandomJoke).Methods(http.MethodGet)
	a.Router.HandleFunc("/jokes", a.postJokes).Methods(http.MethodPost)

	return a
}

func (a *App) getAuthenticatedUser(r *http.Request) (*model.User, error) {
	key := r.Header.Get("Authorization")
	if key == "" {
		return nil, fmt.Errorf("authorization token is required")
	}
	if !strings.HasPrefix(key, "Bearer ") {
		return nil, fmt.Errorf("authorization header value must being with 'Bearer '")
	}

	split := strings.Split(key, " ")
	key = split[1]

	user, err := a.DataProvider.GetUserByApiKey(key)
	if err != nil {
		return nil, fmt.Errorf("invalid authentication token")
	}

	return &user, nil
}

func (a *App) getRandomJoke(w http.ResponseWriter, r *http.Request) {
	joke, err := a.DataProvider.RetrieveRandomJoke()
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, joke)
}

func (a *App) getJokes(w http.ResponseWriter, r *http.Request) {
	var err error

	page, err := getIntFromQuery(r, "page", 1)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}
	limit, err := getIntFromQuery(r, "limit", 20)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	jokes, err := a.DataProvider.RetrieveJokes(page, limit)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, jokes)
}

func getIntFromQuery(r *http.Request, key string, defaultValue int) (int, error) {
	if value, okay := r.URL.Query()[key]; okay {
		intVal, err := strconv.Atoi(value[0])
		if err != nil {
			return 0, fmt.Errorf("could not interpret query param '%s': %w", key, err)
		}

		return intVal, nil
	}

	return defaultValue, nil
}

func (a *App) postJokes(w http.ResponseWriter, r *http.Request) {
	user, err := a.getAuthenticatedUser(r)
	if err != nil {
		errorResponse(w, err, http.StatusUnauthorized)
		return
	}

	var jokes []model.Joke
	err = json.NewDecoder(r.Body).Decode(&jokes)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	jokes, err = a.DataProvider.SaveNewJokes(jokes, user.ID)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonResponse(w, jokes)
}

func errorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	output := map[string]string{
		"error": err.Error(),
	}
	json.NewEncoder(w).Encode(output)
}

func jsonResponse(w http.ResponseWriter, v any) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func main() {
	// TODO: Do not hardcode DB details. These should ideally be passed in as command line flags or fetched from environment variables.
	db, err := dbmodel.NewJokeDB("postgres", "blah", "host.docker.internal", "5433", "postgres")
	if err != nil {
		panic(err)
	}

	a := NewApp(db)
	fmt.Println("Running server on localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}
