: r13 ( c -- o )
  dup 32 or                                    \ tolower
  dup [char] a [char] z 1+ within if
    [char] m > if -13 else 13 then +
  else drop then ;


: rot13
0 do
    dup i +
    dup
    c@
    r13
    swap
    c!
loop
drop
;

: print 2dup type cr ;

s" ThisIsAString"
print

2dup rot13 print
2dup rot13 print
