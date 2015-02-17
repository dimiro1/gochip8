test:
	go test -cover github.com/dimiro1/gochip8...

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
