VERSION := $(shell git describe --always --tags --dirty)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
-include .env
.PHONY: version


# -
# Local
# -


static:
	go get
	CGO_ENABLED=0 GOOS=linux go build -a $(LDFLAGS) -o samp-servers-api .

fast:
	go build $(LDFLAGS) -o samp-servers-api

local: fast
	DEBUG=1 \
	./samp-servers-api

version:
	git tag $(VERSION)
	git push
	git push origin $(VERSION)

test:
	go get
	go test -v -race


# -
# Docker
# -

build:
	docker build --no-cache -t southclaws/samp-servers-api:$(VERSION) .

push:
	docker push southclaws/samp-servers-api:$(VERSION)
	docker tag southclaws/samp-servers-api:$(VERSION) southclaws/samp-servers-api:latest
	docker push southclaws/samp-servers-api:latest
	
run:
	-docker stop samp-servers-api
	-docker rm samp-servers-api
	docker run \
		--name samp-servers-api \
		--network host \
		--env-file .env \
		southclaws/samp-servers-api:$(VERSION)


# -
# Testing
# -

mongodb:
	-docker stop mongodb
	-docker rm mongodb
	-docker stop express
	-docker rm express
	docker run \
		--name mongodb \
		--publish 27017:27017 \
		--detach \
		mongo
	sleep 5
	docker run \
		--name express \
		--publish 8081:8081 \
		--link mongodb:mongo \
		--detach \
		mongo-express
