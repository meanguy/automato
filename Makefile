.PHONY: coverage
coverage:
	go test -coverprofile=build/coverage.out ./...
	go tool cover -html=build/coverage.out -o build/coverage.html
	python3 -m http.server || true
