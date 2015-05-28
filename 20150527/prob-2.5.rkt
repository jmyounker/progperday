; SICP 2.5 - A really twisted representation of positive non-negative integer pairs.

(define (cons-e x y) (* (expt 2 x) (expt 3 y)))

(define (car-e r)
  (define (car-loop r x)
    (if (> (remainder r 2) 0)
        x
        (car-loop (/ r 2) (+ x 1))))
  (car-loop r 0))

(define (cdr-e r)
  (define (cdr-loop r x)
    (if (> (remainder r 3) 0)
        x
        (cdr-loop (/ r 3) (+ x 1))))
  (cdr-loop r 0))
               
(car-e (cons-e 3 2))
(cdr-e (cons-e 3 2))

(car-e (cons-e 1 5))
(cdr-e (cons-e 5 1))
