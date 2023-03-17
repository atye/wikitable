test:
	go test -coverpkg=./... -count=1 -race ./...
	
cover:
	go test -coverprofile=coverage.out github.com/atye/wikitable/internal/model
	go tool cover -html=coverage.out
	rm -f coverage.out