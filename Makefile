.PHONY: build run install clean tidy

BINARY := notion-tui

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY)

install: build
	cp $(BINARY) $(GOPATH)/bin/$(BINARY) 2>/dev/null || cp $(BINARY) ~/go/bin/$(BINARY)

clean:
	rm -f $(BINARY)

tidy:
	go mod tidy
