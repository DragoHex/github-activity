BINARY_NAME=gact

.PHONY: bin
bin:
	go build -o artifacts/${BINARY_NAME} cmd/main.go

.PHONY: run
run:
	go run cmd/main.go

.PHONY: clean
clean:
	go clean
	@rm -r artifacts 2> /dev/null

$(V).SILENT:
