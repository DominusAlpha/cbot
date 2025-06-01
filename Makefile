VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
REGISTRY=chigabiga
TARGETOS=linux
TARGETARCH=amd64
APP=$(shell basename $(shell git remote get-url origin))

format: 
	gofmt -s -w ./

get:
	go get

build: get
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o cbot -ldflags "-X="github.com/dominusalpha/cbot/cmd.appVersion=${VERSION}

image:
	docker --debug build . -t ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}

push:
	docker push ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}

clean:
	rm -rf cbot