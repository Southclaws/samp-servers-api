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
	export MONGODB_USER=samplist
	export MONGODB_PASS=$(MONGO_PASS)
	export MONGODB_HOST=southcla.ws
	export MONGODB_PORT=27017
	export MONGODB_NAME=samplist
	export MONGODB_COLLECTION=servers
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
	docker run \
		--name samp-servers-test \
		-e BIND=localhost:7790 \
		-e MONGODB_USER=samplist \
		-e MONGODB_PASS=$(MONGO_PASS) \
		-e MONGODB_HOST=southcla.ws \
		-e MONGODB_PORT=27017 \
		-e MONGODB_NAME=samplist \
		-e MONGODB_COLLECTION=servers \
		-e QUERY_INTERVAL=0 \
		-e MAX_FAILED_QUERY=0 \
		-e VERIFY_BY_HOST=0 \
		southclaws/samp-servers:$(VERSION)

enter:
	docker run -it --entrypoint=bash southclaws/samp-servers:$(VERSION)

enter-mount:
	docker run -v $(shell pwd)/testspace:/samp -it --entrypoint=bash southclaws/samp-servers:$(VERSION)
