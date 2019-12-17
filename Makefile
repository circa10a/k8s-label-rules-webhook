GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GORUN=$(GOCMD) run
PROJECT=circa10a/k8s-label-rules-webhook
BINARY=webhook

build:
	$(GOBUILD) -o $(BINARY)

run:
	$(GORUN) *.go --file sample-rules.yaml --metrics

compile:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY)-linux-amd64
	GOOS=linux GOARCH=arm go build -o bin/$(BINARY)-linux-arm
	GOOS=linux GOARCH=arm64 go build -o bin/$(BINARY)-linux-arm64
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY)-darwin-amd64

clean:
	$(GOCLEAN)
	rm -f $(BINARY)

docker-build:
	docker build -t $(PROJECT) .

docker-run:
	docker run --rm -v $(shell pwd)/sample-rules.yaml:/rules.yaml \
    -p 8080:8080 \
    $(PROJECT) --file rules.yaml --metrics

docker-dev: docker-build docker-run