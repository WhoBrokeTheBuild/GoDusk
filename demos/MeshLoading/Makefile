CONFIG ?= debug
TAGS   ?=

TARGET := MeshLoading.bin
ifeq ($(OS),Windows_NT)
	TARGET := MeshLoading.exe
endif

SOURCE ?= $(wildcard ../../dusk/*.go ../../dusk/*/*.go *.go)

ifeq ($(CONFIG),debug)
	TAGS += debug
else 
	TAGS += release
	LDFLAGS += -s
endif

all: $(TARGET)

$(TARGET): $(SOURCE)
	go build -o $(TARGET) -tags "$(TAGS)"

run: $(TARGET)
	./$(TARGET)
