: ?Prime
  dup
  2 do
    dup i mod
    0= if drop false unloop exit then
  loop
  drop true
;

: .bool
    IF
    ." true" cr
    ELSE
    ." false" cr
    THEN
;

: primeloop
  3 do
    i dup ?Prime OVER . .bool drop
  loop ;

20 primeloop


