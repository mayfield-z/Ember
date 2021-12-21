BINARY_NAME=Ember
DPDK_TEST_NAME=dpdktest

default:
	@go build -o ./build/$(BINARY_NAME) ./cmd/main/main.go

dpdk_test:
	@go build -o ./build/$(DPDK_TEST_NAME) ./cmd/dpdk_test/main.go

clean:
	@rm -rf ./build