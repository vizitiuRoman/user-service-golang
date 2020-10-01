clean:
	@go clean ./...
	@go mod tidy

build:
	cd cmd && go build -trimpath -gcflags=-trimpath=%CD% -asmflags=-trimpath=%CD% -ldflags "-s -w"

test:
	@go test -race -v -p 2 -coverpkg=./... -covermode=atomic -coverprofile cover.out.tmp ./...

lint:
	@golangci-lint run ./...
	@cd pkg && golangci-lint run ./...

docker-build:
	docker build -t vizitiuroman/user-service-golang .
	docker tag vizitiuroman/user-service-golang vizitiuroman/user-service-golang:1.0.0

