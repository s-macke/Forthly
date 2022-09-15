: '\n' 10 ;
: cr '\n' emit ;
: BL   32 ;
: SPACE BL EMIT ;

( create booleans )
: true 1 ;
: false 0 ;

: not false = ;

( no-operation )
: nop ;

( create basic arithmetic )

: NEGATE 0 swap - ;

: 1+ 1 + ;
: 1- 1 - ;

: >     swap < ;
: <=    > not ;
: >=    < not ;
: <>    = not ;

: 0=    0 = ;
: 0<>   0 <> ;
: 0<    0 < ;
: 0>    0 > ;
: 0<=   0 <= ;
: 0>=   0 >= ;


( Variables and constant )
( "name" -- )
: variable create 0 , ;

( n "name" -- )
( : constant create , does> @ ; )