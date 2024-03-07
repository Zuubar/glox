package test

import "testing"

func TestFunctions(t *testing.T) {
	program1 := `
fun sum(a, b) {
	return a + b;
}

print sum(17, 3);
print sum(17, 3) + sum(10, 5);
`
	program2 := `
fun fib(n) {
  if (n <= 1) return n;
  return fib(n - 2) + fib(n - 1);
}

print fib(5);
`
	program3 := `
var closure = fun(i) {
	return i;
};

print closure(15) + closure(17);
`
	program4 := `
fun adder() {
    var i = 0;
    return fun (x) {
        i = i + x;
        print i;
    };
}

var pos = adder();
var neg = adder();

for (var i = 0; i < 5; i = i + 1) {
    pos(1);
    neg(-1);
}
`
	testPrograms(t, []testCase{
		{program1, "20\n35\n"},
		{program2, "5\n"},
		{program3, "32\n"},
		{program4, "1\n-1\n2\n-2\n3\n-3\n4\n-4\n5\n-5\n"},
	})
}
