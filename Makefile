.PHONY: run
run: build build-docker
	docker run -d --name dadz_db -e POSTGRES_PASSWORD=blah -p 5433:5432 --mount type=bind,source="$(shell pwd)"/db,target=/docker-entrypoint-initdb.d postgres:13.4
	sleep 5
	docker run --rm --name dadz -p 8080:8080 johanmeiring/dadz

.PHONY: build
build:
	env GOOS=linux GOARCH=amd64 go build -o build/dadz cmd/dadz/main.go

.PHONY: build-docker
build-docker:
	docker build -t johanmeiring/dadz .

.PHONY: clean-db
clean-db:
	docker stop dadz_db
	docker rm dadz_db

.PHONY: db
db:
	docker run --name dadz_db -e POSTGRES_PASSWORD=blah -p 5433:5432 --mount type=bind,source="$(shell pwd)"/db,target=/docker-entrypoint-initdb.d postgres:13.4

.PHONY: test-db
test-db:
	docker run -d --name dadz_db_test -e POSTGRES_PASSWORD=blah -p 5434:5432 --mount type=bind,source="$(shell pwd)"/db,target=/docker-entrypoint-initdb.d postgres:13.4

.PHONY: clean-test-db
clean-test-db:
	docker stop dadz_db_test
	docker rm dadz_db_test

.PHONY: test
test: test-db
	sleep 5
	-go test ./...
	docker stop dadz_db_test
	docker rm dadz_db_test

.PHONY: test-cover
test-cover: test-db
	sleep 5
	-go test ./... -cover
	docker stop dadz_db_test
	docker rm dadz_db_test
