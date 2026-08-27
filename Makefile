.PHONY: build-lambda build-api build-processor

build-lambda: build-api build-processor

build-api:
	mkdir -p backend/bin/api
	cd backend && \
		GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
		go build -tags lambda.norpc -trimpath -ldflags="-s -w" \
		-o bin/api/bootstrap ./cmd/api/main.go

build-processor:
	mkdir -p backend/bin/processor
	cd backend && \
		GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
		go build -trimpath -ldflags="-s -w" \
		-o bin/processor/bootstrap ./cmd/processor/main.go
