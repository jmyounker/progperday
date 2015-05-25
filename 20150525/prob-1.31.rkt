(define (prod term a next b)
  (if (> a b) 1
      (* (term a) (prod term (next a) next b))))

(define (square x) (* x x))
(define (inc x) (+ x 1))

; part a
(define (fac n)
  (prod (lambda (x) x) 1 inc n))

(fac 1)
(fac 2)
(fac 3)
(fac 4)
(fac 5)
(fac 6)

(define (pi n)
  (* (/ 8 3)
     (prod
      (lambda (x) (let ((xs (square x))) (/ (* 4 xs) (- (* 4 xs) 1))))
      2 inc (+ 2 n))))

(pi 1)
(pi 2)
(pi 3)
(pi 4)
(pi 5)
(pi 6)

; b
(define (prodi term a next b)
  (define (prod-iter s term a next b)
    (if (> a b) s (prod-iter (* s (term a)) term (next a) next b)))
  (prod-iter 1 term a next b))

(define (faci n)
  (prodi (lambda (x) x) 1 inc n))

(faci 1)
(faci 2)
(faci 3)
(faci 4)
(faci 5)
(faci 6)

