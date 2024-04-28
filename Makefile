clean:
	sudo find . -name '*.go' -exec gofmt -w {} \;

run:
	go mod download
	sudo go run .
