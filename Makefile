build:
	docker build -t dora .

run:
	docker run --env-file=local.env -it --rm -p 8080:8080 dora

test:
	go test -v ./...

rebuild: build run

compose-run:
	docker compose run --env-file=local.env

env:
	cp default.env local.env