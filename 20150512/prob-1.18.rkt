(define (mult-iter x y)
  (set-product-sign x y
     (let ((x (abs x)) (y (abs y)))
              (mult-iter-unsigned (max x y) (min x y)))))

(define (set-product-sign x y p)
  (if (= (sgn x) (sgn y)) p (- p)))

(define (sgn x)
  (cond ((negative? x) -1)
        ((positive? x) 1)
        (else 0)))

(define (mult-iter-unsigned x y)
  (cond ((= y 0) 0)
        ((= y 1) x)
        ((even? y) (mult-iter-unsigned-even x y))
        (else (+ (mult-iter-unsigned-even x y) x))))
  
(define (mult-iter-unsigned-even x y)
  (cond ((= y 0) 0)
        ((= y 1) x)
        (else (mult-iter-unsigned-even (double x) (halve y)))))

(define (double x) (+ x x))
(define (halve x) (floor (/ x 2)))

(mult-iter 5 0)
(mult-iter 5 1)
(mult-iter 5 2)
(mult-iter 5 3)
(mult-iter 5 -1)
(mult-iter -5 -3)
(mult-iter -3 5)
(mult-iter 0 -1)

