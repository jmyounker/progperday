(define (gcd-fast x y)
    (gcd-fast-iter (max (abs x) (abs y)) (min (abs x) (abs y))))

(define (gcd-fast-iter x y)                            
  (let ((r (remainder x y)))
    (if (= r 0)
        y
        (gcd-fast-iter y r))))

(= (gcd-fast 24 16) 8)
(= (gcd-fast 15 1) 1)
(= (gcd-fast 45 9) 9)
(= (gcd-fast 18 45) 9)
(= (gcd-fast 13 7) 1)
(= (gcd-fast 15 10) 5)
(= (gcd-fast 15 -10) 5)
(= (gcd-fast -15 10) 5)
(= (gcd-fast -15 -10) 5)
