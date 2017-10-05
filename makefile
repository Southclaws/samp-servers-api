VERSION := $(shell git rev-parse HEAD)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
MONGO_PASS := $(shell cat MONGO_PASS.private)

.PHONY: version

fast:
	go build $(LDFLAGS) -o samp-servers-api

static:
	CGO_ENABLED=0 GOOS=linux go build -a $(LDFLAGS) -o samp-servers-api .

local: fast
	export BIND=localhost:7790
	export MONGO_USER=samplist
	export MONGO_PASS=$(MONGO_PASS)
	export MONGO_HOST=southcla.ws
	export MONGO_PORT=27017
	export MONGO_NAME=samplist
	export MONGO_COLLECTION=servers
	export QUERY_INTERVAL=0
	export MAX_FAILED_QUERY=0
	export VERIFY_BY_HOST=0
	./main

version:
	git tag $(VERSION)
	git push
	git push origin $(VERSION)

# Docker

build:
	docker build -t southclaws/samp-servers:$(VERSION) .

push: build
	docker push southclaws/samp-servers:$(VERSION)
	
run:
	-docker rm samp-servers-test
	docker run \
		--name samp-servers-test \
		-e BIND=localhost:7790 \
		-e MONGO_USER=samplist \
		-e MONGO_PASS=$(MONGO_PASS) \
		-e MONGO_HOST=southcla.ws \
		-e MONGO_PORT=27017 \
		-e MONGO_NAME=samplist \
		-e MONGO_COLLECTION=servers \
		-e QUERY_INTERVAL=30 \
		-e MAX_FAILED_QUERY=100 \
		-e VERIFY_BY_HOST=0 \
		southclaws/samp-servers:$(VERSION)

enter:
	docker run -it --entrypoint=bash southclaws/samp-servers:$(VERSION)

enter-mount:
	docker run -v $(shell pwd)/testspace:/samp -it --entrypoint=bash southclaws/samp-servers:$(VERSION)
