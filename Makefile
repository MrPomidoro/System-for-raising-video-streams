init:
	go mod tidy

run:
	go run main.go

test:
	go test -short -count=1 -coverprofile=out.out ./...
	go tool -html=out.out
	rm out.out

.PHONY: init, run, test