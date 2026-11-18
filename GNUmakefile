TEST?=./...

default: build

build:
	go build

test:
	go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

vet:
	go vet $(TEST)

testacc:
	TF_LOG=debug KSQL_SERVER_URL=localhost:8088 TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

.PHONY: default build test vet testacc
