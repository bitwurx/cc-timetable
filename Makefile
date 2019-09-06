.PHONY: build
build:
	@docker run \
		--rm \
		-e CGO_ENABLED=0 \
		-v $(PWD):/usr/src/concord-timetable \
		-w /usr/src/concord-timetable \
		golang /bin/sh -c "go get -v -d && go build -a -installsuffix cgo -o main"
	@docker build -t concord/timetable .
	@rm -f main

.PHONY: test
test:
	@docker run \
		-d \
		-e ARANGO_ROOT_PASSWORD=abc123 \
		--name concord-timetable_test__arangodb \
		arangodb/arangodb
	@docker run \
		-d \
		-e ARANGODB_HOST=http://arangodb:8529 \
		-e ARANGODB_NAME=test__concord_timetable \
		-e ARANGODB_USER=root \
		-e ARANGODB_PASS=abc123 \
		-v $(PWD):/go/src/concord-timetable \
		-v $(PWD)/.src:/go/src \
		-w /go/src/concord-timetable \
		--link concord-timetable_test__arangodb:arangodb \
		--name concord-timetable_test \
		golang /bin/sh -c "go get -v -t -d && go test -v -coverprofile=.coverage.out"
	@docker logs -f concord-timetable_test
	@docker rm -f concord-timetable_test
	@docker rm -f concord-timetable_test__arangodb

.PHONY: test-short
test-short:
	@docker run \
		--rm \
		-v $(PWD):/go/src/concord-timetable \
		-v $(PWD)/.src:/go/src \
		-w /go/src/concord-timetable \
		golang /bin/sh -c "go get -v -t -d && go test -short -v -coverprofile=.coverage.out"
