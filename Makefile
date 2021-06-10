.PHONY: run build build-win-64 build-linux-64 clean

BUILD_DIR = ./bin
BINARY_NAME = fantia_downloader

run:
	go run main.go
build: clean build-win-64 build-linux-64
build-win-64:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_x64.exe -ldflags '-w -s'
	upx --lzma $(BUILD_DIR)/$(BINARY_NAME)_x64.exe
build-linux-64:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_x64 -ldflags '-w -s'
	upx --lzma $(BUILD_DIR)/$(BINARY_NAME)_x64
clean:
	rm -rf $(BUILD_DIR)