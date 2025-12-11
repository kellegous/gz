PROTOC_GEN_GO_VERSION := v1.36.10
PROTOC_VERSION := 33.0

GO_MOD = $(shell go list -m)

PROTOS := gz.pb.go

.PHONY: all clean test

.PRECIOUS: $(PROTOS)

all: bin/gz

bin/%: cmd/%/main.go $(PROTOS) $(shell find internal -type f -name '*.go')
	go build -o $@ ./cmd/$*

bin/protoc:
	etc/download-protoc $(PROTOC_VERSION)

bin/protoc-gen-go:
	GOBIN="$(CURDIR)/bin" go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)


%.pb.go: %.proto bin/protoc-gen-go bin/protoc
	bin/protoc --proto_path=. \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go_out=. \
		--go_opt=module=$(GO_MOD) \
		$<

clean:
	rm -rf bin $(PROTOS)

test:
	go test -v=true ./...