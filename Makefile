run:
	docker-compose -f build/docker/docker-compose.yml up

run-sentinel:
	docker-compose -f build/docker/docker-compose.yml up redis sentinel

# adjust --scale to modify number of lasers
run-lasers:
	docker-compose -f build/docker/docker-compose.yml up --scale laser=5 redis laser

run-target:
	docker-compose -f build/docker/docker-compose.yml up target

# adjust --scale to modify number of lasers
run-all:
	docker-compose -f build/docker/docker-compose.yml up --scale laser=5 redis sentinel laser target

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