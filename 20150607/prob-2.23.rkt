(define (for-each-mine f items)
  (if (not (null? items))
      (begin 
        (f (car items))
        (for-each-mine f (cdr items)))))


(for-each-mine (lambda (x) (newline) (display x))
          (list 57 321 88))
;57
;321
;88


