# dadz

`dadz` is a dad joke API served over HTTP. Its name is a callback to good old "icanhascheezburger".

## How to run it

A Makefile has been included for convenience. The following targets are available:

- `make run`: Starts a Docker-based DB instance with data, waits 5 seconds for it to become consistent, and builds and runs the `johanmeiring/dadz` image, thus serving the API at `localhost:8080`.
- `make build`: Builds a `dadz` binary for 64-bit x86 Linux and outputs it to `build/dadz`.
- `make build-docker`: Builds a Docker image tagged `johanmeiring/dadz` containing the `dadz` binary.
- `make clean-db`: Stop and deletes the DB instance.
- `make db`: Starts a docker-based DB instance with data, for use in development.
- `make test-db`: Starts a docker-based DB instance with data, specifically for using with tests.
- `make clean-test-db`: Stops and deletes the test DB instance.
- `make test`: Run the full test suite using the test DB instance.
- `make test-cover`: Run the full test suite using the test DB instance, and output coverage details.

If you simply want to build and run the application, then a simple `make` should do the trick.

## API endpoints

- `GET /random-joke`: Returns a random joke.
- `GET /jokes`: Returns an array of jokes. Accepts `page` and `limit` query string parameters for paging.  The default limit per page is 20 jokes.
- `POST /jokes`: Accepts an array of joke objects and creates them in the DB. Requires the `Authorization` header to be set in the format `Authorization Bearer <api key>`.
  
  Example call:

  `$ curl -X POST -d '[{"intro":"This is","punchline":"Funny"},{"intro":"This is also","punchline":"very funny"}]' -H "Authorization: Bearer 9e44bed1ab8845508297079d4a0117b4" localhost:8080/jokes`

Example joke:

```json
{
  "id": 79,
  "intro": "Why did the crab never share?",
  "punchline": "Because he's shellfish.",
  "user": {
    "id": 1,
    "name": "Joe"
  }
}
```

## Approach and assumptions

The code is structured in such a way that there is a very low amount of coupling and dependency. Data can be provided to 
the application using any provider that implements the `model.Provider` interface. This makes it possible to have many 
different providers (HTTP API, Redis, text file, etc), and also simplifies testing by allowing mock providers to be 
implemented. The `model` package provides the middleware between the HTTP and data provider implementation layers; data
providers should implement the `model.Provider` interface in such a way that the functions accept and return values
of `model.Joke` and `model.User`, which are used by the HTTP layer.

The structure of the data assumes that each joke needs to be split into an "intro" and a "punchline", for maximum flexible
usability for the calling application.

It is also assumed that the host system (dev and prod) has Docker installed.

## Future improvements

I went slightly over the recommended time allotment of 4 hours. Given more time, the following improvements would be made:

- Input validation on the `POST /jokes` endpoint.
- Testing API endpoints with actual DB data instead of a mocked data provider.
- Passing DB credentials to the application using either command line flags, environment variables or a config file (the currently hardcoded DB credentials is horrifically poor practice).
