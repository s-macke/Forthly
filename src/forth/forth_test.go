package forth

import (
	"strings"
	"testing"
)

func assertEqual(t *testing.T, left, right string) {
	right = strings.TrimSpace(right)
	left = strings.TrimSpace(left)
	if left != right {
		t.Errorf("assert failed with: '%s' is not equal to '%s'", left, right)
	}
}

func assertStackEmpty(t *testing.T, f *Forth) {
	if f.stack.Size() != 0 {
		t.Errorf("assert failed with: stack is not empty")
	}
}

func TestStack(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec("1 DUP . .")
	assertEqual(t, result, "1 1")
	assertStackEmpty(t, f)

	result, _ = f.Exec("3 4 OVER . . .")
	assertEqual(t, result, "3 4 3")
	assertStackEmpty(t, f)

	result, _ = f.Exec("3 4 DROP .")
	assertEqual(t, result, "3")
	assertStackEmpty(t, f)
}

func TestInterpretedMath(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec("1 2 + .")
	assertEqual(t, result, "3")
	assertStackEmpty(t, f)

	result, _ = f.Exec("1 2 - .")
	assertEqual(t, result, "-1")
	assertStackEmpty(t, f)

	result, _ = f.Exec("3 2 * .")
	assertEqual(t, result, "6")
	assertStackEmpty(t, f)

	result, _ = f.Exec("6 2 / .")
	assertEqual(t, result, "3")
	assertStackEmpty(t, f)
}

func TestCompiledMath(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec("1 2 + .")
	assertEqual(t, result, "3")
	assertStackEmpty(t, f)

	result, _ = f.Exec("1 2 - .")
	assertEqual(t, result, "-1")
	assertStackEmpty(t, f)

	result, _ = f.Exec("3 2 * .")
	assertEqual(t, result, "6")
	assertStackEmpty(t, f)

	result, _ = f.Exec("6 2 / .")
	assertEqual(t, result, "3")
	assertStackEmpty(t, f)
}

func TestVariable(t *testing.T) {
	f := NewForth(false)

	result, _ := f.Exec("VARIABLE var1")
	assertEqual(t, result, "")
	assertStackEmpty(t, f)

	result, _ = f.Exec("100 var1 !")
	assertEqual(t, result, "")
	assertStackEmpty(t, f)

	result, _ = f.Exec("var1 @ .")
	assertEqual(t, result, "100")
	assertStackEmpty(t, f)
}

func TestConstant(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec("42 CONSTANT FOO")
	assertEqual(t, result, "")
	assertStackEmpty(t, f)

	result, _ = f.Exec("FOO .")
	assertEqual(t, result, "42")
	assertStackEmpty(t, f)
}

func TestFunc(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec(`: square DUP * ;
                              3 square .`)
	assertEqual(t, result, "9")
	assertStackEmpty(t, f)
}

func TestCompare(t *testing.T) {
	f := NewForth(false)

	result, _ := f.Exec("0 NOT .")
	assertEqual(t, result, "-1")
	assertStackEmpty(t, f)

	result, _ = f.Exec("-1 NOT .")
	assertEqual(t, result, "0")
	assertStackEmpty(t, f)

	result, _ = f.Exec("1 NOT .")
	assertEqual(t, result, "0")
	assertStackEmpty(t, f)

	result, _ = f.Exec("5 5 = .")
	assertEqual(t, result, "-1")
	assertStackEmpty(t, f)

	result, _ = f.Exec("4 3 = .")
	assertEqual(t, result, "0")
	assertStackEmpty(t, f)

	result, _ = f.Exec("6 5 > .")
	assertEqual(t, result, "-1")
	assertStackEmpty(t, f)

	result, _ = f.Exec("6 5 < .")
	assertEqual(t, result, "0")
	assertStackEmpty(t, f)

	result, _ = f.Exec("5 5 <> .")
	assertEqual(t, result, "0")
	assertStackEmpty(t, f)

}

func TestIfThen(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec(": test IF 2 . THEN 3 . ;")
	result, _ = f.Exec("1 test")
	assertEqual(t, result, "2 3")
	assertStackEmpty(t, f)
	result, _ = f.Exec("0 test")
	assertEqual(t, result, "3")
	assertStackEmpty(t, f)
}

func TestUnless(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec(": test 0 UNLESS 2 . THEN 3 . ;")
	result, _ = f.Exec("test")
	assertEqual(t, result, "2 3")
	assertStackEmpty(t, f)
}

func TestR(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec("1 >R")
	assertEqual(t, result, "")
	assertStackEmpty(t, f)
	result, _ = f.Exec("R> .")
	assertEqual(t, result, "1")
	assertStackEmpty(t, f)

	result, _ = f.Exec("4 >R R@ R@ . .")
	assertEqual(t, result, "4 4")
	assertStackEmpty(t, f)

	result, _ = f.Exec(": rtest 1 >R 2 >R R> R> ;")
	result, _ = f.Exec("3 rtest . . .")
	assertEqual(t, result, "1 2 3")
	assertStackEmpty(t, f)
}

