(#%require (only racket/base random))

(define (make-interval a b) (cons (min a b) (max a b)))

(define (lower-bound x) (car x))
(define (upper-bound x) (cdr x))

(define (add-interval x y)
  (make-interval (+ (lower-bound x)
                    (lower-bound y))
                 (+ (upper-bound x)
                    (upper-bound y))))

(define (sub-interval x y)
  (add-interval x (neg-interval y)))

(define (mul-interval x y)
  (define (upper-interval? i) (and (>= (lower-bound i) 0) (>= (upper-bound i) 0)))
  (define (lower-interval? i) (and (< (lower-bound i) 0) (< (upper-bound i) 0)))
  (define (split-interval? i) (and (< (lower-bound i) 0) (>= (upper-bound i) 0)))
  (cond ((and (upper-interval? x) (upper-interval? y))
         (make-interval
          (* (lower-bound x) (lower-bound y))
          (* (upper-bound x) (upper-bound y))))
        
         ((and (lower-interval? x) (lower-interval? y))
          (make-interval
           (* (upper-bound x) (upper-bound y))
           (* (lower-bound x) (lower-bound y))))
         
         ((and (lower-interval? x) (upper-interval? y))
          (make-interval
           (* (lower-bound x) (upper-bound y))
           (* (upper-bound x) (lower-bound y))))
         
         ((and (upper-interval? x) (lower-interval? y))
          (make-interval
           (* (upper-bound x) (lower-bound y))
           (* (lower-bound x) (upper-bound y))))
         
         ((and (split-interval? x) (upper-interval? y))
          (make-interval
           (* (lower-bound x) (upper-bound y))
           (* (upper-bound x) (upper-bound y))))
         
         ((and (upper-interval? x) (split-interval? y))
          (make-interval
           (* (upper-bound x) (lower-bound y))
           (* (upper-bound x) (upper-bound y))))
         
         ((and (lower-interval? x) (split-interval? y))
          (make-interval
           (* (lower-bound x) (upper-bound y))
           (* (lower-bound x) (lower-bound y))))
         
         ((and (split-interval? x) (lower-interval? y))
          (make-interval
           (* (upper-bound x) (lower-bound y))
           (* (lower-bound x) (lower-bound y))))
          (else 
           (let ((p1 (* (lower-bound x)
                        (lower-bound y)))
                 (p2 (* (lower-bound x)
                        (upper-bound y)))
                 (p3 (* (upper-bound x)
                        (lower-bound y)))
                 (p4 (* (upper-bound x)
                        (upper-bound y))))
             (make-interval (min p1 p2 p3 p4)
                            (max p1 p2 p3 p4))))))

(define (mul-exp-interval x y) 
  (let ((p1 (* (lower-bound x)
               (lower-bound y)))
        (p2 (* (lower-bound x)
               (upper-bound y)))
        (p3 (* (upper-bound x)
               (lower-bound y)))
        (p4 (* (upper-bound x)
               (upper-bound y))))
    (make-interval (min p1 p2 p3 p4)
                   (max p1 p2 p3 p4))))
       
(define (recip-interval x)
  (make-interval (/ 1.0 (upper-bound x))
                 (/ 1.0 (lower-bound x))))

(define (neg-interval x)
  (make-interval (- (upper-bound x))
                 (- (lower-bound x))))

(define (div-interval x y)
  (if (contains-interval? y 0) #f
      (mul-interval x (recip-interval y))))

(define (contains-interval? x a)
  (and (<= (lower-bound x) a) (<= a (upper-bound x))))

(define (make-center-percent c t)
  (let ((w (* c (/ t 100)))) 
  (make-interval (- c w) (+ c w))))

(define (center-interval i) (/ (+ (upper-bound i) (lower-bound i)) 2))

(define (percent-interval i)
  (let ((c (center-interval i)))
        (abs (/ (* (- (upper-bound i) c) 100) c))))

(define (equal-interval? x y)
  (and (= (lower-bound x) (lower-bound y))
       (= (upper-bound x) (upper-bound y))))

(equal-interval? (make-center-percent 10 50) (make-interval 5 15))

(= (center-interval (make-center-percent 10 50)) 10)

(= (center-interval (make-center-percent -10 50)) -10)

(= (percent-interval (make-center-percent 10 50)) 50)

(= (percent-interval (make-center-percent -10 50)) 50)

(equal-interval? 
 (add-interval (make-interval 2 3) (make-interval 7 8))
 (make-interval 9 11))

(equal-interval? 
 (mul-interval (make-interval 2 3) (make-interval 7 8))
 (make-interval 14 24))

(equal-interval? 
 (div-interval (make-interval 2 3) (make-interval 7 8))
 (make-interval 0.25 0.42857142857142855))

(equal-interval?
 (recip-interval (make-interval 2 3))
 (make-interval (/ 1 3) 0.5))

(equal-interval?
 (neg-interval (make-interval 2 3))
 (make-interval -3 -2))

(equal-interval?
 (sub-interval (make-interval 2 3) (make-interval 7 8))
 (make-interval -6 -4))

(equal-interval?
 (sub-interval (make-interval 2 3) (make-interval 2 3))
 (make-interval -1 1))

(not (div-interval (make-interval 2 3) (make-interval -0.5 0.5)))

;; Faster to code than to come up with tall the test cases.
(define (test-mul-interval n)
  (define (rnd) (- (random 10000) 5000))
  (if (= n 0) #t
      (let ((i1 (make-interval (rnd) (rnd)))
            (i2 (make-interval (rnd) (rnd))))
        (if (not (equal-interval? (mul-interval i1 i2) (mul-exp-interval i1 i2)))
            (begin
              (display (list n i1 i2 ":" (mul-interval i1 i2) (mul-exp-interval i1 i2)))
              (newline)
              #f)
            (test-mul-interval (- n 1))))))

(test-mul-interval 200)
