# NOTE: llvm is installed with 'brew install llvm'
# NOTE: this is only known to work with go 1.4

GOROOT=/usr/local/go
GOPATH=$(shell cd -- "../../../.." && pwd)
GO=$(GOROOT)/bin/go

LLVM_CONFIG=/usr/local/opt/llvm/bin/llvm-config
LLVM_CFLAGS=$(shell $(LLVM_CONFIG) --cflags)
LLVM_LDFLAGS=$(shell $(LLVM_CONFIG) --ldflags)
LLVM_LIBDIR=$(shell $(LLVM_CONFIG) --libdir)
LLVM_VER=$(shell $(LLVM_CONFIG) --version)

CGO_CFLAGS=$(LLVM_CFLAGS) -I ../include
CGO_LDFLAGS=$(LLVM_LDFLAGS) -Wl,-L $(LLVM_LIBDIR) -lLLVM-$(LLVM_VER)

TAGS=-tags llvm$(LLVM_VER)

GO_ENV=GOROOT=$(GOROOT) GOPATH=$(GOPATH) CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)"

clean:
	$(GO_ENV) $(GO) clean -i github.com/axw/gollvm/llvm
	rm -f parse

get:
	$(GO_ENV) $(GO) get $(TAGS) $* github.com/axw/gollvm/llvm

build:
	$(GO_FLAGS) $(GO) build

test:
	$(GO_FLAGS) $(GO) test

all: clean get build test

.PHONY: all build clean get test
