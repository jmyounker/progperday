(define (charmichael? n)
  (verify-charmichael-iter n (- n 1)))

(define (verify-charmichael-iter n x)
  (if (= x 1) #t
    (and (fermet-test x n) (verify-charmichael-iter n (- x 1)))))

(define (fermet-test a n)
 (= a (expmod a n n)))

(define (expmod base exp m)
  (cond ((= exp 0) 1)
        ((even? exp)
         (remainder
          (square (expmod base (/ exp 2) m))
          m))
        (else
         (remainder
          (* base (expmod base (- exp 1) m))
             m))))

(define (square x) (* x x))

(charmichael? 561)
(charmichael? 1105)
(charmichael? 1729)
(charmichael? 2465)
(charmichael? 2821)
(charmichael? 6601)
