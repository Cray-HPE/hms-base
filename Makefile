NAME ?= hms-base 
VERSION ?= $(shell cat .version)

all : image unittest coverage

image:
		docker build --pull ${DOCKER_ARGS} --tag '${NAME}:${VERSION}' .

unittest: buildbase
		docker build -t cray/hms-base-testing -f Dockerfile.testing .

coverage:
		docker build -t cray/hms-base-coverage -f Dockerfile.coverage .

buildbase: buildbase
		docker build -t cray/hms-base-build-base -f Dockerfile.build-base .