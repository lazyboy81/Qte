APP_NAME=tqt
CMD=./cmd/tqt
BUILD_DIR=build

.PHONY: build build-all clean

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD)

build-all:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME).exe $(CMD)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD)

clean:
	rm -rf $(BUILD_DIR)

run:
	go run $(CMD)/main.go