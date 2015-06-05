(define (reverse-list ls)
  (define (reverse in out)
    (if (null? in) out
        (reverse (cdr in) (cons (car in) out))))
  (reverse ls (list)))

(reverse-list (list))
(reverse-list (list 1))
(reverse-list (list 1 2 3 4 5))
