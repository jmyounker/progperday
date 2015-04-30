#!/bin/sh
export PATH=$PATH:/usr/local/Cellar/llvm/3.5.1/bin

ver="`llvm-config --version`"
export CGO_CFLAGS="`llvm-config --cflags` -I ../include"
export CGO_LDFLAGS="`llvm-config --ldflags` -Wl,-L`llvm-config --libdir` -lLLVM-$ver"

case "$ver" in
*svn)
	tags="-tags llvmsvn"
	;;
*)
	tags="-tags llvm$ver"
	;;
esac

go clean -i github.com/axw/gollvm/llvm
go get $tags $* github.com/axw/gollvm/llvm
