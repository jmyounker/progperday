(define (same-parity i . ls)
  (define (parity x) (if (even? x) 1 2))
  (let ((desired-parity (parity i)))
    (define (same-parity-rec in)
      (cond ((null? in) (list))
            ((= (parity (car in)) desired-parity)
             (cons (car in) (same-parity-rec (cdr in))))
            (else (same-parity-rec (cdr in)))))
  (cons i (same-parity-rec ls))))

(same-parity 1 2 3 4 5 6 7)
;(1 3 5 7)

(same-parity 2 3 4 5 6 7)
;(2 4 6)

