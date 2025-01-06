all: run

run:
	go run ./cmd/... -config=dev.yml

build:
	CGO_ENABLED=0 go build -o push2sentinel ./cmd/...

test:
	go test -v ./...

clean:
	rm -r dist/ push2sentinel || true

update:
	go get -u ./...
	go mod tidy
