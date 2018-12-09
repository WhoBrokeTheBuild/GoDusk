
GIT_SHORT := $(shell git rev-parse --short HEAD)

EXT := 
ifeq ($(OS),Windows_NT)
	EXT := .exe
endif

TARGET := example$(EXT)

all:
	go build -ldflags "-X main.GIT_SHORT=$(GIT_SHORT)"

run: all
	./$(TARGET)
