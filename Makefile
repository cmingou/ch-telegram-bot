.PHONY: all build clean

all: build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ch-telegram-bot-linux main.go token.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ch-telegram-bot-mac main.go token.go

clean:
	@echo "\033[32m----- Clear all environment -----\033[0m"
	rm ch-telegram-*