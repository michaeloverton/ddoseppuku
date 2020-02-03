run:
	docker-compose -f build/docker/docker-compose.yml up

run-sentinel:
	docker-compose -f build/docker/docker-compose.yml up redis sentinel

run-target:
	docker-compose -f build/docker/docker-compose.yml up target

# adjust scale to number of lasers
run-sentinel-laser:
	docker-compose -f build/docker/docker-compose.yml up --scale laser=1 redis sentinel laser

run-all:
	docker-compose -f build/docker/docker-compose.yml up --scale laser=3 redis sentinel laser target

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