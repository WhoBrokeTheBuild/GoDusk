CONFIG ?= debug
TAGS   ?=

TARGET := example
ifeq ($(OS),Windows_NT)
	TARGET := example.exe
endif

SOURCE ?= $(wildcard ../dusk/*.go ../dusk/obj/*.go ../dusk/fbx/*.go *.go)

GOFLAGS ?= 
LDFLAGS ?= -X main.Build=$(shell git rev-parse --short HEAD)

ifeq ($(CONFIG),debug)
	ifneq ($(OS),Windows_NT)
		# go-bindata -debug breaks on windows
    	TAGS += debug
	endif
else 
	LDFLAGS += -s
endif

all: $(TARGET)

$(TARGET):
	go build -o $(TARGET) $(GOFLAGS) -tags "$(TAGS)" -ldflags "$(LDFLAGS)"

run: $(TARGET)
	./$(TARGET)
