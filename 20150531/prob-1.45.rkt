(define dx 0.00001)

(define (smooth f)
  (lambda (x) (/ (+ (f (- x dx) x (+ x dx)) 3))))

(define (smooth-n f n)
  (repeated (smooth f) n))

(define (compose f g) (lambda (x) (f (g x))))

(define (repeated f n)
  (cond ((= 0 n) (lambda (x) x))
        ((= 1 n) f)
        (else (compose f (repeated f (- n 1))))))

(define tolerance 0.00001)
(define tries 40)

(define (fixed-point f first-guess)
  (define (close-enough? v1 v2)
    (< (abs (- v1 v2)) tolerance))
  (define (try guess n)
    (let ((next (f guess)))
;      (display (exact->inexact guess))
;      (display "-> ")
;      (display (exact->inexact next))
;      (newline)
      (if (= n 0) -1
          (if (close-enough? guess next)
              (begin
;                (display (- tries n))
;                (display " steps")
;                (newline)
                next)
              (try next (- n 1))))))
  (try first-guess tries))


(define (damp f) (lambda (x) (/ (+ x (f x)) 2)))

;; Return equation for finding the nth root via fixed point.
(define (rfp-eq b n) (lambda (x) (/ b (expt x (- n 1)))))

(fixed-point (rfp-eq 2 1) 4)
(fixed-point (damp (rfp-eq 2 1)) 4)

(expt (fixed-point (damp (rfp-eq 5 2)) 4.0) 2)
(expt (fixed-point (damp (rfp-eq 5 3)) 4.0) 3)
(expt (fixed-point (damp (rfp-eq 5 4)) 4.0) 4)
(expt (fixed-point (damp (damp (rfp-eq 5 4))) 4.0) 4)
(expt (fixed-point (damp (damp (rfp-eq 5 5))) 4.0) 5)
(expt (fixed-point (damp (damp (rfp-eq 5 6))) 4.0) 6)
(expt (fixed-point (damp (damp (rfp-eq 5 7))) 4.0) 6)

(define (converge? dn pow)
  (begin
    (let ((v (fixed-point ((repeated damp dn) (rfp-eq 5 pow)) 4.0)))
    (display (list dn pow v))
    (not (= -1 v)))))

(converge? 0 2) ; #f
(converge? 1 2) ; #t
(converge? 1 3) ; #t
(converge? 1 4) ; #f
(converge? 2 4) ; #t
(converge? 2 5) ; #t
(converge? 2 6) ; #t
(converge? 2 7) ; #t
(converge? 2 8) ; #f
(converge? 3 8) ; #t
(converge? 4 8) ; #t
(converge? 3 15) ; #t
(converge? 3 16) ; #f
(converge? 4 16) ; #t
(converge? 4 30) ; #t
(converge? 4 31) ; #f
(converge? 5 31) ; #t
(converge? 32 64) ; #t
(converge? 8 128) ; #t
(converge? 64 128) ; #t
(converge? 64 256) ; #t
;(converge? 3 9) ; #t
;(converge? 3 10) ; #t
;(converge? 3 11) ; #t
;(converge? 3 12) ; #t
;(converge? 3 12) ; #t
