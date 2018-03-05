VERSION := $(shell cat VERSION)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
-include .env
.PHONY: version


# -
# Local
# -


fast:
	go build $(LDFLAGS) -o samp-servers-api

static:
	CGO_ENABLED=0 GOOS=linux go build -a $(LDFLAGS) -o samp-servers-api .

local: fast
	DEBUG=1 \
	./samp-servers-api

version:
	git tag $(VERSION)
	git push
	git push origin $(VERSION)

test:
	go test -v -race


# -
# Docker
# -


build:
	docker build --no-cache -t southclaws/samp-servers:$(VERSION) .

push:
	docker push southclaws/samp-servers:$(VERSION)
	
run:
	-docker stop samp-servers-api
	-docker rm samp-servers-api
	docker run \
		--name samp-servers-api \
		--network host \
		--env-file .env \
		southclaws/samp-servers:$(VERSION)

run-prod:
	-docker stop samp-servers-api
	-docker rm samp-servers-api
	docker run \
		--name samp-servers-api \
		--detach \
		--publish 7790:80 \
		--restart on-failure:10 \
		--env-file .env \
		southclaws/samp-servers:$(VERSION)
	docker network connect samp-servers samp-servers-api


# -
# Testing
# -


mongodb:
	-docker stop mongodb
	-docker rm mongodb
	docker run \
		--name mongodb \
		-p 27017:27017 \
		-d \
		mongo
