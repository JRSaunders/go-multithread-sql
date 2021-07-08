DOCKER_USER ?=
DOCKER_PASS_FILE ?=
DOCKER_PASS ?=
TAG ?= $$(git log --format="%H" -n 1)
PHP_FPM_PORT ?= 9000
NGINX_PORT ?= 8080
VERSION ?= latest

login:
ifdef DOCKER_PASS_FILE
	cat $(DOCKER_PASS_FILE) | docker login -u ${DOCKER_USER} --password-stdin
else
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}
endif

build:
	docker build -t jrsaunders/gothreaded:${TAG} . \
	&& docker tag jrsaunders/gothreaded:${TAG} jrsaunders/gothreaded:latest

build-new:
	docker build --platform amd64 -t jrsaunders/gothreaded:amd64 .

push:
	docker push jrsaunders/gothreaded

build-go:
	cd src && env GOOS=linux GOARCH=amd64 go build -o ../gothreaded

build-go-osx:
	cd src && go build -o ../gothreaded-osx



