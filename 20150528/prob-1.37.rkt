(define (cont-frac n d k)
  (define (cont-frac-iter r k)
    (if (= k 0) r
        (cont-frac-iter (/ (n k) (+ (d k) r)) (- k 1))))
  (cont-frac-iter 0 k))

(define (gmean-rcp k)
           (cont-frac (lambda (i) 1.0)
                      (lambda (i) 1.0)
                      k))

(define (loop f a b)
  (if (not (= a b)) (begin (f a) (loop f (+ a 1) b))))

(loop (lambda (x) (display (gmean-rcp x)) (newline)) 1 10)

