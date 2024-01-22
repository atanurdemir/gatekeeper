NAME=dist/gatekeeper
ENTRY_POINT=cmd/main.go

.PHONY: build
## build: Compile the packages.
build:
	@go build -o $(NAME) $(ENTRY_POINT)

.PHONY: dev
## dev: Run in development mode with Air (hot-reload).
dev:
	@air -c .air.toml

.PHONY: run
## run: Build and Run in development mode.
run: build
	@./$(NAME)

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f $(NAME)

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

.PHONY: test
## test: Run tests with verbose mode
test:
	@go test -v ./tests/*

.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo