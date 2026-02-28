.PHONY: build-EventsFunction

build-EventsFunction:
	GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -o bootstrap ./cmd/api
	cp bootstrap $(ARTIFACTS_DIR)/
