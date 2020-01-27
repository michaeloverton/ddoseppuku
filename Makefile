run:
	docker-compose -f build/docker/docker-compose.yml up

build: build-laser build-sentinel build-target

build-laser:
	docker build \
		-t michaeloverton/laser \
		-f build/docker/laser/Dockerfile \
		.

build-sentinel:
	docker build \
		-t michaeloverton/sentinel \
		-f build/docker/sentinel/Dockerfile \
		.

build-target:
	docker build \
		-t michaeloverton/target \
		-f build/docker/target/Dockerfile \
		.

vet:
	go vet ./...