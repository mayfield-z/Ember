BINARY_NAME=Ember

default:
	@go build -o ./build/$(BINARY_NAME) ./cmd/main/main.go

clean:
	@rm -rf ./build