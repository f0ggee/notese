test: go test ./..


lint: golangci-test run



check: fmt lint run