TARGET := Lighting$(shell go env GOEXE)

SOURCE ?= $(wildcard ../../dusk/*.go ../../dusk/*/*.go *.go)

all: $(TARGET)

$(TARGET): $(SOURCE)
	go build -o $(TARGET)

run: $(TARGET)
	./$(TARGET)
