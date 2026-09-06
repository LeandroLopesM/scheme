(define bar (vector 1 2 #\c))

(display bar)
(newline)

(vector-set! bar 2 #\!)
(display bar)
(newline)