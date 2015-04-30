GOROOT=/usr/local/go
GOPATH=$(shell cd -- "../../../.." && pwd)
GO=$(GOROOT)/bin/go

LLVM_CONFIG=/usr/local/opt/llvm/bin/llvm-config

VER=$(shell $(LLVM_CONFIG) --version)

LLVM_CFLAGS=$(shell $(LLVM_CONFIG) --cflags)
LLVM_LDFLAGS=$(shell $(LLVM_CONFIG) --ldflags)
LLVM_LIBDIR=$(shell $(LLVM_CONFIG) --libdir)

CGO_CFLAGS=$(LLVM_CFLAGS) -I ../include

CGO_LDFLAGS=$(LLVM_LDFLAGS) -Wl,-L $(LLVM_LIBDIR) -lLLVM-$(VER)

TAGS=-tags llvm$(VER)

GO_FLAGS=GOROOT=$(GOROOT) GOPATH=$(GOPATH) CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)"

clean:
	$(GO_FLAGS) $(GO) clean -i github.com/axw/gollvm/llvm
	rm -f parse

get:
	$(GO_FLAGS) $(GO) get $(TAGS) $* github.com/axw/gollvm/llvm

build:
	$(GO_FLAGS) $(GO) build

all: get build

.PHONY: all build clean get