// from https://forth-standard.org/standard/core/CASE
func TestCase(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec(`
       : cs1 CASE
          1 OF 111 ENDOF
          2 OF 222 ENDOF
          3 OF 333 ENDOF
          >R 999 R>
          ENDCASE ;`)

	result, _ = f.Exec("1 cs1 .")
	assertEqual(t, result, "111")
	assertStackEmpty(t, f)
	result, _ = f.Exec("2 cs1 .")
	assertEqual(t, result, "222")
	assertStackEmpty(t, f)
	result, _ = f.Exec("3 cs1 .")
	assertEqual(t, result, "333")
	assertStackEmpty(t, f)
	result, _ = f.Exec("4 cs1 .")
	assertEqual(t, result, "999")
	assertStackEmpty(t, f)
}

func TestCase2(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec(`
       : cs2 >R CASE
   -1 OF CASE R@ 1 OF 100 ENDOF
                2 OF 200 ENDOF
                >R -300 R>
        ENDCASE
     ENDOF
   -2 OF CASE R@ 1 OF -99 ENDOF
                >R -199 R>
        ENDCASE
     ENDOF
     >R 299 R>
   ENDCASE R> DROP ;`)

	result, _ = f.Exec("-1 1 cs2 .")
	assertEqual(t, result, "100")
	assertStackEmpty(t, f)
	result, _ = f.Exec("-1 2 cs2 .")
	assertEqual(t, result, "200")
	result, _ = f.Exec("-1 3 cs2 .")
	assertEqual(t, result, "-300")
	result, _ = f.Exec("-2 1 cs2 .")
	assertEqual(t, result, "-99")
	result, _ = f.Exec("-2 2 cs2 .")
	assertEqual(t, result, "-199")
	result, _ = f.Exec("0 2 cs2 .")
	assertEqual(t, result, "299")
	assertStackEmpty(t, f)
}

func TestChar(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec("CHAR ! .")
	assertEqual(t, result, "33")
	assertStackEmpty(t, f)

	result, _ = f.Exec("'(' . ')' .")
	assertEqual(t, result, "40 41")
	assertStackEmpty(t, f)

	result, _ = f.Exec("'.' . '+' .")
	assertEqual(t, result, "46 43")
	assertStackEmpty(t, f)

}

func TestComment(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec("1 ( 2 3 then loop do this is my random comment ) 4 5 .s")
	assertEqual(t, result, "<3> 1 4 5")
}

func TestRecurse(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec(": rec 1- DUP . ?DUP 0 > IF RECURSE THEN ;")
	result, _ = f.Exec("4 rec")
	assertEqual(t, result, "3 2 1 0")
	assertStackEmpty(t, f)
}

func TestLoop(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec(": iter 5 0 do i . loop ;")
	result, _ = f.Exec("iter")
	assertEqual(t, result, "0 1 2 3 4")
	assertStackEmpty(t, f)
}

func TestRepeat(t *testing.T) {
	f := NewForth(false)
	// from http://lars.nocrew.org/forth2012/core/WHILE.html
	result, _ := f.Exec(": looptest BEGIN DUP 5 < WHILE DUP 1+ REPEAT ;")

	result, _ = f.Exec(" 0 looptest . . . . . .")
	assertEqual(t, result, "5 4 3 2 1 0")
	assertStackEmpty(t, f)

	result, _ = f.Exec("6 looptest .")
	assertEqual(t, result, "6")
	assertStackEmpty(t, f)
}

func TestFactorial(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec(`
        : fac
          dup 0> if
          dup 1- recurse *
          else
          drop 1
          then ;
    `)
	assertStackEmpty(t, f)
	assertEqual(t, result, "")

	result, _ = f.Exec("8 fac .")
	assertEqual(t, result, "40320")
	assertStackEmpty(t, f)
}

func TestErase(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec(`  
  CREATE BUF 20 ALLOT  
  BUF 20 ERASE
    `)
	assertStackEmpty(t, f)
	assertEqual(t, result, "")
	address := f.Find("BUF")
	for i := 0; i < 20; i++ {
		if f.heap[int(address)+i+2].(int) != 0 {
			t.Errorf("assert failed with: BUF[%d] is not zero, but %v", i, f.heap[int(address)+i])
		}
	}
}

func TestTrimString(t *testing.T) {
	f := NewForth(false)
	result, _ := f.Exec(`2 5 2 /STRING . .`)
	assertStackEmpty(t, f)
	assertEqual(t, result, "3 4")
}
