default: fib

all: fib tmpl

fib: fib.s
	clang -x assembler fib.s -o fib

fib.s: fib.bc
	/usr/local/Cellar/llvm/3.5.1/bin/llc fib.bc -o fib.s
	# Get rid of directive that breaks clang.  There's probably a real
	# solution to this, but I'm not spending the time for this exercise.
	sed -i '' '/\.macosx_version_min/d' fib.s

fib.bc: fib.ll
	/usr/local/Cellar/llvm/3.5.1/bin/llvm-as fib.ll -o fib.bc


# tmpl is a tet program used to generate assembly skeleton & verify
# the build process
tmpl: tmpl.s
	clang -x assembler tmpl.s -o tmpl

tmpl.bc: tmpl.ll
	/usr/local/Cellar/llvm/3.5.1/bin/llvm-as tmpl.ll -o tmpl.bc

tmpl.s: tmpl.bc
	/usr/local/Cellar/llvm/3.5.1/bin/llc tmpl.bc -o tmpl.s
	# Get rid of directive that breaks clang.  There's probably a real
	# solution to this, but I'm not spending the time for this exercise.
	sed -i '' '/\.macosx_version_min/d' tmpl.s

tmpl.ll: tmpl.C
	clang -S -emit-llvm tmpl.C -o tmpl.ll

clean:
	rm -f tmpl tmpl.s tmpl.bc tmpl.ll
	rm -f fib fib.s fib.bc

.phony: all clean default
