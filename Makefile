GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GORUN=$(GOCMD) run
GOBUILDFLAGS=-ldflags="-s -w"
PROJECT=circa10a/k8s-label-rules-webhook
BINARY=webhook

build:
	$(GOBUILD) -o $(BINARY)

run:
	$(GORUN) *.go --file sample-rules.yaml --metrics

compile:
	GOOS=linux GOARCH=amd64 go build $(GOBUILDFLAGS) -o bin/$(BINARY)-linux-amd64
	GOOS=linux GOARCH=arm go build $(GOBUILDFLAGS) -o bin/$(BINARY)-linux-arm
	GOOS=linux GOARCH=arm64 go build $(GOBUILDFLAGS) -o bin/$(BINARY)-linux-arm64
	GOOS=darwin GOARCH=amd64 go build $(GOBUILDFLAGS) -o bin/$(BINARY)-darwin-amd64

clean:
	$(GOCLEAN)
	rm -f $(BINARY)

# https://github.com/swaggo/gin-swagger
docs: SHELL:=/bin/bash
docs:
# Static Site
	@if ! which npm 1>/dev/null; then\
		echo "npm needs to be installed to create static site with gitbook";\
		exit 1;\
	fi;\

	@if ! which gitbook 1>/dev/null; then\
		echo "gitbook needs to be installed to create static site";\
		npm install gitbook-cli -g;\
	fi;\
	gitbook install
	gitbook build . ./docs

swagger-docs:
# Swagger
	#swag init
	#sed -i 's;"//;"/;g' docs/swagger.json


docker-build:
	docker build -t $(PROJECT) .

docker-run:
	docker run --rm -v $(shell pwd)/sample-rules.yaml:/rules.yaml \
    -p 8080:8080 \
    $(PROJECT) --file rules.yaml --metrics

docker-dev: docker-build docker-run