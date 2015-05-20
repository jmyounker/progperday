;; SICP 1.29 - Calculate cubes via Simpson's rule

(define (sum term a next b)
  (if (> a b)
      0
      (+ (term a)
         (sum term (next a) next b))))

(define (simp f a b n)
  (let ((h (/ (- b a) n)))
    (define (term k)
      (define (coeff k) 
        (cond ((= k 0) 1)
              ((= k n) 1)
              ((even? k) 2)
              (else 4)))
      (define (fk k)
        (f (+ a (* k h))))
      (* (coeff k) (fk k)))
    (define (inc x) (+ x 1))
    
    (* (/ h 3) (sum term 0 inc n))))

(define (cube x) (* (* x x) x))

(simp cube 0 1 100)
(simp cube 0 1 1000)