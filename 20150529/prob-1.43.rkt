(define (compose f g) (lambda (x) (f (g x))))

(define (repeated f n)
  (if (= 1 n) f
      (compose f (repeated f (- n 1)))))

(define (inc x) (+ x 1))
(define (square x) (* x x))

((repeated inc 1) 0)
((repeated inc 2) 0)
((repeated inc 3) 0)
((repeated inc 4) 0)

((repeated square 2) 5)