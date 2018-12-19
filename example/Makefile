CONFIG ?= debug
TAGS   ?=

TARGET := example$(shell go env GOEXE)
SOURCE ?= $(wildcard ../dusk/*.go ../dusk/*/*.go *.go)

GOFLAGS ?= 
LDFLAGS ?= -X main.Build=$(shell git rev-parse --short HEAD)

ifeq ($(CONFIG),debug)
	# go-bindata -debug breaks on windows
	TAGS += debug
else 
	TAGS += release
	LDFLAGS += -s
endif

all: $(TARGET)

$(TARGET): $(SOURCE)
	go build -o $(TARGET) $(GOFLAGS) -tags "$(TAGS)" -ldflags "$(LDFLAGS)"

run: $(TARGET)
	./$(TARGET)
