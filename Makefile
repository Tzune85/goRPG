cover:
	go test ./game/... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test:
	go test ./game/... -v

build:
	go build -o goRpg ./...
