: fac
  dup 0> if
    dup 1- recurse *
  else
    drop 1
  then
;

\ show the faculty of 5
5 fac .
