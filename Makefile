compose:
	docker compose build
	docker compose up

copy-test-env:
	cp .test.env .env

mongo-run:
	docker run \
		-d \
		--name mongo-test \
		-p 27017:27017 \
		-e MONGO_INITDB_ROOT_USERNAME=user \
		-e MONGO_INITDB_ROOT_PASSWORD=password \
		mongo:latest | :

mongo-stop:
	docker stop mongo-test
	docker rm mongo-test

mongo-sh:
	mongosh --host 127.0.0.1 --port 27017 -u user -p password

ginkgo:
	ginkgo -r -v --randomize-all --randomize-suites

tests: mongo-run ginkgo mongo-stop