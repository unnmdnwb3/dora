build:
	docker build -t dora .

run:
	docker run --env-file=local.env -it --rm -p 8080:8080 dora

rebuild: build run

compose:
	docker compose build
	docker compose up

copy-env:
	cp default.env local.env

unit-tests:
	go test -v ./...

mongo-run:
	docker run \
		-d \
		--name mongo-test \
		-p 27017:27017 \
		-e MONGO_INITDB_ROOT_USERNAME=user \
		-e MONGO_INITDB_ROOT_PASSWORD=password \
		mongo:latest

mongo-stop:
	docker stop mongo-test
	docker rm mongo-test

mongo-sh:
	mongosh --host 127.0.0.1 --port 27017 -u user -p password

tests: mongo-run unit-tests mongo-stop