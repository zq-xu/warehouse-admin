WebserverImg ?= zqxu1993/warehouse-admin-webserver:v0.0.11

VERSION=$(shell git rev-parse --short HEAD)

DatabaseAddress=169.254.166.121
DatabasePort=3306
DatabaseName=warehouse_admin
DatabaseUsername=root
DatabasePassword=root

webserver: docker-build docker-push

docker-build:
	docker build --no-cache --build-arg version=$(VERSION) -t ${WebserverImg} -f build/webserver/Dockerfile .

docker-push:
	docker push ${WebserverImg}

run:
	docker run \
	--name warehouse-server \
	-d \
	-p 8080:8080 \
	-e DatabaseAddress=${DatabaseAddress} \
	-e DatabasePort=${DatabasePort} \
	-e DatabaseName=${DatabaseName} \
	-e DatabaseUsername=${DatabaseUsername} \
	-e DatabasePassword=${DatabasePassword} \
	${WebserverImg}