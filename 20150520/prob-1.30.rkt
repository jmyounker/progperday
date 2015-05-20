;; SICP 1.30 - Transform "sum" function from a recursive to an interative process.
;;
;; A recursive process returns its initial state from the final recursion call.
;; An iterative process returns its result from the final recursion call.
;;
;; Put another way:
;;
;; A recursive process computes its result on the way back up the recursion tree.
;; An iterative process computes its result on the way down the recursion tree.

(define (sum term a next b)
  (define (iter a result)
    (if (> a b)
      result
      (iter (next a) (+ (term a) result))))
  (iter a 0))

;; The "simp" function is from the previous exercise, 1.29.
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