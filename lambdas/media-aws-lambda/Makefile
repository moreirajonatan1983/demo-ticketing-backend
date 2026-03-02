.PHONY: build-MediaFunction

build-MediaFunction:
	GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -o bootstrap ./cmd/api
	cp bootstrap $(ARTIFACTS_DIR)/
