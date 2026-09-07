(define bar (cons 123 (cons "X" #\y)))

(display "Bar is " (car bar) ", " (car (cdr bar)))