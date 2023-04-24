EBPF_LOADER := scope-ebpf
BPFTOOL ?= bpftool
CLANG ?= clang
CFLAGS := -O2 -g -Wall -Werror $(CFLAGS)
BTF_VMLINUX ?= /sys/kernel/btf/vmlinux
EBPF_DIR := internal/ebpf

ARCH := $(shell uname -m)
GOARCH := $(subst aarch64,arm64,$(subst x86_64,amd64,$(ARCH)))

GO ?= $(shell which go 2>&1)
ifeq (,$(GO))
$(error "error: \`go\` not in PATH; install or set GO to it's path")
endif

# Define a variable to store the list of Go files
GO_FILES := $(shell find . -name "*.go" ! -name "*bpfel*.go" -type f)

all: build
build: scope-ebpf

clean:
	$(RM) bin/${EBPF_LOADER}
	$(RM) internal/ebpf/vmlinux.h
	$(RM) internal/ebpf/oom/bpf_bpfel_x86.o
	$(RM) internal/ebpf/oom/bpf_bpfel_x86.go

scope-ebpf: generate
	$(GO) build -ldflags="-extldflags=-static" -o bin/${EBPF_LOADER} ./cmd/scope-ebpf

fmt:
	@for file in $(GO_FILES); do \
		$(GO) fmt $$file; \
	done

generate: export BPF_CLANG := $(CLANG)
generate: export BPF_CFLAGS := $(CFLAGS)
generate: vmlinux
	$(GO) generate internal/ebpf/oom/oom.go
# TODO:
# @$(foreach entry,$(wildcard $(EBPF_DIR)/*), \
# 	if [ -d "$(entry)" ]; then \
# 	    echo "Go generate directory: $(entry)/$(entry).go"; \
# 		$(GO) generate $(entry)/$(entry).go; \
# 	fi; \
# )

help:
	@echo "Available targets:"
	@echo "  all             - Default target, builds the scope-ebpf binary"
	@echo "  build           - Builds the scope-ebpf binary"
	@echo "  clean           - Cleans up build artifacts"
	@echo "  scope-ebpf      - Builds the scope-ebpf binary"
	@echo "  fmt             - Formats Go source files"
	@echo "  generate        - Generates Go code for ebpf programs"
	@echo "  vet             - Runs Go vet on source files"
	@echo "  vmlinux         - Generates vmlinux.h header file"

vet:
	@for file in $(GO_FILES); do \
		$(GO) vet $$file; \
	done

vmlinux:
	$(BPFTOOL) btf dump file $(BTF_VMLINUX) format c > internal/ebpf/vmlinux.h

.PHONY: all build clean fmt generate help scope-ebpf vet vmlinux
