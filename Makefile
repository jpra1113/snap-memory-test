init:
	glide install
build-release:
	CGO_ENABLED=0 go build -a -installsuffix cgo
build-local:
	go build .
docker-build: 
	sudo docker build -t jpra1113/snap:memory-test .
docker-push: build-release docker-build
	sudo docker push jpra1113/snap:memory-test
