(define (mult x y)
  (set-product-sign x y
     (let ((x (abs x)) (y (abs y)))
              (mult-iter 0 (max x y) (min x y)))))

(define (set-product-sign x y p)
  (if (= (sgn x) (sgn y)) p (- p)))

(define (sgn x)
  (cond ((negative? x) -1)
        ((positive? x) 1)
        (else 0)))

(define (mult-iter a x y)
  (cond ((= y 0) a)
        ((even? y) (mult-iter a (double x) (halve y)))
        (else (mult-iter (+ a x) x (- y 1)))))
  
(define (double x) (+ x x))
(define (halve x) (/ x 2))

(mult 5 0)
(mult 5 1)
(mult 5 2)
(mult 5 3)
(mult 5 4)

(mult 5 -1)
(mult -5 -1)
(mult -5 1)
(mult 0 -1)

