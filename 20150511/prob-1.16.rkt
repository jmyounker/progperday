#lang R5RS

(define (expo b n)
  (if (evenp n)
      (expo-even b n)
      (* b (expo-even b (- n 1) ))))

(define (expo-even b n)
  (if (= n 1)
      b
  (expo-even (* b b) (/ n 2))))

(define (evenp n)
 (let ((h (/ n 2)))
      (= (floor h) h))
)

