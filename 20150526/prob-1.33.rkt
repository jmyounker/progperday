(#%require (only racket/base random))

(define (facc filter combiner null-value term a next b)
  (define (facc-iter r a)
    (cond ((> a b) r) 
          ((filter a) (facc-iter (combiner r (term a)) (next a)))
          (else (facc-iter r (next a)))))
  (facc-iter null-value a))

;; quick and dirty hack to allow testing
(define (prime? n)
  (cond ((= n 2) #t)
        ((= n 3) #t)
        ((= n 5) #t)
        ((= n 7) #t)
        ((= n 11) #t)
        ((= n 13) #t)
        (else #f)))

(define (square x) (* x x))
(define (inc x) (+ x 1))

(define (sum-sq-primes a b)
  (facc prime? + 0 square a inc b))

(define (prod-gcd n) 
  (facc (lambda (x) (= (GCD x n) 1)) * 1 (lambda (x) x) 1 inc (- n 1)))

(sum-sq-primes 1 10)
(prod-gcd 5)
