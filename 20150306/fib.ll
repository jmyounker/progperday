; ModuleID = 'fib'
target datalayout = "e-m:o-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-apple-macosx10.9.0"

@.str = private unnamed_addr constant [6 x i8] c"%llu\0A\00", align 1

; Function Attrs: ssp uwtable
define i32 @main(i32 %argc, i8** %argv) #0 {
entry:
  %argsBad = icmp ne i32 %argc, 2
  br i1 %argsBad, label %failure, label %goodtogo

goodtogo:                                        ; preds = %entry
  %argv_1_ptr = getelementptr inbounds i8** %argv, i64 1
  %argv_1 = load i8** %argv_1_ptr, align 8
  %fib_no = call i32 @atoi(i8* %argv_1)
  %fib = call i64 @fib(i32 %fib_no)
  ; yup, I'm not checking the return status from printf
  %unused = call i32 (i8*, ...)* @printf(i8* getelementptr inbounds ([6 x i8]* @.str, i32 0, i32 0), i64 %fib)
  ret i32 0

failure:                                         ; preds = %entry
  ret i32 127

}

define i64 @fib(i32 %count) {
entry:
  switch i32 %count, label %loop [
    i32 0, label %fibZero
    i32 1, label %fibOne
  ]

fibZero:                                          ; preds = %entry
  ret i64 0

fibOne:                                           ; preds = %entry
  ret i64 1

loop:                                             ; preds = %loop, %entry
  %i = phi i32 [ %count, %entry ], [ %inext, %loop ]
  %p0 = phi i64 [ 0, %entry ], [ %p1, %loop ]
  %p1 = phi i64 [ 1, %entry ], [ %pnext, %loop ]
  %pnext = add i64 %p0, %p1
  %inext = sub i32 %i, 1
  %isDone = icmp eq i32 %inext, 1
  br i1 %isDone, label %done, label %loop

done:                                             ; preds = %loop
  ret i64 %pnext
}

declare i32 @atoi(i8*) #1

declare i32 @printf(i8*, ...) #1

attributes #0 = { ssp uwtable "less-precise-fpmad"="false" "no-frame-pointer-elim"="true" "no-frame-pointer-elim-non-leaf" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { "less-precise-fpmad"="false" "no-frame-pointer-elim"="true" "no-frame-pointer-elim-non-leaf" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }

!llvm.ident = !{!0}

!0 = metadata !{metadata !"Apple LLVM version 6.0 (clang-600.0.56) (based on LLVM 3.5svn)"}

