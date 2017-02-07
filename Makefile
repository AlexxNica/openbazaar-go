##
## Building
##

deploy:
	./deploy.sh

build:
	./build.sh

build_linux:
	./build.sh linux/amd64

##
## docker
##
DOCKER_PROFILE ?= ob1company
DOCKER_IMAGE_NAME ?= $(DOCKER_PROFILE)/openbazaard

build_docker:
	docker build -t $(DOCKER_IMAGE_NAME) .

push_docker:
	docker push $(DOCKER_IMAGE_NAME)

docker: build_linux build_docker push_docker

##
## Cleanup
##

clean_build:
	rm -f ./dist/*

clean_docker:
	docker rmi -f $(DOCKER_IMAGE_NAME); true

clean: clean_build clean_docker
