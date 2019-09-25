
all: build

.PHONY: build
build:
	go build -o datagen cmd/datagen/main.go

