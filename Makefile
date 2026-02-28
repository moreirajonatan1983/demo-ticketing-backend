.PHONY: build-EventsFunction build-CheckoutFunction

build-EventsFunction:
	GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -o bootstrap ./cmd/api/events
	cp bootstrap $(ARTIFACTS_DIR)/

build-CheckoutFunction:
	GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -o bootstrap ./cmd/api/checkout
	cp bootstrap $(ARTIFACTS_DIR)/
