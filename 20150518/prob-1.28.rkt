(#%require (only racket/base random))

(define (prime? n)
  (if (= (remainder n 2) 0) #f
      (miller-raben-test (+ (random (- n 2)) 1) n)))

(define (miller-raben-test a n)
 (=
  (expmod a (- n 1) n)
  1))


(define (expmod base exp m)
(signal-non-trivial-sqrt
  (cond ((= exp 0) 1)
         ((even? exp)
           (remainder
            (square (expmod base (/ exp 2) m))
            m))
        (else
         (remainder
          (* base (expmod base (- exp 1) m))
             m)))
  exp m))

(define (signal-non-trivial-sqrt x exp m)
  (cond ((= x 1) x)
        ((= x (- m 1)) x)
        (else
         (if (= (remainder (square x) m) 1) 0 x))))

(define (square x) (* x x))

(prime? 561)
(prime? 1105)
(prime? 55)
(prime? 53)
(prime? 7)
(prime? 29)
