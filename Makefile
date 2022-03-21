test:
	go clean -testcache &&  go test ./...

ci-test:
	scripts/ci-test.sh

lint:
	./bin/golangci-lint run -v --timeout 5m --out-format checkstyle | tee golangci-lint-report.xml
