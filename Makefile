CicdWebserverImg ?= 192.168.1.240/beluga/cicd-webserver:develop

VERSION=$(shell git rev-parse --short HEAD)

cicd-webserver: cicd-webserver-build cicd-webserver-push

cicd-webserver-build:
	docker build --no-cache --build-arg version=$(VERSION) -t ${CicdWebserverImg} -f build/cicd-webserver/Dockerfile .

cicd-webserver-push:
	docker push ${CicdWebserverImg}
