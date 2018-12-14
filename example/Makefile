GIT_SHORT := $(shell git rev-parse --short HEAD)

GOFLAGS ?= -tags debug
LDFLAGS ?= -ldflags "-X main.Build=$(GIT_SHORT)"

EXT := 
ifeq ($(OS),Windows_NT)
	EXT := .exe
endif

TARGET := example$(EXT)

all:
	go build $(GOFLAGS) $(LDFLAGS)

run: all
	./$(TARGET)
