.PHONY: all

all: local

local: build run

build:
	mkdir -p dist && go build -o dist/mekong ./cmd/mekong

run:
	export MEKONG_CONFIG_FILE=dist/config.yaml && ./dist/mekong